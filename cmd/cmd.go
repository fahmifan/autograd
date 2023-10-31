package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/config"
	db "github.com/fahmifan/autograd/db/migrations"
	"github.com/fahmifan/autograd/fs"
	"github.com/fahmifan/autograd/httpsvc"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/auth/auth_cmd"
	"github.com/fahmifan/autograd/pkg/core/user_management/user_management_cmd"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	"github.com/fahmifan/autograd/pkg/pb/autograd/v1/autogradv1connect"
	"github.com/fahmifan/autograd/pkg/service"
	"github.com/fahmifan/autograd/repository"
	"github.com/fahmifan/autograd/usecase"
	"github.com/fahmifan/autograd/worker"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Execute() error {
	var rootCmd = &cobra.Command{
		Use:   "autograd",
		Short: "Autograd is a auto grader for programming assignment",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.AddCommand(serverCmd())
	rootCmd.AddCommand(adminCmd())
	rootCmd.AddCommand(loginCmd())

	return rootCmd.Execute()
}

func mustInitService() *service.Service {
	gormDB := db.MustSQLite()
	return service.NewService(gormDB, config.JWTKey(), core.MediaConfig{
		RootFolder:   config.FileUploadPath(),
		ObjectStorer: fs.NewLocalStorage(),
	})
}

func serverCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Run autograd server",
		RunE: func(cmd *cobra.Command, args []string) error {
			redisPool := config.NewRedisPool(config.RedisWorkerHost())
			gormDB := db.MustSQLite()
			broker := worker.NewBroker(redisPool)
			localStorage := fs.NewLocalStorage()

			userRepo := repository.NewUserRepository(gormDB)
			userUsecase := usecase.NewUserUsecase(userRepo)
			submissionRepo := repository.NewSubmissionRepo(gormDB)
			assignmentRepo := repository.NewAssignmentRepository(gormDB)

			assignmentUsecase := usecase.NewAssignmentUsecase(assignmentRepo, submissionRepo)
			submissionUsecase := usecase.NewSubmissionUsecase(submissionRepo, usecase.SubmissionUsecaseWithBroker(broker))
			mediaUsecase := usecase.NewMediaUsecase(config.FileUploadPath(), localStorage)
			graderUsecase := usecase.NewGraderUsecase(submissionUsecase, assignmentUsecase)

			service := mustInitService()

			server := httpsvc.NewServer(
				config.Port(),
				config.FileUploadPath(),
				httpsvc.WithUserUsecase(userUsecase),
				httpsvc.WithAssignmentUsecase(assignmentUsecase),
				httpsvc.WithSubmissionUsecase(submissionUsecase),
				httpsvc.WithMediaUsecase(mediaUsecase),
				httpsvc.WithGormDB(gormDB),
				httpsvc.WithService(service),
				httpsvc.WithJWTKey(config.JWTKey()),
			)

			wrk := worker.NewWorker(redisPool, worker.WithGrader(graderUsecase))

			go func() {
				logrus.Info("run server")
				server.Run()
			}()

			go func() {
				logrus.Info("run worker")
				wrk.Start()
			}()

			// Wait for a signal to quit:
			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
			<-signalChan

			logrus.Info("stopping server")
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			server.Stop(ctx)

			logrus.Info("stopping worker")
			time.AfterFunc(time.Second*30, func() {
				os.Exit(1)
			})
			wrk.Stop()
			logrus.Info("worker stopped")

			return nil
		},
	}
}

func adminCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin",
		Short: "Admin command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(runCreateAdminUser())
	cmd.AddCommand(cmdAdminUser())

	return cmd
}

func runCreateAdminUser() *cobra.Command {
	service := mustInitService()

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create admin user",
	}

	req := user_management_cmd.CreateAdminUserRequest{}
	cmd.Flags().StringVar(&req.Name, "name", "", "admin name")
	cmd.Flags().StringVar(&req.Email, "email", "", "admin email")
	cmd.Flags().StringVar(&req.Password, "password", "", "admin password")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("email")
	cmd.MarkFlagRequired("password")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := service.InternalCreateAdminUser(cmd.Context(), req)
		if err != nil {
			return err
		}

		fmt.Println("Admin user created with id:", res.String())
		return nil
	}

	return cmd
}

func cmdAdminUser() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Admin user command",
	}

	cmd.AddCommand(runAdminCreateUser())

	return cmd
}

func initServiceClient() autogradv1connect.AutogradServiceClient {
	httpclient := http.DefaultClient
	setHeader := func(uf connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {
			ar.Header().Set("Authorization", "Bearer "+config.AutogradAuthToken())
			res, err := uf(ctx, ar)
			if err != nil {
				return nil, err
			}

			return res, nil
		}
	}
	interceptor := connect.WithInterceptors(connect.UnaryInterceptorFunc(setHeader))
	client := autogradv1connect.NewAutogradServiceClient(httpclient, config.AutogradServerURL(), interceptor)

	return client
}

func runAdminCreateUser() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
	}

	req := &autogradv1.CreateManagedUserRequest{}
	cmd.Flags().StringVar(&req.Name, "name", "", "user name")
	cmd.Flags().StringVar(&req.Email, "email", "", "user email")
	cmd.Flags().StringVar(&req.Role, "role", "", "user role")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("email")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("role")

	client := initServiceClient()

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := client.CreateManagedUser(cmd.Context(), &connect.Request[autogradv1.CreateManagedUserRequest]{
			Msg: req,
		})
		if err != nil {
			fmt.Println("CreateUser failed:", err)
			return err
		}

		fmt.Println("User created with id:", res.Msg.GetId())
		return nil
	}

	return cmd
}

func loginCmd() *cobra.Command {
	service := mustInitService()

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login user",
	}

	req := auth_cmd.InternalLoginRequest{}
	cmd.Flags().StringVar(&req.Email, "email", "", "email")
	cmd.Flags().StringVar(&req.Password, "password", "", "password")

	cmd.MarkFlagRequired("email")
	cmd.MarkFlagRequired("password")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		_, token, err := service.InternalLogin(cmd.Context(), req)
		if err != nil {
			fmt.Println("Login failed:", err)
			return err
		}

		fmt.Printf("User logged in with token:\n\n%s\n", token)
		return nil
	}

	return cmd
}
