package cachelayer

import (
	"context"
	"fmt"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/aggregate"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/repository/facade"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/cache"
)

type UserCacheLayer struct {
	facade.UserRepository
	cacheStore cache.Cache
}

func NewUserCacheLayer(userRepository facade.UserRepository, cacheStore cache.Cache) *UserCacheLayer {
	return &UserCacheLayer{UserRepository: userRepository, cacheStore: cacheStore}
}

func (u *UserCacheLayer) GetUserByID(ctx context.Context, id int) (*aggregate.User, error) {
	var user *aggregate.User
	_ = u.cacheStore.Get(fmt.Sprintf("user:%d", id), &user)
	if user != nil {
		return user, nil
	}

	user, err := u.UserRepository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := u.cacheStore.SetWithDefaultExpiry(fmt.Sprintf("user:%d", id), user); err != nil {
		return nil, err
	}

	return user, nil
}
