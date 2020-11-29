package httpsvc

import (
	"net/http"
	"time"

	"github.com/miun173/autograd/dto"
	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/utils"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type submissionRequest struct {
	AssignmentID int64  `json:"assignmentID"`
	SubmittedBy  int64  `json:"submittedBy"`
	FileURL      string `json:"fileURL"`
}

func (s *submissionRequest) toModel() *model.Submission {
	return &model.Submission{
		AssignmentID: s.AssignmentID,
		SubmittedBy:  s.SubmittedBy,
		FileURL:      s.FileURL,
	}
}

type submissionRes struct {
	ID           string  `json:"id"`
	AssignmentID string  `json:"assignmentID"`
	SubmittedBy  string  `json:"submittedBy"`
	FileURL      string  `json:"fileURL"`
	Grade        float64 `json:"grade"`
	Feedback     string  `json:"feedback"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

func submissionResFromModel(m *model.Submission) *submissionRes {
	return &submissionRes{
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

	submission := submissionReq.toModel()
	err = s.submissionUsecase.Create(c.Request().Context(), submission)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, submissionResFromModel(submission))
}

type uploadRequest struct {
	SourceCode string `json:"sourceCode"`
}

func (u *uploadRequest) toDTO() *dto.Upload {
	return &dto.Upload{
		SourceCode: u.SourceCode,
	}
}

type uploadRes struct {
	FileURL string `json:"fileURL"`
}

func uploadResFromDTO(u *dto.Upload) *uploadRes {
	return &uploadRes{
		FileURL: u.FileURL,
	}
}

func (s *Server) handleUpload(c echo.Context) error {
	uploadReq := &uploadRequest{}
	err := c.Bind(uploadReq)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	upload := uploadReq.toDTO()
	err = s.submissionUsecase.Upload(c.Request().Context(), upload)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, uploadResFromDTO(upload))
}

func (s *Server) handleGetAssignmentSubmission(c echo.Context) error {
	assignmentID := utils.StringToInt64(c.Param("assignmentID"))
	submissions, err := s.submissionUsecase.FindByAssignmentID(c.Request().Context(), assignmentID)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, submissions)
}
