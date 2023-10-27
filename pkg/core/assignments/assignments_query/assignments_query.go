package assignments_query

import (
	"context"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/assignments"
	"github.com/fahmifan/autograd/pkg/logs"
	autogradv1 "github.com/fahmifan/autograd/pkg/pb/autograd/v1"
	"github.com/google/uuid"
)

type AssignmentsQuery struct {
	*core.Ctx
}

func (query *AssignmentsQuery) FindAssignment(
	ctx context.Context,
	req *connect.Request[autogradv1.FindByIDRequest],
) (*connect.Response[autogradv1.Assignment], error) {
	assignmentID, err := uuid.Parse(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	assignmentReader := assignments.AssignmentReader{}
	assignment, err := assignmentReader.FindByID(ctx, query.GormDB, assignmentID)
	if core.IsDBNotFoundErr(err) {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	if err != nil {
		logs.ErrCtx(ctx, err, "AssignmentsQuery: FindAssignment: FindByID")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &connect.Response[autogradv1.Assignment]{
		Msg: toAssignmentProto(assignment),
	}, nil
}

func toAssignmentProto(assignment assignments.Assignment) *autogradv1.Assignment {
	return &autogradv1.Assignment{
		Id:          assignment.ID.String(),
		Name:        assignment.Name,
		Description: assignment.Description,
		Metadata:    assignment.ProtoMeta(),
		CaseInputFile: &autogradv1.Assignment_AssignmentFile{
			Id:       assignment.CaseInputFile.ID.String(),
			Url:      assignment.CaseInputFile.URL,
			Metadata: assignment.CaseInputFile.ProtoMeta(),
		},
		CaseOutputFile: &autogradv1.Assignment_AssignmentFile{
			Id:       assignment.CaseOutputFile.ID.String(),
			Url:      assignment.CaseOutputFile.URL,
			Metadata: assignment.CaseInputFile.ProtoMeta(),
		},
	}
}
