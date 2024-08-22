package persistence

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/aggregate"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/repository/facade"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/persistence/dao"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db      *gorm.DB
	roleDao dao.IRoleDao
}

var _ facade.RoleRepository = (*RoleRepository)(nil)

func NewRbacRepository(db *gorm.DB, roleDao dao.IRoleDao) *RoleRepository {
	return &RoleRepository{
		db:      db,
		roleDao: roleDao,
	}
}

func (r *RoleRepository) GetByID(ctx context.Context, roleID int) (*aggregate.Role, error) {
	rolePO, err := r.roleDao.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	return &aggregate.Role{
		ID:        aggregate.RoleID(rolePO.ID),
		Name:      rolePO.Name,
		CreatedAt: rolePO.CreatedAt,
		UpdatedAt: rolePO.UpdatedAt,
	}, nil
}
