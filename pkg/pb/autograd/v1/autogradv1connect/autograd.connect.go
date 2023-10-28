// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: autograd/v1/autograd.proto

package autogradv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion0_1_0

const (
	// AutogradServiceName is the fully-qualified name of the AutogradService service.
	AutogradServiceName = "autograd.v1.AutogradService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// AutogradServicePingProcedure is the fully-qualified name of the AutogradService's Ping RPC.
	AutogradServicePingProcedure = "/autograd.v1.AutogradService/Ping"
	// AutogradServiceCreateUserProcedure is the fully-qualified name of the AutogradService's
	// CreateUser RPC.
	AutogradServiceCreateUserProcedure = "/autograd.v1.AutogradService/CreateUser"
	// AutogradServiceFindAssignmentProcedure is the fully-qualified name of the AutogradService's
	// FindAssignment RPC.
	AutogradServiceFindAssignmentProcedure = "/autograd.v1.AutogradService/FindAssignment"
	// AutogradServiceFindSubmissionProcedure is the fully-qualified name of the AutogradService's
	// FindSubmission RPC.
	AutogradServiceFindSubmissionProcedure = "/autograd.v1.AutogradService/FindSubmission"
	// AutogradServiceCreateAssignmentProcedure is the fully-qualified name of the AutogradService's
	// CreateAssignment RPC.
	AutogradServiceCreateAssignmentProcedure = "/autograd.v1.AutogradService/CreateAssignment"
	// AutogradServiceUpdateAssignmentProcedure is the fully-qualified name of the AutogradService's
	// UpdateAssignment RPC.
	AutogradServiceUpdateAssignmentProcedure = "/autograd.v1.AutogradService/UpdateAssignment"
	// AutogradServiceDeleteAssignmentProcedure is the fully-qualified name of the AutogradService's
	// DeleteAssignment RPC.
	AutogradServiceDeleteAssignmentProcedure = "/autograd.v1.AutogradService/DeleteAssignment"
	// AutogradServiceCreateSubmissionProcedure is the fully-qualified name of the AutogradService's
	// CreateSubmission RPC.
	AutogradServiceCreateSubmissionProcedure = "/autograd.v1.AutogradService/CreateSubmission"
	// AutogradServiceUpdateSubmissionProcedure is the fully-qualified name of the AutogradService's
	// UpdateSubmission RPC.
	AutogradServiceUpdateSubmissionProcedure = "/autograd.v1.AutogradService/UpdateSubmission"
	// AutogradServiceDeleteSubmissionProcedure is the fully-qualified name of the AutogradService's
	// DeleteSubmission RPC.
	AutogradServiceDeleteSubmissionProcedure = "/autograd.v1.AutogradService/DeleteSubmission"
)

// AutogradServiceClient is a client for the autograd.v1.AutogradService service.
type AutogradServiceClient interface {
	Ping(context.Context, *connect.Request[v1.Empty]) (*connect.Response[v1.PingResponse], error)
	CreateUser(context.Context, *connect.Request[v1.CreateUserRequest]) (*connect.Response[v1.CreatedResponse], error)
	// Assignment Submission
	// Assignment Queries
	FindAssignment(context.Context, *connect.Request[v1.FindByIDRequest]) (*connect.Response[v1.Assignment], error)
	FindSubmission(context.Context, *connect.Request[v1.FindByIDRequest]) (*connect.Response[v1.Submission], error)
	// Assignment Mutations
	CreateAssignment(context.Context, *connect.Request[v1.CreateAssignmentRequest]) (*connect.Response[v1.CreatedResponse], error)
	UpdateAssignment(context.Context, *connect.Request[v1.UpdateAssignmentRequest]) (*connect.Response[v1.Empty], error)
	DeleteAssignment(context.Context, *connect.Request[v1.DeleteByIDRequest]) (*connect.Response[v1.Empty], error)
	CreateSubmission(context.Context, *connect.Request[v1.CreateSubmissionRequest]) (*connect.Response[v1.CreatedResponse], error)
	UpdateSubmission(context.Context, *connect.Request[v1.UpdateSubmissionRequest]) (*connect.Response[v1.Empty], error)
	DeleteSubmission(context.Context, *connect.Request[v1.DeleteByIDRequest]) (*connect.Response[v1.Empty], error)
}

// NewAutogradServiceClient constructs a client for the autograd.v1.AutogradService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAutogradServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) AutogradServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &autogradServiceClient{
		ping: connect.NewClient[v1.Empty, v1.PingResponse](
			httpClient,
			baseURL+AutogradServicePingProcedure,
			opts...,
		),
		createUser: connect.NewClient[v1.CreateUserRequest, v1.CreatedResponse](
			httpClient,
			baseURL+AutogradServiceCreateUserProcedure,
			opts...,
		),
		findAssignment: connect.NewClient[v1.FindByIDRequest, v1.Assignment](
			httpClient,
			baseURL+AutogradServiceFindAssignmentProcedure,
			opts...,
		),
		findSubmission: connect.NewClient[v1.FindByIDRequest, v1.Submission](
			httpClient,
			baseURL+AutogradServiceFindSubmissionProcedure,
			opts...,
		),
		createAssignment: connect.NewClient[v1.CreateAssignmentRequest, v1.CreatedResponse](
			httpClient,
			baseURL+AutogradServiceCreateAssignmentProcedure,
			opts...,
		),
		updateAssignment: connect.NewClient[v1.UpdateAssignmentRequest, v1.Empty](
			httpClient,
			baseURL+AutogradServiceUpdateAssignmentProcedure,
			opts...,
		),
		deleteAssignment: connect.NewClient[v1.DeleteByIDRequest, v1.Empty](
			httpClient,
			baseURL+AutogradServiceDeleteAssignmentProcedure,
			opts...,
		),
		createSubmission: connect.NewClient[v1.CreateSubmissionRequest, v1.CreatedResponse](
			httpClient,
			baseURL+AutogradServiceCreateSubmissionProcedure,
			opts...,
		),
		updateSubmission: connect.NewClient[v1.UpdateSubmissionRequest, v1.Empty](
			httpClient,
			baseURL+AutogradServiceUpdateSubmissionProcedure,
			opts...,
		),
		deleteSubmission: connect.NewClient[v1.DeleteByIDRequest, v1.Empty](
			httpClient,
			baseURL+AutogradServiceDeleteSubmissionProcedure,
			opts...,
		),
	}
}

// autogradServiceClient implements AutogradServiceClient.
type autogradServiceClient struct {
	ping             *connect.Client[v1.Empty, v1.PingResponse]
	createUser       *connect.Client[v1.CreateUserRequest, v1.CreatedResponse]
	findAssignment   *connect.Client[v1.FindByIDRequest, v1.Assignment]
	findSubmission   *connect.Client[v1.FindByIDRequest, v1.Submission]
	createAssignment *connect.Client[v1.CreateAssignmentRequest, v1.CreatedResponse]
	updateAssignment *connect.Client[v1.UpdateAssignmentRequest, v1.Empty]
	deleteAssignment *connect.Client[v1.DeleteByIDRequest, v1.Empty]
	createSubmission *connect.Client[v1.CreateSubmissionRequest, v1.CreatedResponse]
	updateSubmission *connect.Client[v1.UpdateSubmissionRequest, v1.Empty]
	deleteSubmission *connect.Client[v1.DeleteByIDRequest, v1.Empty]
}

// Ping calls autograd.v1.AutogradService.Ping.
func (c *autogradServiceClient) Ping(ctx context.Context, req *connect.Request[v1.Empty]) (*connect.Response[v1.PingResponse], error) {
	return c.ping.CallUnary(ctx, req)
}

// CreateUser calls autograd.v1.AutogradService.CreateUser.
func (c *autogradServiceClient) CreateUser(ctx context.Context, req *connect.Request[v1.CreateUserRequest]) (*connect.Response[v1.CreatedResponse], error) {
	return c.createUser.CallUnary(ctx, req)
}

// FindAssignment calls autograd.v1.AutogradService.FindAssignment.
func (c *autogradServiceClient) FindAssignment(ctx context.Context, req *connect.Request[v1.FindByIDRequest]) (*connect.Response[v1.Assignment], error) {
	return c.findAssignment.CallUnary(ctx, req)
}

// FindSubmission calls autograd.v1.AutogradService.FindSubmission.
func (c *autogradServiceClient) FindSubmission(ctx context.Context, req *connect.Request[v1.FindByIDRequest]) (*connect.Response[v1.Submission], error) {
	return c.findSubmission.CallUnary(ctx, req)
}

// CreateAssignment calls autograd.v1.AutogradService.CreateAssignment.
func (c *autogradServiceClient) CreateAssignment(ctx context.Context, req *connect.Request[v1.CreateAssignmentRequest]) (*connect.Response[v1.CreatedResponse], error) {
	return c.createAssignment.CallUnary(ctx, req)
}

// UpdateAssignment calls autograd.v1.AutogradService.UpdateAssignment.
func (c *autogradServiceClient) UpdateAssignment(ctx context.Context, req *connect.Request[v1.UpdateAssignmentRequest]) (*connect.Response[v1.Empty], error) {
	return c.updateAssignment.CallUnary(ctx, req)
}

// DeleteAssignment calls autograd.v1.AutogradService.DeleteAssignment.
func (c *autogradServiceClient) DeleteAssignment(ctx context.Context, req *connect.Request[v1.DeleteByIDRequest]) (*connect.Response[v1.Empty], error) {
	return c.deleteAssignment.CallUnary(ctx, req)
}

// CreateSubmission calls autograd.v1.AutogradService.CreateSubmission.
func (c *autogradServiceClient) CreateSubmission(ctx context.Context, req *connect.Request[v1.CreateSubmissionRequest]) (*connect.Response[v1.CreatedResponse], error) {
	return c.createSubmission.CallUnary(ctx, req)
}

// UpdateSubmission calls autograd.v1.AutogradService.UpdateSubmission.
func (c *autogradServiceClient) UpdateSubmission(ctx context.Context, req *connect.Request[v1.UpdateSubmissionRequest]) (*connect.Response[v1.Empty], error) {
	return c.updateSubmission.CallUnary(ctx, req)
}

// DeleteSubmission calls autograd.v1.AutogradService.DeleteSubmission.
func (c *autogradServiceClient) DeleteSubmission(ctx context.Context, req *connect.Request[v1.DeleteByIDRequest]) (*connect.Response[v1.Empty], error) {
	return c.deleteSubmission.CallUnary(ctx, req)
}

// AutogradServiceHandler is an implementation of the autograd.v1.AutogradService service.
type AutogradServiceHandler interface {
	Ping(context.Context, *connect.Request[v1.Empty]) (*connect.Response[v1.PingResponse], error)
	CreateUser(context.Context, *connect.Request[v1.CreateUserRequest]) (*connect.Response[v1.CreatedResponse], error)
	// Assignment Submission
	// Assignment Queries
	FindAssignment(context.Context, *connect.Request[v1.FindByIDRequest]) (*connect.Response[v1.Assignment], error)
	FindSubmission(context.Context, *connect.Request[v1.FindByIDRequest]) (*connect.Response[v1.Submission], error)
	// Assignment Mutations
	CreateAssignment(context.Context, *connect.Request[v1.CreateAssignmentRequest]) (*connect.Response[v1.CreatedResponse], error)
	UpdateAssignment(context.Context, *connect.Request[v1.UpdateAssignmentRequest]) (*connect.Response[v1.Empty], error)
	DeleteAssignment(context.Context, *connect.Request[v1.DeleteByIDRequest]) (*connect.Response[v1.Empty], error)
	CreateSubmission(context.Context, *connect.Request[v1.CreateSubmissionRequest]) (*connect.Response[v1.CreatedResponse], error)
	UpdateSubmission(context.Context, *connect.Request[v1.UpdateSubmissionRequest]) (*connect.Response[v1.Empty], error)
	DeleteSubmission(context.Context, *connect.Request[v1.DeleteByIDRequest]) (*connect.Response[v1.Empty], error)
}

// NewAutogradServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAutogradServiceHandler(svc AutogradServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	autogradServicePingHandler := connect.NewUnaryHandler(
		AutogradServicePingProcedure,
		svc.Ping,
		opts...,
	)
	autogradServiceCreateUserHandler := connect.NewUnaryHandler(
		AutogradServiceCreateUserProcedure,
		svc.CreateUser,
		opts...,
	)
	autogradServiceFindAssignmentHandler := connect.NewUnaryHandler(
		AutogradServiceFindAssignmentProcedure,
		svc.FindAssignment,
		opts...,
	)
	autogradServiceFindSubmissionHandler := connect.NewUnaryHandler(
		AutogradServiceFindSubmissionProcedure,
		svc.FindSubmission,
		opts...,
	)
	autogradServiceCreateAssignmentHandler := connect.NewUnaryHandler(
		AutogradServiceCreateAssignmentProcedure,
		svc.CreateAssignment,
		opts...,
	)
	autogradServiceUpdateAssignmentHandler := connect.NewUnaryHandler(
		AutogradServiceUpdateAssignmentProcedure,
		svc.UpdateAssignment,
		opts...,
	)
	autogradServiceDeleteAssignmentHandler := connect.NewUnaryHandler(
		AutogradServiceDeleteAssignmentProcedure,
		svc.DeleteAssignment,
		opts...,
	)
	autogradServiceCreateSubmissionHandler := connect.NewUnaryHandler(
		AutogradServiceCreateSubmissionProcedure,
		svc.CreateSubmission,
		opts...,
	)
	autogradServiceUpdateSubmissionHandler := connect.NewUnaryHandler(
		AutogradServiceUpdateSubmissionProcedure,
		svc.UpdateSubmission,
		opts...,
	)
	autogradServiceDeleteSubmissionHandler := connect.NewUnaryHandler(
		AutogradServiceDeleteSubmissionProcedure,
		svc.DeleteSubmission,
		opts...,
	)
	return "/autograd.v1.AutogradService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case AutogradServicePingProcedure:
			autogradServicePingHandler.ServeHTTP(w, r)
		case AutogradServiceCreateUserProcedure:
			autogradServiceCreateUserHandler.ServeHTTP(w, r)
		case AutogradServiceFindAssignmentProcedure:
			autogradServiceFindAssignmentHandler.ServeHTTP(w, r)
		case AutogradServiceFindSubmissionProcedure:
			autogradServiceFindSubmissionHandler.ServeHTTP(w, r)
		case AutogradServiceCreateAssignmentProcedure:
			autogradServiceCreateAssignmentHandler.ServeHTTP(w, r)
		case AutogradServiceUpdateAssignmentProcedure:
			autogradServiceUpdateAssignmentHandler.ServeHTTP(w, r)
		case AutogradServiceDeleteAssignmentProcedure:
			autogradServiceDeleteAssignmentHandler.ServeHTTP(w, r)
		case AutogradServiceCreateSubmissionProcedure:
			autogradServiceCreateSubmissionHandler.ServeHTTP(w, r)
		case AutogradServiceUpdateSubmissionProcedure:
			autogradServiceUpdateSubmissionHandler.ServeHTTP(w, r)
		case AutogradServiceDeleteSubmissionProcedure:
			autogradServiceDeleteSubmissionHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedAutogradServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAutogradServiceHandler struct{}

func (UnimplementedAutogradServiceHandler) Ping(context.Context, *connect.Request[v1.Empty]) (*connect.Response[v1.PingResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("autograd.v1.AutogradService.Ping is not implemented"))
}

func (UnimplementedAutogradServiceHandler) CreateUser(context.Context, *connect.Request[v1.CreateUserRequest]) (*connect.Response[v1.CreatedResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("autograd.v1.AutogradService.CreateUser is not implemented"))
}

func (UnimplementedAutogradServiceHandler) FindAssignment(context.Context, *connect.Request[v1.FindByIDRequest]) (*connect.Response[v1.Assignment], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("autograd.v1.AutogradService.FindAssignment is not implemented"))
}

func (UnimplementedAutogradServiceHandler) FindSubmission(context.Context, *connect.Request[v1.FindByIDRequest]) (*connect.Response[v1.Submission], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("autograd.v1.AutogradService.FindSubmission is not implemented"))
}

func (UnimplementedAutogradServiceHandler) CreateAssignment(context.Context, *connect.Request[v1.CreateAssignmentRequest]) (*connect.Response[v1.CreatedResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("autograd.v1.AutogradService.CreateAssignment is not implemented"))
}

func (UnimplementedAutogradServiceHandler) UpdateAssignment(context.Context, *connect.Request[v1.UpdateAssignmentRequest]) (*connect.Response[v1.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("autograd.v1.AutogradService.UpdateAssignment is not implemented"))
}

func (UnimplementedAutogradServiceHandler) DeleteAssignment(context.Context, *connect.Request[v1.DeleteByIDRequest]) (*connect.Response[v1.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("autograd.v1.AutogradService.DeleteAssignment is not implemented"))
}

func (UnimplementedAutogradServiceHandler) CreateSubmission(context.Context, *connect.Request[v1.CreateSubmissionRequest]) (*connect.Response[v1.CreatedResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("autograd.v1.AutogradService.CreateSubmission is not implemented"))
}

func (UnimplementedAutogradServiceHandler) UpdateSubmission(context.Context, *connect.Request[v1.UpdateSubmissionRequest]) (*connect.Response[v1.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("autograd.v1.AutogradService.UpdateSubmission is not implemented"))
}

func (UnimplementedAutogradServiceHandler) DeleteSubmission(context.Context, *connect.Request[v1.DeleteByIDRequest]) (*connect.Response[v1.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("autograd.v1.AutogradService.DeleteSubmission is not implemented"))
}
