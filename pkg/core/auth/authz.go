package auth

import (
	"github.com/samber/lo"
)

// Role ..
type Role string

// ToString ..
func (u Role) ToString() string {
	return string(u)
}

func ValidRole(role Role) bool {
	return lo.Contains(_validRoles, role)
}

// roles ..
const (
	RoleAdmin   = Role("admin")
	RoleStudent = Role("student")
)

var _validRoles = []Role{
	RoleAdmin,
	RoleStudent,
}

const _ok = true

type Permission int

const (
	CreateAssignment Permission = iota
	UpdateAssignment
	ViewAssignment
	ViewAnyAssignments
	DeleteAssignment
	GradeAssignment
	ViewAnySubmissions

	CreateSubmission
	CreateSubmissionForOther
	UpdateSubmission
	ViewSubmission
	DeleteSubmission
	DeleteSubmissionForOther

	ViewAnyUsers
	CreateAnyUser
	UpdateUser
	CreateUser

	CreateMedia

	ViewCourse
	CreateCourse
	UpdateCourse
)

var policy = map[Role]map[Permission]bool{
	RoleAdmin: {
		ViewAnyUsers:             _ok,
		CreateAssignment:         _ok,
		UpdateAssignment:         _ok,
		ViewAssignment:           _ok,
		ViewAnyAssignments:       _ok,
		DeleteAssignment:         _ok,
		GradeAssignment:          _ok,
		CreateSubmission:         _ok,
		ViewSubmission:           _ok,
		ViewAnySubmissions:       _ok,
		CreateAnyUser:            _ok,
		CreateSubmissionForOther: _ok,
		CreateMedia:              _ok,
		ViewCourse:               _ok,
		CreateCourse:             _ok,
		UpdateCourse:             _ok,
	},
	RoleStudent: {
		ViewAssignment:   _ok,
		UpdateSubmission: _ok,
		DeleteSubmission: _ok,
		UpdateUser:       _ok,
		CreateMedia:      _ok,
		CreateSubmission: _ok,
		ViewCourse:       _ok,
	},
}

// Granted check if role is granted with a permission
func (r Role) Granted(perm Permission) bool {
	role, ok := policy[r]
	if !ok {
		return false
	}

	return role[perm]
}

func (r Role) Can(perms ...Permission) bool {
	for _, perm := range perms {
		if !r.Granted(perm) {
			return false
		}
	}

	return true
}
