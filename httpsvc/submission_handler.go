package httpsvc

import (
	"net/http"
	"time"

	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/utils"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// SubmissionRequest ..
type submissionRequest struct {
	AssignmentID int64  `json:"assignmentID"`
	SubmittedBy  int64  `json:"submittedBy"`
	FileURL      string `json:"fileURL"`
}

// SubmissionResponse ..
type submissionResponse struct {
	ID           string  `json:"id"`
	AssignmentID string  `json:"assignmentID"`
	SubmittedBy  string  `json:"submittedBy"`
	FileURL      string  `json:"fileURL"`
	Grade        float64 `json:"grade"`
	Feedback     string  `json:"feedback"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

func submissionRequestToModel(r *submissionRequest) *model.Submission {
	return &model.Submission{
		AssignmentID: r.AssignmentID,
		SubmittedBy:  r.SubmittedBy,
		FileURL:      r.FileURL,
	}
}

func submissionModelToResponse(m *model.Submission) *submissionResponse {
	return &submissionResponse{
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
	submissionReq := &submissionRequest{}
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

type cursorResponse struct {
	Size      int64       `json:"size"`
	Page      int64       `json:"page"`
	Sort      string      `json:"sort"`
	TotalPage int64       `json:"totalPage"`
	TotalData int64       `json:"totalData"`
	Data      interface{} `json:"data"`
}

func newSubmissionResponses(submissions []*model.Submission) (submissionRes []*submissionResponse) {
	for _, submission := range submissions {
		submissionRes = append(submissionRes, submissionModelToResponse(submission))
	}

	return
}

func newCursorResponse(c model.Cursor, data interface{}, count int64) *cursorResponse {
	return &cursorResponse{
		Size:      c.GetSize(),
		Page:      c.GetPage(),
		Sort:      c.GetSort(),
		TotalPage: c.GetTotalPage(count),
		TotalData: count,
		Data:      data,
	}
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
