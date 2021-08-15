package httpsvc

import (
	"net/http"

	"github.com/fahmifan/autograd/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// CreateSubmission godoc
// @Summary create a user submission
// @ID CreateSubmission
// @Accept json
// @Produce json
// @Param user body SubmissionReq true "submission request"
// @Success 200 {object} SubmissionRes
// @Failure 400 {object} Error
// @Router /api/v1/submissions [post]
func (s *Server) handleCreateSubmission(c echo.Context) error {
	user := getUserFromCtx(c)
	req := &SubmissionReq{}
	err := c.Bind(req)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	submission := submissionCreateReqToModel(req)
	submission.SubmittedBy = user.ID
	err = s.submissionUsecase.Create(c.Request().Context(), submission)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, submissionModelToRes(submission))
}

// DeleteSubmission godoc
// @Summary delete a submission
// @ID DeleteSubmission
// @Accept json
// @Produce json
// @Param id path string true "submission id"
// @Success 200 {object} SubmissionRes
// @Failure 400,404 {object} Error
// @Router /api/v1/submissions/{id} [delete]
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
	user := getUserFromCtx(c)

	var submission *model.Submission
	var err error
	if user.Role.GrantedAny(model.ViewAnySubmissions) {
		submission, err = s.submissionUsecase.FindByID(c.Request().Context(), id)
	} else {
		submission, err = s.submissionUsecase.FindByIDAndSubmitter(c.Request().Context(), id, user.ID)
	}
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, submissionModelToRes(submission))
}

func (s *Server) handleUpdateSubmission(c echo.Context) error {
	submissionReq := &SubmissionUpdate{}
	err := c.Bind(submissionReq)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	submission := submissionUpdateReqToModel(submissionReq)
	ctx := c.Request().Context()
	user := getUserFromCtx(c)

	_, err = s.submissionUsecase.FindByIDAndSubmitter(ctx, submission.ID, user.ID)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	err = s.submissionUsecase.Update(c.Request().Context(), submission)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, submissionModelToRes(submission))
}
