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

type Resource string
type Action string

// roles ..
const (
	RoleAdmin   = Role(1)
	RoleStudent = Role(2)
)

// resources ..
const (
	ResourceUser       Resource = "users"
	ResourceAssignment Resource = "assignments"
	ResourceSubmission Resource = "submissions"
)

// actions ..
const (
	ActionGetAny  Action = "get_any"
	ActionGetSelf Action = "get_self"

	ActionEditAny  Action = "edit_any"
	ActionEditSelf Action = "edit_self"

	ActionCreateAny  Action = "create_any"
	ActionCreateSelf Action = "create_self"

	ActionGradeAny Action = "grade_any"
)

const _ok = true

var roleResourcesActions = map[Role]map[Resource]map[Action]bool{
	RoleAdmin: {
		ResourceUser: {
			ActionEditAny: _ok,
			ActionGetAny:  _ok,
			ActionGetSelf: _ok,
		},
		ResourceAssignment: {
			ActionCreateAny: _ok,
			ActionGetAny:    _ok,
			ActionGradeAny:  _ok,
		},
	},
	RoleStudent: {
		ResourceUser: {
			ActionGetSelf:  _ok,
			ActionEditSelf: _ok,
		},
		ResourceSubmission: {
			ActionCreateSelf: _ok,
			ActionEditSelf:   _ok,
			ActionGetSelf:    _ok,
		},
	},
}

// HasAccess check if given Role has access to Resouce rsc with Action act
func (r Role) HasAccess(rsc Resource, act Action) bool {
	roleRscActions, ok := roleResourcesActions[r]
	if !ok {
		return false
	}

	rscActions, ok := roleRscActions[rsc]
	if !ok {
		return false
	}

	isGranted, ok := rscActions[act]
	if !ok {
		return false
	}

	return isGranted
}
