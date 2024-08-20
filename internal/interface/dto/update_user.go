package dto

type UpdateUserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   int    `json:"role_id"`
}
