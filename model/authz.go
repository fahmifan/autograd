package model

// Role ..
type Role int

// ToString ..
func (u Role) ToString() string {
	switch u {
	case RoleAdmin:
		return "ADMIN"
	case RoleStudent:
		return "STUDENT"
	default:
		return ""
	}
}

// ParseRole ..
func ParseRole(s string) Role {
	switch s {
	case "ADMIN":
		return RoleAdmin
	case "STUDENT":
		return RoleStudent
	default:
		return RoleStudent
	}
}

// roles ..
const (
	RoleAdmin   = Role(1)
	RoleStudent = Role(2)
)

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
	UpdateSubmission
	ViewSubmission
	DeleteSubmission

	UpdateUser
	CreateUser

	CreateMedia
	ViewMedia
	ViewAnyMedia
)

var policy = map[Role]map[Permission]bool{
	RoleAdmin: {
		CreateAssignment:   _ok,
		UpdateAssignment:   _ok,
		ViewAssignment:     _ok,
		ViewAnyAssignments: _ok,
		DeleteAssignment:   _ok,
		GradeAssignment:    _ok,
		ViewSubmission:     _ok,
		ViewAnySubmissions: _ok,
		CreateUser:         _ok,
		CreateMedia:        _ok,
		ViewAnyMedia:       _ok,
	},
	RoleStudent: {
		ViewAssignment:   _ok,
		UpdateSubmission: _ok,
		DeleteSubmission: _ok,
		UpdateUser:       _ok,
		CreateMedia:      _ok,
		CreateSubmission: _ok,
		ViewMedia:        _ok,
		ViewSubmission:   _ok,
	},
}

// GrantedAny check if the Role is granted any permissions.
// If no permission given, it returns true
func (r Role) GrantedAny(perm ...Permission) bool {
	if len(perm) == 0 {
		return true
	}
	for _, p := range perm {
		if r.granted(p) {
			return true
		}
	}

	return false
}

// Granted check if role is granted with a permission
func (r Role) granted(perm Permission) bool {
	role, ok := policy[r]
	if !ok {
		return false
	}

	return role[perm]
}
