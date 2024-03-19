package admin_courses_cmd

import (
	"context"
	"errors"
	"time"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/admin_courses"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/fahmifan/autograd/pkg/dbconn"
	"github.com/fahmifan/autograd/pkg/logs"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	"github.com/fahmifan/autograd/pkg/xsqlc"
	"github.com/fahmifan/ulids"
)

type AdminCoursesCmd struct {
	*core.Ctx
}

func (cmd *AdminCoursesCmd) CreateAdminCourse(
	ctx context.Context,
	req *connect.Request[autogradv1.CreateAdminCourseRequest],
) (*connect.Response[autogradv1.CreatedResponse], error) {
	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return nil, core.ErrUnauthenticated
	}

	if !authUser.Role.Can(auth.CreateCourse) {
		return nil, core.ErrPermissionDenied
	}

	adminReader := admin_courses.AdminReader{}
	courseWriter := admin_courses.CourseWriter{}

	admin, err := adminReader.FindAdminByID(ctx, cmd.SqlDB, authUser.UserID)
	if err != nil {
		if core.IsDBNotFoundErr(err) {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("admin not found"))
		}

		logs.ErrCtx(ctx, err, "AdminCoursesCmd: CreateAdminCourses: SaveCourse")
		return nil, core.ErrInternalServer
	}

	course, err := admin_courses.CreateCourse(admin_courses.CreateCourseRequest{
		Admin:       admin,
		NewID:       ulids.New(),
		Name:        req.Msg.Name,
		Description: req.Msg.Description,
	})
	if err != nil {
		return nil, err
	}

	err = courseWriter.Save(ctx, cmd.SqlDB, &course)
	if err != nil {
		logs.ErrCtx(ctx, err, "AdminCoursesCmd: CreateAdminCourses: SaveCourse")
		return nil, core.ErrInternalServer
	}

	return &connect.Response[autogradv1.CreatedResponse]{
		Msg: &autogradv1.CreatedResponse{
			Id: course.ID.String(),
		},
	}, nil
}

func (cmd *AdminCoursesCmd) UpdateAdminCourse(
	ctx context.Context,
	req *connect.Request[autogradv1.UpdateAdminCourseRequest],
) (*connect.Response[autogradv1.Empty], error) {
	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return nil, core.ErrUnauthenticated
	}

	if !authUser.Role.Can(auth.UpdateCourse) {
		return nil, core.ErrPermissionDenied
	}

	id, err := ulids.Parse(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	courseReader := admin_courses.CourseReader{}
	courseWriter := admin_courses.CourseWriter{}

	err = dbconn.SqlcTransaction(ctx, cmd.SqlDB, func(d xsqlc.DBTX) error {
		course, err := courseReader.FindCourseByID(ctx, cmd.SqlDB, id)
		if err != nil {
			if core.IsDBNotFoundErr(err) {
				return connect.NewError(connect.CodeNotFound, errors.New("course not found"))
			}

			logs.ErrCtx(ctx, err, "AdminCoursesCmd: UpdateAdminCourse: FindCourseByID")
			return core.ErrInternalServer
		}

		course, err = course.Update(admin_courses.UpdateCourseRequest{
			Now:         time.Now(),
			Name:        req.Msg.Name,
			Description: req.Msg.Description,
			IsActive:    req.Msg.Active,
		})
		if err != nil {
			return err
		}

		err = courseWriter.Save(ctx, cmd.SqlDB, &course)
		if err != nil {
			logs.ErrCtx(ctx, err, "AdminCoursesCmd: UpdateAdminCourse: SaveCourse")
			return core.ErrInternalServer
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &connect.Response[autogradv1.Empty]{}, nil
}
