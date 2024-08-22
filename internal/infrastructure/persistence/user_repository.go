package persistence

import (
	"context"
	"fmt"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/aggregate"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/repository/facade"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/persistence/dao"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/persistence/mapper"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/persistence/po"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/transaction"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/util"
)

type UserRepository struct {
	factory    transaction.DBFactory
	UserDao    *dao.UserDao
	UserMapper *mapper.UserMapper
}

var _ facade.UserRepository = (*UserRepository)(nil)

func NewUserRepository(factory transaction.DBFactory, userDao *dao.UserDao) *UserRepository {
	return &UserRepository{
		factory:    factory,
		UserDao:    userDao,
		UserMapper: mapper.NewUserMapper(),
	}
}

func (r *UserRepository) Save(ctx context.Context, aggregate *aggregate.User) error {
	db := r.factory.GetDB(ctx)
	return db.Create(r.UserMapper.ToPO(aggregate)).Error
}

func (r *UserRepository) Update(ctx context.Context, aggregate *aggregate.User) error {
	db := r.factory.GetDB(ctx)
	diff := aggregate.DetectChanges()
	if diff != nil {
		if diff.UserChanged {
			fmt.Println("123")
			err := db.Debug().Model(&po.User{}).
				Where("id = ?", aggregate.ID).
				Updates(map[string]interface{}{
					"Username":  aggregate.Username,
					"password":  aggregate.Password,
					"role_id":   aggregate.RoleID,
					"UpdatedAt": aggregate.UpdatedAt,
				}).Error
			if err != nil {
				return err
			}
		}
	}

	aggregate.Attach()
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID int) (*aggregate.User, error) {
	tx := r.factory.GetDB(ctx)
	userPo, _ := r.UserDao.GetByID(ctx, tx, userID)
	if userPo == nil {
		return nil, fmt.Errorf("User %d not found", userID)
	}
	userDO := r.UserMapper.ToDO(userPo)
	userDO.Attach()
	return userDO, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*aggregate.User, error) {
	db := r.factory.GetDB(ctx)
	userPO, err := r.UserDao.GetByUsername(ctx, db, username)
	if err != nil {
		return nil, err
	}
	if userPO == nil {
		return nil, nil
	}

	return r.UserMapper.ToDO(userPO), nil
}

func (r *UserRepository) SearchUsers(ctx context.Context, opts *aggregate.SearchUserOpts) ([]*aggregate.User, *util.Pagination, error) {
	db := r.factory.GetDB(ctx)
	lis, pagination, err := r.UserDao.Search(ctx, db, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("search users error: %w", err)
	}
	users := make([]*aggregate.User, 0, len(lis))
	for _, userPO := range lis {
		users = append(users, r.UserMapper.ToDO(userPO))
	}

	return users, pagination, nil
}

func (r *UserRepository) Remove(ctx context.Context, userID int) error {
	db := r.factory.GetDB(ctx)
	return r.UserDao.DeleteByID(ctx, db, userID)
}
