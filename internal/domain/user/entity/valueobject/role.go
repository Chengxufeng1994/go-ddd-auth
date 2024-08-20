package valueobject

type Role struct {
	RoleID int
}

func NewRole(roleID int) *Role {
	return &Role{
		RoleID: roleID,
	}
}
