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
	AutogradServiceName = "shopper.v1.AutogradService"
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
	AutogradServicePingProcedure = "/shopper.v1.AutogradService/Ping"
	// AutogradServiceCreateUserProcedure is the fully-qualified name of the AutogradService's
	// CreateUser RPC.
	AutogradServiceCreateUserProcedure = "/shopper.v1.AutogradService/CreateUser"
)

// AutogradServiceClient is a client for the shopper.v1.AutogradService service.
type AutogradServiceClient interface {
	Ping(context.Context, *connect.Request[v1.Empty]) (*connect.Response[v1.PingResponse], error)
	CreateUser(context.Context, *connect.Request[v1.CreateUserRequest]) (*connect.Response[v1.CreatedResponse], error)
}

// NewAutogradServiceClient constructs a client for the shopper.v1.AutogradService service. By
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
	}
}

// autogradServiceClient implements AutogradServiceClient.
type autogradServiceClient struct {
	ping       *connect.Client[v1.Empty, v1.PingResponse]
	createUser *connect.Client[v1.CreateUserRequest, v1.CreatedResponse]
}

// Ping calls shopper.v1.AutogradService.Ping.
func (c *autogradServiceClient) Ping(ctx context.Context, req *connect.Request[v1.Empty]) (*connect.Response[v1.PingResponse], error) {
	return c.ping.CallUnary(ctx, req)
}

// CreateUser calls shopper.v1.AutogradService.CreateUser.
func (c *autogradServiceClient) CreateUser(ctx context.Context, req *connect.Request[v1.CreateUserRequest]) (*connect.Response[v1.CreatedResponse], error) {
	return c.createUser.CallUnary(ctx, req)
}

// AutogradServiceHandler is an implementation of the shopper.v1.AutogradService service.
type AutogradServiceHandler interface {
	Ping(context.Context, *connect.Request[v1.Empty]) (*connect.Response[v1.PingResponse], error)
	CreateUser(context.Context, *connect.Request[v1.CreateUserRequest]) (*connect.Response[v1.CreatedResponse], error)
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
	return "/shopper.v1.AutogradService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case AutogradServicePingProcedure:
			autogradServicePingHandler.ServeHTTP(w, r)
		case AutogradServiceCreateUserProcedure:
			autogradServiceCreateUserHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedAutogradServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAutogradServiceHandler struct{}

func (UnimplementedAutogradServiceHandler) Ping(context.Context, *connect.Request[v1.Empty]) (*connect.Response[v1.PingResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("shopper.v1.AutogradService.Ping is not implemented"))
}

func (UnimplementedAutogradServiceHandler) CreateUser(context.Context, *connect.Request[v1.CreateUserRequest]) (*connect.Response[v1.CreatedResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("shopper.v1.AutogradService.CreateUser is not implemented"))
}