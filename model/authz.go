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
	ActionView   Action = "get"
	ActionEdit   Action = "edit"
	ActionCreate Action = "create"
	ActionDelete Action = "delete"
	ActionGrade  Action = "grade"
)

const _ok = true

var roleResourcesActions = map[Role]map[Resource]map[Action]bool{
	RoleAdmin: {
		ResourceUser: {
			ActionView: _ok,
		},
		ResourceAssignment: {
			ActionCreate: _ok,
			ActionGrade:  _ok,
		},
		ResourceSubmission: {
			ActionView: _ok,
		},
	},
	RoleStudent: {
		ResourceUser: {
			ActionView: _ok,
			ActionEdit: _ok,
		},
		ResourceSubmission: {
			ActionCreate: _ok,
			ActionEdit:   _ok,
			ActionView:   _ok,
			ActionDelete: _ok,
		},
	},
}

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

	CreateMedia
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
	},
	RoleStudent: {
		ViewAssignment:   _ok,
		UpdateSubmission: _ok,
		DeleteSubmission: _ok,
		UpdateUser:       _ok,
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

// Authorized check if given Role has access to Resouce rsc with Action act
// if an action is not explecitly stated then it was consider as unauthorized
func (r Role) Authorized(rsc Resource, acts ...Action) bool {
	roleRscActions, ok := roleResourcesActions[r]
	if !ok {
		return false
	}

	rscActions, ok := roleRscActions[rsc]
	if !ok {
		return false
	}

	for _, act := range acts {
		isGranted, ok := rscActions[act]
		if ok {
			return isGranted
		}
	}

	return false
}
