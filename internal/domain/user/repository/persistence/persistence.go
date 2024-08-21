package persistence

import (
	"context"
	"fmt"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/entity"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/repository/dao"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/repository/facade"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/repository/mapper"
	"gorm.io/gorm"
)

type UserRepository struct {
	db         *gorm.DB
	UserDao    *dao.UserDao
	UserMapper *mapper.UserMapper
}

var _ facade.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *gorm.DB, userDao *dao.UserDao) *UserRepository {
	return &UserRepository{
		db:         db,
		UserDao:    userDao,
		UserMapper: mapper.NewUserMapper(),
	}
}

func (r *UserRepository) SaveUser(ctx context.Context, aggregate *entity.User) error {
	po := r.UserMapper.ToPO(aggregate)
	return r.db.Transaction(func(tx *gorm.DB) error {
		return r.UserDao.Save(ctx, tx, po)
	})
}

func (r *UserRepository) UpdateUser(ctx context.Context, id int, fn facade.UpdateFn) error {
	var err error
	tx := r.db.Debug().Begin()

	defer func() {
		err = r.finishTransaction(ctx, err, tx)
	}()

	existingUser, err := r.getUserByID(ctx, tx, id)
	if err != nil {
		return err
	}

	updatedUser, err := fn(existingUser)

	err = r.SaveUser(ctx, updatedUser)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) finishTransaction(_ context.Context, err error, tx *gorm.DB) error {
	if err != nil {
		return tx.Rollback().Error
	}

	return tx.Commit().Error
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	return r.getUserByID(ctx, r.db, userID)
}

func (r *UserRepository) getUserByID(ctx context.Context, tx *gorm.DB, userID int) (*entity.User, error) {
	userPo, _ := r.UserDao.GetByID(ctx, tx, userID)
	if userPo == nil {
		return nil, fmt.Errorf("User %d not found", userID)
	}
	return r.UserMapper.ToDO(userPo), nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	userPO, _ := r.UserDao.GetByUsername(ctx, r.db, username)
	if userPO == nil {
		return nil, fmt.Errorf("User %s not found", username)
	}

	return r.UserMapper.ToDO(userPO), nil
}

func (r *UserRepository) DeleteUserByID(ctx context.Context, userID int) error {
	return r.UserDao.DeleteByID(ctx, r.db, userID)
}
