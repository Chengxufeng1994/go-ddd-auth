package facade

import (
	"context"
	"time"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/token/entity"
)

type TokenRepository interface {
	SaveAccessToken(ctx context.Context, authDetails *entity.AuthDetails, accessToken *entity.AuthToken, ttl time.Duration) error
	GetAccessToken(ctx context.Context, authDetails *entity.AuthDetails) (*entity.AuthToken, error)
	SaveRefreshToken(ctx context.Context, authDetails *entity.AuthDetails, refreshToken *entity.AuthToken, ttl time.Duration) error
}
