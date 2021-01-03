package httpsvc

import (
	"net/http"

	"github.com/mashingan/smapping"
	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/utils"

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

	submission := &model.Submission{}
	err = smapping.FillStruct(submission, smapping.MapFields(submissionReq))
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	err = s.submissionUsecase.Create(c.Request().Context(), submission)
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

func (s *Server) handleGetAssignmentSubmission(c echo.Context) error {
	assignmentID := utils.StringToInt64(c.Param("ID"))
	cursor := getCursorFromContext(c)
	submissions, count, err := s.submissionUsecase.FindAllByAssignmentID(c.Request().Context(), cursor, assignmentID)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	submissionRes := newSubmissionResponses(submissions)

	return c.JSON(http.StatusOK, newCursorRes(cursor, submissionRes, count))
}

func (s *Server) handleUpload(c echo.Context) error {
	uploadReq := &uploadReq{}
	err := c.Bind(uploadReq)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	fileURL, err := s.submissionUsecase.Upload(c.Request().Context(), uploadReq.SourceCode)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	uploadRes := &uploadRes{FileURL: fileURL}

	return c.JSON(http.StatusOK, uploadRes)
}
