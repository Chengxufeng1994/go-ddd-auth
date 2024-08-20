package service

import (
	"context"
	"reflect"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/entity"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/repository/facade"
)

type RoleDomainService struct {
	roleRepository facade.RoleRepository
}

func NewRoleDomainService(roleRepository facade.RoleRepository) *RoleDomainService {
	return &RoleDomainService{
		roleRepository: roleRepository,
	}
}

func (svc *RoleDomainService) Exists(ctx context.Context, roleID int) (bool, error) {
	ent, err := svc.roleRepository.GetByID(ctx, roleID)
	if err != nil {
		return false, err
	}

	if reflect.DeepEqual(ent, &entity.Role{}) {
		return false, nil
	}

	return true, nil
}
