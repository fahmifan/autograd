package httpsvc

import (
	"net/http"
	"time"

	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/utils"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type submissionRequest struct {
	AssignmentID int64
	SubmittedBy  int64
}

func (s *submissionRequest) toModel() *model.Submission {
	return &model.Submission{
		AssignmentID: s.AssignmentID,
		SubmittedBy:  s.SubmittedBy,
	}
}

type submissionRes struct {
	ID           string  `json:"id"`
	AssignmentID string  `json:"assignment_id"`
	SubmittedBy  string  `json:"submitted_by"`
	FileURL      string  `json:"file_url"`
	Grade        float64 `json:"grade"`
	Feedback     string  `json:"feedback"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

func submissionResFromModel(m *model.Submission, fileURL string) *submissionRes {
	return &submissionRes{
		ID:           utils.Int64ToString(m.ID),
		AssignmentID: utils.Int64ToString(m.AssignmentID),
		SubmittedBy:  utils.Int64ToString(m.SubmittedBy),
		FileURL:      fileURL,
		Grade:        m.Grade,
		Feedback:     m.Feedback,
		CreatedAt:    m.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:    m.UpdatedAt.Format(time.RFC3339Nano),
	}
}

func (s *Server) handleCreateSubmission(c echo.Context) error {

	form, err := c.MultipartForm()

	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	files := form.File["files"]
	fileURLs := []string{}

	for _, file := range files {

		fileURL, err := s.submissionUsecase.Upload(file)

		if err != nil {
			logrus.Error(err)
			return responseError(c, err)
		}

		fileURLs = append(fileURLs, fileURL)
	}

	submissionReq := &submissionRequest{}
	err = c.Bind(submissionReq)

	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	submission := submissionReq.toModel()
	err = s.submissionUsecase.Create(c.Request().Context(), submission, fileURLs)

	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	submissionResArr := []*submissionRes{}

	for _, fileURL := range fileURLs {
		submissionResArr = append(submissionResArr, submissionResFromModel(submission, fileURL))
	}

	return c.JSON(http.StatusOK, map[string][]*submissionRes{"submissions": submissionResArr})
}
