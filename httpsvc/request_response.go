package httpsvc

import (
	"time"

	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/utils"
)

type submissionRequest struct {
	AssignmentID int64  `json:"assignmentID"`
	SubmittedBy  int64  `json:"submittedBy"`
	FileURL      string `json:"fileURL"`
}

func submissionRequestToModel(r *submissionRequest) *model.Submission {
	return &model.Submission{
		AssignmentID: r.AssignmentID,
		SubmittedBy:  r.SubmittedBy,
		FileURL:      r.FileURL,
	}
}

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

func newSubmissionResponses(submissions []*model.Submission) (submissionRes []*submissionResponse) {
	for _, submission := range submissions {
		submissionRes = append(submissionRes, submissionModelToResponse(submission))
	}

	return
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

type cursorResponse struct {
	Size      int64       `json:"size"`
	Page      int64       `json:"page"`
	Sort      string      `json:"sort"`
	TotalPage int64       `json:"totalPage"`
	TotalData int64       `json:"totalData"`
	Data      interface{} `json:"data"`
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
