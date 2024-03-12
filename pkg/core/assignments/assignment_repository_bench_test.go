package assignments_test

import (
	"testing"
	"time"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/assignments"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

var asg = assignments.Assignment{
	ID:          uuid.New(),
	Name:        "Assignment 1",
	Description: "Description 1",
	DeadlineAt:  time.Time{},
	Assigner: assignments.Assigner{
		ID:     uuid.New(),
		Name:   "Assigner 1",
		Active: true,
	},
	CaseInputFile: assignments.CaseFile{
		ID:   uuid.New(),
		URL:  "http://example.com",
		Type: "input",
		TimestampMetadata: core.TimestampMetadata{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: null.NewTime(time.Now(), false),
		},
	},
	CaseOutputFile: assignments.CaseFile{
		ID:   uuid.New(),
		URL:  "http://example.com",
		Type: "input",
		TimestampMetadata: core.TimestampMetadata{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: null.NewTime(time.Now(), false),
		},
	},
	TimestampMetadata: core.NewTimestampMeta(time.Time{}),
}

func BenchmarkAssignmentNoPrealloc_10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var res []assignments.Assignment
		for range 10 {
			res = append(res, asg)
		}
	}
}

func BenchmarkAssignmentPrealloc_10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res := make([]assignments.Assignment, 10)
		for i := range 10 {
			res[i] = asg
		}
	}
}

func BenchmarkAssignmentNoPrealloc_20(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var res []assignments.Assignment
		for range 20 {
			res = append(res, asg)
		}
	}
}

func BenchmarkAssignmentPrealloc_20(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res := make([]assignments.Assignment, 20)
		for i := range 20 {
			res[i] = asg
		}
	}
}

func BenchmarkAssignmentNoPrealloc_50(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var res []assignments.Assignment
		for range 50 {
			res = append(res, asg)
		}
	}
}

func BenchmarkAssignmentPrealloc_50(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res := make([]assignments.Assignment, 50)
		for i := range 50 {
			res[i] = asg
		}
	}
}
