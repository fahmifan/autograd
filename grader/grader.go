package grader

import (
	"github.com/fahmifan/autograd/model"
)

// Type types of graders
type Type string

// grader types
const (
	TypeCPP = Type("cpp")
)

// New new grader factory
func New(t Type) model.GraderEngine {
	switch t {
	case TypeCPP:
		return &CPPGrader{}
	default:
		return nil
	}
}
