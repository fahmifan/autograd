package student_courses_cmdquery

import (
	"context"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/fahmifan/autograd/pkg/core/student_courses"
	"github.com/fahmifan/autograd/pkg/dbconn"
	"github.com/fahmifan/autograd/pkg/logs"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	"github.com/fahmifan/autograd/pkg/xsqlc"
)

type StudentCoursesQuery struct {
	*core.Ctx
}

func (query *StudentCoursesQuery) FindAllStudentEnrolledCourses(
	ctx context.Context,
	req *connect.Request[autogradv1.FindAllStudentEnrolledCoursesRequest],
) (
	*connect.Response[autogradv1.FindAllStudentEnrolledCoursesResponse],
	error,
) {
	authUser, _ := auth.GetUserFromCtx(ctx)
	if !authUser.Role.Can(auth.ViewStudentEnrolledCourses) {
		return nil, core.ErrPermissionDenied
	}

	sqlcQuery, err := dbconn.NewSqlcFromGorm(query.GormDB)
	if err != nil {
		logs.ErrCtx(ctx, err, "StudentCoursesQuery: FindAllStudentEnrolledCourses: NewSqlcFromGorm")
		return nil, core.ErrInternalServer
	}

	pagin := core.PaginationRequestFromProto(req.Msg.GetPaginationRequest())

	rows, err := sqlcQuery.FindAllStudentEnrolledCourses(ctx, xsqlc.FindAllStudentEnrolledCoursesParams{
		UserID:     authUser.UserID.String(),
		UserType:   string(student_courses.CourseUserTypeStudent),
		PageLimit:  pagin.Limit,
		PageOffset: pagin.Offset(),
	})
	if err != nil {
		logs.ErrCtx(ctx, err, "StudentCoursesQuery: FindAllStudentEnrolledCourses: FindAllStudentEnrolledCourses")
		return nil, core.ErrInternalServer
	}

	total, err := sqlcQuery.CountAllStudentEnrolledCourses(ctx, xsqlc.CountAllStudentEnrolledCoursesParams{
		UserID:   authUser.UserID.String(),
		UserType: string(student_courses.CourseUserTypeStudent),
	})
	if err != nil && !core.IsErrDBNotFound(err) {
		logs.ErrCtx(ctx, err, "StudentCoursesQuery: FindAllStudentEnrolledCourses: CountAllStudentEnrolledCourses")
		return nil, core.ErrInternalServer
	}

	protoCourses := make([]*autogradv1.FindAllStudentEnrolledCoursesResponse_Course, len(rows))
	for i := range rows {
		protoCourses[i] = &autogradv1.FindAllStudentEnrolledCoursesResponse_Course{
			Id:          rows[i].ID,
			Name:        rows[i].Name,
			Description: rows[i].Description,
		}
	}

	return &connect.Response[autogradv1.FindAllStudentEnrolledCoursesResponse]{
		Msg: &autogradv1.FindAllStudentEnrolledCoursesResponse{
			PaginationMetadata: pagin.WithTotal(int32(total)).ProtoPagination(),
			Courses:            protoCourses,
		},
	}, nil
}

func (query *StudentCoursesQuery) FindStudentCourseDetail(
	ctx context.Context,
	req *connect.Request[autogradv1.FindByIDRequest],
) (
	*connect.Response[autogradv1.FindStudentCourseDetailResponse],
	error,
) {
	sqlcQuery, err := dbconn.NewSqlcFromGorm(query.GormDB)
	if err != nil {
		logs.ErrCtx(ctx, err, "StudentCoursesQuery: FindStudentCourseDetail: NewSqlcFromGorm")
		return nil, core.ErrInternalServer
	}

	course, err := sqlcQuery.FindCourseByID(ctx, req.Msg.GetId())
	if err != nil {
		if core.IsErrDBNotFound(err) {
			return nil, core.ErrNotFound
		}
		logs.ErrCtx(ctx, err, "StudentCoursesQuery: FindStudentCourseDetail: FindCourseByID")
		return nil, core.ErrInternalServer
	}

	return &connect.Response[autogradv1.FindStudentCourseDetailResponse]{
		Msg: &autogradv1.FindStudentCourseDetailResponse{
			Course: &autogradv1.FindStudentCourseDetailResponse_Course{
				Id:          course.ID,
				Name:        course.Name,
				Description: course.Description,
			},
		},
	}, nil
}
