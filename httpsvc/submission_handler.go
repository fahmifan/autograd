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
	submissionReq := &submissionRequest{}
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

	return c.JSON(http.StatusOK, submissionModelToResponse(submission))
}

func (s *Server) handleUpload(c echo.Context) error {
	uploadReq := &uploadRequest{}
	err := c.Bind(uploadReq)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	upload := &model.Upload{}
	err = smapping.FillStruct(upload, smapping.MapFields(uploadReq))
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	err = s.submissionUsecase.Upload(c.Request().Context(), upload)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	uploadRes := &uploadResponse{}
	err = smapping.FillStruct(uploadRes, smapping.MapFields(upload))
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, uploadRes)
}

func (s *Server) handleGetAssignmentSubmission(c echo.Context) error {
	assignmentID := utils.StringToInt64(c.Param("assignmentID"))
	cursor := getCursorFromContext(c)
	submissions, count, err := s.submissionUsecase.FindAllByAssignmentID(c.Request().Context(), cursor, assignmentID)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	submissionRes := newSubmissionResponses(submissions)

	return c.JSON(http.StatusOK, newCursorResponse(cursor, submissionRes, count))
}
