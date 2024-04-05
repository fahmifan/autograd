package admin_courses_query

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"github.com/coocood/freecache"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/fahmifan/autograd/pkg/dbconn"
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

func (query *AdminCoursesQuery) FindAdminCourseDetail(
	ctx context.Context,
	req *connect.Request[autogradv1.FindByIDRequest],
) (
	*connect.Response[autogradv1.FindAdminCourseDetailResponse],
	error,
) {
	sqlcQuery, err := dbconn.NewSqlcFromGorm(query.GormDB)
	if err != nil {
		logs.ErrCtx(ctx, err, "AdminCoursesQuery: FindAdminCourseDetail: NewSqlcFromGorm")
		return nil, core.ErrInternalServer
	}

	course, err := sqlcQuery.FindCourseByID(ctx, req.Msg.GetId())
	if err != nil {
		if core.IsDBNotFoundErr(err) {
			return nil, core.ErrNotFound
		}

		logs.ErrCtx(ctx, err, "AdminCoursesQuery: FindAdminCourseDetail: NewSqlcFromGorm")
		return nil, core.ErrInternalServer
	}

	return &connect.Response[autogradv1.FindAdminCourseDetailResponse]{
		Msg: &autogradv1.FindAdminCourseDetailResponse{
			Course: &autogradv1.FindAdminCourseDetailResponse_Course{
				Id:          course.ID,
				Name:        course.Name,
				Description: course.Description,
				Active:      course.IsActive,
			},
		},
	}, nil
}

func (query *AdminCoursesQuery) FindAllCourseStudents(
	ctx context.Context,
	req *connect.Request[autogradv1.FindAllCourseStudentsRequest],
) (
	*connect.Response[autogradv1.FindAllCourseStudentsResponse],
	error,
) {
	sqlcQuery, err := dbconn.NewSqlcFromGorm(query.GormDB)
	if err != nil {
		logs.ErrCtx(ctx, err, "AdminCoursesQuery: FindAllCourseStudents: NewSqlcFromGorm")
		return nil, core.ErrInternalServer
	}

	pagin := core.PaginationRequestFromProto(req.Msg.GetPaginationRequest())

	totalCacheKey := NewCacheKey("total", "course", "course_id", req.Msg.GetCourseId())
	total, err := GetOrCache(query.Ctx, totalCacheKey, 60, func() (int64, error) {
		total, err := sqlcQuery.CountAllStudentsByCourse(ctx, req.Msg.GetCourseId())
		if err != nil {
			logs.ErrCtx(ctx, err, "AdminCoursesQuery: FindAllCourseStudents: FindAllStudentsByCourse")
			return 0, core.ErrInternalServer
		}

		return total, nil
	})
	if err != nil {
		logs.ErrCtx(ctx, err, "AdminCoursesQuery: FindAllCourseStudents: GetOrCache")
		return nil, core.ErrInternalServer
	}

	students, err := sqlcQuery.FindAllStudentsByCourse(ctx, xsqlc.FindAllStudentsByCourseParams{
		CourseID:   req.Msg.GetCourseId(),
		PageOffset: pagin.Offset(),
		PageLimit:  pagin.Limit,
	})
	if err != nil {
		logs.ErrCtx(ctx, err, "AdminCoursesQuery: FindAllCourseStudents: FindAllStudentsByCourse")
		return nil, core.ErrInternalServer
	}

	resStudents := make([]*autogradv1.FindAllCourseStudentsResponse_Student, len(students))
	for i := range students {
		resStudents[i] = &autogradv1.FindAllCourseStudentsResponse_Student{
			Id:   students[i].ID,
			Name: students[i].Name,
		}
	}

	return &connect.Response[autogradv1.FindAllCourseStudentsResponse]{
		Msg: &autogradv1.FindAllCourseStudentsResponse{
			Students:           resStudents,
			PaginationMetadata: pagin.WithTotal(int32(total)).ProtoPagination(),
		},
	}, nil
}

type CacheKey string

func NewCacheKey(arg ...string) CacheKey {
	return CacheKey(strings.Join(arg, ":"))
}

func GetOrCache[T any](coreCtx *core.Ctx, key CacheKey, expireInSecond int, fetch func() (T, error)) (T, error) {
	val, err, _ := coreCtx.Flight.Do(string(key), func() (interface{}, error) {
		var tval T

		val, err := coreCtx.Cache.Get([]byte(key))
		if err == nil {
			// found
			err = json.Unmarshal(val, &tval)
			return tval, err
		}
		if !errors.Is(err, freecache.ErrNotFound) {
			return tval, err
		}

		tval, err = fetch()
		if err != nil {
			return tval, err
		}

		val, err = json.Marshal(tval)
		if err != nil {
			return tval, err
		}

		err = coreCtx.Cache.Set([]byte(key), val, expireInSecond)
		if err != nil {
			return tval, err
		}

		return tval, nil
	})
	defer coreCtx.Flight.Forget(string(key))

	var tval T
	if err != nil {
		return tval, nil
	}

	return val.(T), nil
}
