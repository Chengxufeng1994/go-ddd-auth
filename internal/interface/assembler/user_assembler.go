package assembler

import (
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/aggregate"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/valueobject"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/interface/dto"
)

type UserAssembler struct{}

func NewUserAssembler() *UserAssembler {
	return &UserAssembler{}
}

func (u *UserAssembler) ToDTO(user *aggregate.User) *dto.UserDto {
	return &dto.UserDto{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		RoleID:    int(*user.RoleID),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (u *UserAssembler) ToDO(dto *dto.UserDto) *aggregate.User {
	roleID, _ := valueobject.NewRoleID(dto.RoleID)

	return &aggregate.User{
		ID:        dto.ID,
		Username:  dto.Username,
		Password:  dto.Password,
		RoleID:    &roleID,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}
