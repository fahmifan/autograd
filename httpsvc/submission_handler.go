package httpsvc

import (
	"net/http"

	"github.com/fahmifan/autograd/utils"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) handleCreateSubmission(c echo.Context) error {
	submissionReq := &submissionReq{}
	err := c.Bind(submissionReq)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	submission := submissionCreateReqToModel(submissionReq)
	err = s.submissionUsecase.Create(c.Request().Context(), submission)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, submissionModelToRes(submission))
}

func (s *Server) handleDeleteSubmission(c echo.Context) error {
	id := utils.StringToInt64(c.Param("ID"))
	submission, err := s.submissionUsecase.DeleteByID(c.Request().Context(), id)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, submissionModelToRes(submission))
}

func (s *Server) handleGetSubmission(c echo.Context) error {
	id := utils.StringToInt64(c.Param("ID"))
	submission, err := s.submissionUsecase.FindByID(c.Request().Context(), id)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
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
	err = s.submissionUsecase.Update(c.Request().Context(), submission)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, submissionModelToRes(submission))
}
