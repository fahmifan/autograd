package httpsvc

import (
	"net/http"

	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) handleCreateSubmission(c echo.Context) error {
	req := &submissionReq{}
	err := c.Bind(req)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	submission := submissionCreateReqToModel(req)
	err = s.submissionUsecase.Create(c.Request().Context(), submission)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, submissionModelToRes(submission))
}

func (s *Server) handleDeleteSubmission(c echo.Context) error {
	id := c.Param("id")
	submission, err := s.submissionUsecase.DeleteByID(c.Request().Context(), id)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, submissionModelToRes(submission))
}

func (s *Server) handleGetSubmission(c echo.Context) error {
	id := c.Param("id")
	submission, err := s.submissionUsecase.FindByID(c.Request().Context(), id)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	user := getUserFromCtx(c)
	if user.Role == auth.RoleStudent && !submission.IsOwnedBy(*user) {
		return responseError(c, ErrUnauthorized)
	}

	return c.JSON(http.StatusOK, submissionModelToRes(submission))
}

func (s *Server) handleUpdateSubmission(c echo.Context) error {
	submissionReq := &submissionReq{}
	err := c.Bind(submissionReq)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	submission := submissionUpdateReqToModel(submissionReq)
	ctx := c.Request().Context()
	user := getUserFromCtx(c)

	oldSubm, err := s.submissionUsecase.FindByID(ctx, submission.ID)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	if !oldSubm.IsOwnedBy(*user) {
		return responseError(c, ErrUnauthorized)
	}

	err = s.submissionUsecase.Update(c.Request().Context(), submission)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, submissionModelToRes(submission))
}
