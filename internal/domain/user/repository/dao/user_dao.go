package dao

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/repository/po"
	"gorm.io/gorm"
)

type IUserDao interface {
	Save(ctx context.Context, tx *gorm.DB, po *po.User) error
	GetByID(ctx context.Context, tx *gorm.DB, id int) (*po.User, error)
	DeleteByID(ctx context.Context, tx *gorm.DB, id int) error
	GetByUsername(ctx context.Context, tx *gorm.DB, username string) (*po.User, error)
}

type UserDao struct {
	db *gorm.DB
}

var _ IUserDao = (*UserDao)(nil)

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) Save(ctx context.Context, tx *gorm.DB, userPO *po.User) error {
	return tx.WithContext(ctx).
		Model(&po.User{}).
		Where("id = ?", userPO.ID).
		Save(&userPO).Error
}

func (dao *UserDao) GetByID(ctx context.Context, tx *gorm.DB, id int) (*po.User, error) {
	var user po.User
	err := tx.WithContext(ctx).
		Model(&po.User{}).
		Preload("Role").
		Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (dao *UserDao) DeleteByID(ctx context.Context, tx *gorm.DB, id int) error {
	return tx.WithContext(ctx).
		Model(&po.User{}).
		Where("id = ?", id).
		Delete(&po.User{}).Error
}

func (dao *UserDao) GetByUsername(ctx context.Context, tx *gorm.DB, username string) (*po.User, error) {
	var user po.User
	err := tx.Debug().WithContext(ctx).
		Model(&po.User{}).
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
