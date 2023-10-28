package user_management_query

import (
	"context"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/user_management"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
)

type UserManagementQuery struct {
	*core.Ctx
}

func (query *UserManagementQuery) FindAllManagedUsers(
	ctx context.Context,
	req *connect.Request[autogradv1.FindAllManagedUsersRequest],
) (*connect.Response[autogradv1.FindAllManagedUsersResponse], error) {
	res, err := user_management.ManagedUserReader{}.FindAll(ctx, query.GormDB, user_management.FindAllManagedUsersRequest{
		Page:  req.Msg.GetPage(),
		Limit: req.Msg.GetLimit(),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &connect.Response[autogradv1.FindAllManagedUsersResponse]{
		Msg: &autogradv1.FindAllManagedUsersResponse{
			ManagedUsers: toManagedUserProtos(res.Users),
			PaginationMetadata: &autogradv1.PaginationMetadata{
				Total: res.Count,
				Page:  req.Msg.GetPage(),
				Limit: req.Msg.GetLimit(),
			},
		},
	}, nil
}

func toManagedUserProtos(users []user_management.ManagedUser) []*autogradv1.ManagedUser {
	var userProtos []*autogradv1.ManagedUser
	for _, user := range users {
		userProtos = append(userProtos, toManagedUserProto(user))
	}
	return userProtos
}

func toManagedUserProto(user user_management.ManagedUser) *autogradv1.ManagedUser {
	return &autogradv1.ManagedUser{
		Id:                user.ID.String(),
		Name:              user.Name,
		Email:             user.Email,
		Role:              string(user.Role),
		TimestampMetadata: user.TimestampMetadata.ProtoTimestampMetadata(),
	}
}
