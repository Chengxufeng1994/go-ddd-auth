package valueobject

import "fmt"

type RoleID int

func NewRoleID(roleID int) (RoleID, error) {
	if roleID <= 0 {
		return 0, fmt.Errorf("invalid role id: %d", roleID)
	}
	return RoleID(roleID), nil
}
