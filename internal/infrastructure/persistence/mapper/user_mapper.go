package mapper

import (
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/aggregate"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/valueobject"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/persistence/po"
)

type UserMapper struct{}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

func (m *UserMapper) ToDO(user *po.User) *aggregate.User {
	roleID, _ := valueobject.NewRoleID(int(user.RoleID))

	return &aggregate.User{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		RoleID:    &roleID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (m *UserMapper) ToPO(user *aggregate.User) *po.User {
	return &po.User{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		RoleID:    int(*user.RoleID),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
