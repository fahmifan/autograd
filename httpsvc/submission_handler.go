package httpsvc

import (
	"net/http"
	"time"

	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/utils"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func submissionRequestToModel(r *model.SubmissionRequest) *model.Submission {
	return &model.Submission{
		AssignmentID: r.AssignmentID,
		SubmittedBy:  r.SubmittedBy,
		FileURL:      r.FileURL,
	}
}

func submissionModelToResponse(m *model.Submission) *model.SubmissionResponse {
	return &model.SubmissionResponse{
		ID:           utils.Int64ToString(m.ID),
		AssignmentID: utils.Int64ToString(m.AssignmentID),
		SubmittedBy:  utils.Int64ToString(m.SubmittedBy),
		FileURL:      m.FileURL,
		Grade:        m.Grade,
		Feedback:     m.Feedback,
		CreatedAt:    m.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:    m.UpdatedAt.Format(time.RFC3339Nano),
	}
}

func (s *Server) handleCreateSubmission(c echo.Context) error {
	submissionReq := &model.SubmissionRequest{}
	err := c.Bind(submissionReq)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	submission := submissionRequestToModel(submissionReq)
	err = s.submissionUsecase.Create(c.Request().Context(), submission)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, submissionModelToResponse(submission))
}

func uploadRequestToModel(r *model.UploadRequest) *model.Upload {
	return &model.Upload{
		SourceCode: r.SourceCode,
	}
}

func uploadModelToResponse(m *model.Upload) *model.UploadResponse {
	return &model.UploadResponse{
		FileURL: m.FileURL,
	}
}

func (s *Server) handleUpload(c echo.Context) error {
	uploadReq := &model.UploadRequest{}
	err := c.Bind(uploadReq)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	upload := uploadRequestToModel(uploadReq)
	err = s.submissionUsecase.Upload(c.Request().Context(), upload)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, uploadModelToResponse(upload))
}

func cursorRequestToModel(r *model.CursorRequest) *model.Cursor {
	return &model.Cursor{
		Size: r.Size,
		Page: r.Page,
		Sort: r.Sort,
	}
}

func cursorModelToResponse(m *model.Cursor) *model.CursorResponse {
	return &model.CursorResponse{
		Size:      m.Size,
		Page:      m.Page,
		Sort:      m.Sort,
		Data:      m.Data,
		TotalPage: m.TotalPage,
		TotalData: m.TotalData,
	}
}

func (s *Server) handleGetAssignmentSubmission(c echo.Context) error {
	assignmentID := utils.StringToInt64(c.Param("assignmentID"))
	cursorRequest := generateCursorRequest(c.QueryParams())
	cursor := cursorRequestToModel(cursorRequest)
	cursor.Offset = calculateCursorOffsetValue(cursor)
	submissions, err := s.submissionUsecase.FindAllByAssignmentID(c.Request().Context(), assignmentID)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	cursor.TotalData = int64(len(submissions))
	cursor.TotalPage = calculateCursorTotalPageValue(cursor)

	submissions, err = s.submissionUsecase.FindCursorByAssignmentID(c.Request().Context(), cursor, assignmentID)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	cursor.Data = submissions

	return c.JSON(http.StatusOK, cursorModelToResponse(cursor))
}
