package admin_courses_query

import (
	"context"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/fahmifan/autograd/pkg/logs"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	"github.com/fahmifan/autograd/pkg/xsqlc"
)

type AdminCoursesQuery struct {
	*core.Ctx
}

func (query *AdminCoursesQuery) FindAllAdminCourses(
	ctx context.Context,
	req *connect.Request[autogradv1.FindAllPaginationRequest],
) (
	*connect.Response[autogradv1.FindAllAdminCoursesResponse],
	error,
) {
	authUser, ok := auth.GetUserFromCtx(ctx)
	if !ok {
		return nil, core.ErrUnauthenticated
	}

	if !authUser.Role.Can(auth.ViewCourse) {
		return nil, core.ErrPermissionDenied
	}

	pagin := core.PaginationRequestFromProto(req.Msg.GetPaginationRequest())

	sqlcQuery := xsqlc.New(query.SqlDB)
	courses, err := sqlcQuery.FindAllCoursesByUser(ctx, xsqlc.FindAllCoursesByUserParams{
		UserID:     authUser.UserID.String(),
		PageLimit:  pagin.Limit,
		PageOffset: pagin.Offset(),
	})
	if err != nil {
		logs.ErrCtx(ctx, err, "AdminCoursesQuery: FindAllCoursesByUser: FindAllCoursesByUser")
		return nil, core.ErrInternalServer
	}

	totalCount, err := sqlcQuery.CountAllCoursesByUser(ctx, authUser.UserID.String())
	if err != nil && !core.IsDBNotFoundErr(err) {
		logs.ErrCtx(ctx, err, "AdminCoursesQuery: FindAllCoursesByUser: CountAllAdminCourses")
		return nil, core.ErrInternalServer
	}

	protoCourses := make([]*autogradv1.FindAllAdminCoursesResponse_Course, len(courses))
	for i := range courses {
		protoCourses[i] = &autogradv1.FindAllAdminCoursesResponse_Course{
			Id:          courses[i].ID,
			Name:        courses[i].Name,
			Description: courses[i].Description,
		}
	}

	res := &connect.Response[autogradv1.FindAllAdminCoursesResponse]{
		Msg: &autogradv1.FindAllAdminCoursesResponse{
			Courses:            protoCourses,
			PaginationMetadata: pagin.WithTotal(int32(totalCount)).ProtoPagination(),
		},
	}

	return res, nil
}
