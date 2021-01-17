package grader

import (
	"github.com/miun173/autograd/model"
)

// Type types of graders
type Type int

// grader types
const (
	TypeCPP = Type(0)
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
