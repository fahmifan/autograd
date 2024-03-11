package student_assignment_cmd

import (
	"context"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/grading/grading_cmd"
	"github.com/fahmifan/autograd/pkg/jobqueue"
	"github.com/fahmifan/autograd/pkg/logs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const JobGradeSubmission = jobqueue.JobType("grade_submission")

type GradeStudentSubmissionHandler struct {
	*core.Ctx
}

type GradeStudentSubmissionPayload struct {
	SubmissionID uuid.UUID
}

func (handler *GradeStudentSubmissionHandler) JobType() jobqueue.JobType {
	return JobGradeSubmission
}

func (handler *GradeStudentSubmissionHandler) Handle(ctx context.Context, tx *gorm.DB, payload jobqueue.Payload) error {
	req := GradeStudentSubmissionPayload{}
	err := jobqueue.UnmarshalPayload(payload, &req)
	if err != nil {
		return logs.ErrWrapCtx(ctx, err, "GradeStudentSubmissionHandler: Handle: json.Unmarshal")
	}

	gradingCmd := &grading_cmd.GradingCmd{Ctx: handler.Ctx}

	_, err = gradingCmd.InternalGradeSubmissionTx(ctx, tx, grading_cmd.InternalGradeSubmissionRequest{
		SubmissionID: req.SubmissionID,
	})
	if err != nil {
		return logs.ErrWrapCtx(ctx, err, "GradeStudentSubmissionHandler: Handle: InternalGradeSubmissionTx")
	}

	return nil
}
