package httpsvc

import (
	"time"

	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/utils"
)

type assignmentReq struct {
	ID                string `json:"id,omitempty"`
	AssignedBy        string `json:"assignedBy"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	CaseInputFileURL  string `json:"caseInputFileURL"`
	CaseOutputFileURL string `json:"caseOutputFileURL"`
}

func assignmentCreateReqToModel(r *assignmentReq) *model.Assignment {
	return &model.Assignment{
		AssignedBy:        utils.StringToInt64(r.AssignedBy),
		Name:              r.Name,
		Description:       r.Description,
		CaseInputFileURL:  r.CaseInputFileURL,
		CaseOutputFileURL: r.CaseOutputFileURL,
	}
}

func assigmentUpdateReqToModel(r *assignmentReq) *model.Assignment {
	return &model.Assignment{
		ID:                utils.StringToInt64(r.ID),
		AssignedBy:        utils.StringToInt64(r.AssignedBy),
		Name:              r.Name,
		Description:       r.Description,
		CaseInputFileURL:  r.CaseInputFileURL,
		CaseOutputFileURL: r.CaseOutputFileURL,
	}
}

type assignmentRes struct {
	ID                string `json:"id"`
	AssignedBy        string `json:"assignedBy"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	CaseInputFileURL  string `json:"caseInputFileURL"`
	CaseOutputFileURL string `json:"caseOutputFileURL"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
	DeletedAt         string `json:"deletedAt,omitempty"`
}

func assignmentModelToRes(m *model.Assignment) *assignmentRes {
	return &assignmentRes{
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

func newAssignmentResponses(assignments []*model.Assignment) (assignmentRes []*assignmentRes) {
	for _, assignment := range assignments {
		assignmentRes = append(assignmentRes, assignmentModelToRes(assignment))
	}

	return
}

func assignmentModelToDeleteRes(m *model.Assignment) *assignmentRes {
	return &assignmentRes{
		ID:                utils.Int64ToString(m.ID),
		AssignedBy:        utils.Int64ToString(m.AssignedBy),
		Name:              m.Name,
		Description:       m.Description,
		CaseInputFileURL:  m.CaseInputFileURL,
		CaseOutputFileURL: m.CaseOutputFileURL,
		CreatedAt:         m.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:         m.UpdatedAt.Format(time.RFC3339Nano),
		DeletedAt:         m.DeletedAt.Time.Format(time.RFC3339Nano),
	}
}

type cursorRes struct {
	Size      string      `json:"size"`
	Page      string      `json:"page"`
	Sort      string      `json:"sort"`
	TotalPage string      `json:"totalPage"`
	TotalData string      `json:"totalData"`
	Data      interface{} `json:"data"`
}

func newCursorRes(c model.Cursor, data interface{}, count int64) *cursorRes {
	return &cursorRes{
		Size:      utils.Int64ToString(c.GetSize()),
		Page:      utils.Int64ToString(c.GetPage()),
		Sort:      c.GetSort(),
		TotalPage: utils.Int64ToString(c.GetTotalPage(count)),
		TotalData: utils.Int64ToString(count),
		Data:      data,
	}
}

type uploadReq struct {
	SourceCode string `json:"sourceCode"`
}

type uploadRes struct {
	FileURL string `json:"fileURL"`
}

type submissionReq struct {
	ID           string `json:"id,omitempty"`
	AssignmentID int64  `json:"assignmentID"`
	SubmittedBy  int64  `json:"submittedBy"`
	FileURL      string `json:"fileURL"`
}

func submissionCreateReqToModel(s *submissionReq) *model.Submission {
	return &model.Submission{
		AssignmentID: s.AssignmentID,
		SubmittedBy:  s.SubmittedBy,
		FileURL:      s.FileURL,
	}
}

func submissionUpdateReqToModel(s *submissionReq) *model.Submission {
	return &model.Submission{
		ID:           utils.StringToInt64(s.ID),
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

func submissionModelToRes(m *model.Submission) *submissionRes {
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

func newSubmissionResponses(submissions []*model.Submission) (submissionRes []*submissionRes) {
	for _, submission := range submissions {
		submissionRes = append(submissionRes, submissionModelToRes(submission))
	}

	return
}
