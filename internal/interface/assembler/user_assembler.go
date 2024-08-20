package assembler

import (
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/entity"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/entity/valueobject"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/interface/dto"
)

type UserAssembler struct{}

func NewUserAssembler() *UserAssembler {
	return &UserAssembler{}
}

func (u *UserAssembler) ToDTO(user *entity.User) *dto.UserDto {
	return &dto.UserDto{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		RoleID:    user.Role.RoleID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (u *UserAssembler) ToDO(dto *dto.UserDto) *entity.User {
	return &entity.User{
		ID:        dto.ID,
		Username:  dto.Username,
		Password:  dto.Password,
		Role:      valueobject.NewRole(dto.RoleID),
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}
