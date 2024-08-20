package dao

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/repository/po"
	"gorm.io/gorm"
)

type IRoleDao interface {
	GetByID(ctx context.Context, roleID int) (*po.Role, error)
}

type GormRoleDao struct {
	db *gorm.DB
}

var _ IRoleDao = (*GormRoleDao)(nil)

func NewRoleDao(db *gorm.DB) *GormRoleDao {
	return &GormRoleDao{
		db: db,
	}
}

func (dao *GormRoleDao) GetByID(ctx context.Context, roleID int) (*po.Role, error) {
	var role po.Role
	if err := dao.db.WithContext(ctx).Model(&po.Role{}).Where("id = ?", roleID).First(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}
