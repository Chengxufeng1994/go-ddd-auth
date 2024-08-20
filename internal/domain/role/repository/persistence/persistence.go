package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/entity"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/repository/dao"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/repository/facade"
)

type RoleRepository struct {
	db      *gorm.DB
	roleDao dao.IRoleDao
}

var _ facade.RoleRepository = (*RoleRepository)(nil)

func NewRoleRepository(db *gorm.DB, roleDao dao.IRoleDao) *RoleRepository {
	return &RoleRepository{
		db:      db,
		roleDao: roleDao,
	}
}

func (r *RoleRepository) GetByID(ctx context.Context, roleID int) (*entity.Role, error) {
	rolePO, err := r.roleDao.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	return &entity.Role{
		ID:        rolePO.ID,
		Name:      rolePO.Name,
		CreatedAt: rolePO.CreatedAt,
		UpdatedAt: rolePO.UpdatedAt,
	}, nil
}
