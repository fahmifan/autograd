package httpsvc

import (
	"time"

	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/utils"
)

type assignmentUpdateRequest struct {
	ID                string `json:"id"`
	AssignedBy        string `json:"assignedBy"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	CaseInputFileURL  string `json:"caseInputFileURL"`
	CaseOutputFileURL string `json:"caseOutputFileURL"`
}

func assigmentUpdateReqToModel(r *assignmentUpdateRequest) *model.Assignment {
	return &model.Assignment{
		ID:                utils.StringToInt64(r.ID),
		AssignedBy:        utils.StringToInt64(r.AssignedBy),
		Name:              r.Name,
		Description:       r.Description,
		CaseInputFileURL:  r.CaseInputFileURL,
		CaseOutputFileURL: r.CaseOutputFileURL,
	}
}

type assignmentRequest struct {
	AssignedBy        string `json:"assignedBy"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	CaseInputFileURL  string `json:"caseInputFileURL"`
	CaseOutputFileURL string `json:"caseOutputFileURL"`
}

func assignmentRequestToModel(r *assignmentRequest) *model.Assignment {
	return &model.Assignment{
		AssignedBy:        utils.StringToInt64(r.AssignedBy),
		Name:              r.Name,
		Description:       r.Description,
		CaseInputFileURL:  r.CaseInputFileURL,
		CaseOutputFileURL: r.CaseOutputFileURL,
	}
}

type assignmentResponse struct {
	ID                string `json:"id"`
	AssignedBy        string `json:"assignedBy"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	CaseInputFileURL  string `json:"caseInputFileURL"`
	CaseOutputFileURL string `json:"caseOutputFileURL"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

func assignmentModelToResponse(m *model.Assignment) *assignmentResponse {
	return &assignmentResponse{
		ID:                utils.Int64ToString(m.ID),
		AssignedBy:        utils.Int64ToString(m.AssignedBy),
		Name:              m.Name,
		Description:       m.Description,
		CaseInputFileURL:  m.CaseInputFileURL,
		CaseOutputFileURL: m.CaseOutputFileURL,
		CreatedAt:         m.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:         m.UpdatedAt.Format(time.RFC3339Nano),
	}
}

type cursorResponse struct {
	Size      string      `json:"size"`
	Page      string      `json:"page"`
	Sort      string      `json:"sort"`
	TotalPage string      `json:"totalPage"`
	TotalData string      `json:"totalData"`
	Data      interface{} `json:"data"`
}

func newCursorResponse(c model.Cursor, data interface{}, count int64) *cursorResponse {
	return &cursorResponse{
		Size:      utils.Int64ToString(c.GetSize()),
		Page:      utils.Int64ToString(c.GetPage()),
		Sort:      c.GetSort(),
		TotalPage: utils.Int64ToString(c.GetTotalPage(count)),
		TotalData: utils.Int64ToString(count),
		Data:      data,
	}
}

type uploadRequest struct {
	SourceCode string `json:"sourceCode"`
}

type uploadResponse struct {
	FileURL string `json:"fileURL"`
}

type submissionRequest struct {
	AssignmentID int64  `json:"assignmentID"`
	SubmittedBy  int64  `json:"submittedBy"`
	FileURL      string `json:"fileURL"`
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
