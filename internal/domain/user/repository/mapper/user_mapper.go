package mapper

import (
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/entity"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/entity/valueobject"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/repository/po"
)

type UserMapper struct{}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

func (m *UserMapper) ToDO(user *po.User) *entity.User {
	return &entity.User{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		Role:      valueobject.NewRole(user.RoleID),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (m *UserMapper) ToPO(user *entity.User) *po.User {
	return &po.User{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		RoleID:    user.Role.RoleID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
