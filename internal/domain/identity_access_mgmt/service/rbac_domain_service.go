package service

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/aggregate"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/repository/facade"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/rbac"
)

type RBACDomainService struct {
	roleRepository facade.RoleRepository
	rbacAdapter    rbac.RBACAdapter
}

func NewRbacDomainService(
	roleRepository facade.RoleRepository,
	rbacAdapter rbac.RBACAdapter,
) *RBACDomainService {
	return &RBACDomainService{
		roleRepository: roleRepository,
		rbacAdapter:    rbacAdapter,
	}
}

func (svc *RBACDomainService) getRoleByID(ctx context.Context, roleID int) (*aggregate.Role, error) {
	ent, err := svc.roleRepository.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(ent, &aggregate.Role{}) {
		return nil, fmt.Errorf("role %d not found", roleID)
	}

	return ent, nil
}

func (svc *RBACDomainService) GetRoleByID(ctx context.Context, roleID int) (*aggregate.Role, error) {
	return svc.getRoleByID(ctx, roleID)
}

func (svc *RBACDomainService) ExistingRole(ctx context.Context, roleID int) (bool, error) {
	ent, err := svc.getRoleByID(ctx, roleID)
	if err != nil || ent == nil {
		return false, err
	}

	return true, nil
}

func (svc *RBACDomainService) VerifyPermission(ctx context.Context, roleID int, resource, action string) (bool, error) {
	role, err := svc.roleRepository.GetByID(ctx, roleID)
	if err != nil {
		return false, err
	}
	return svc.rbacAdapter.VerifyPermission(role.Name, resource, action)
}
