package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/entity"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/repository/facade"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/cache"
)

type BadgerTokenRepository struct {
	db cache.Cache
}

var _ facade.TokenRepository = (*BadgerTokenRepository)(nil)

func NewTokenRepository(db cache.Cache) *BadgerTokenRepository {
	return &BadgerTokenRepository{
		db: db,
	}
}

func buildAccessTokenKey(userID int) string {
	key := fmt.Sprintf("access_token:user_id:%d", userID)
	return key
}

func buildRefreshTokenKey(userID int) string {
	key := fmt.Sprintf("refresh_token:user_id:%d", userID)
	return key
}

func (r *BadgerTokenRepository) SaveAccessToken(ctx context.Context, authDetails *entity.AuthDetails, accessToken *entity.AuthToken, ttl time.Duration) error {
	return r.db.SetWithExpiry(buildAccessTokenKey(authDetails.User.UserID), accessToken, ttl)
}

func (r *BadgerTokenRepository) GetAccessToken(ctx context.Context, authDetails *entity.AuthDetails) (*entity.AuthToken, error) {

	var authToken entity.AuthToken
	err := r.db.Get(buildAccessTokenKey(authDetails.User.UserID), &authToken)
	if err != nil {
		return nil, err
	}

	return &authToken, nil
}

func (r *BadgerTokenRepository) SaveRefreshToken(ctx context.Context, authDetails *entity.AuthDetails, refreshToken *entity.AuthToken, ttl time.Duration) error {
	return r.db.SetWithExpiry(buildRefreshTokenKey(authDetails.User.UserID), refreshToken, ttl)
}
