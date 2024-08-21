package token

import (
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/token/entity"
)

type TokenEnhancer interface {
	SignToken(authToken *entity.AuthToken, authDetails *entity.AuthDetails) (*entity.AuthToken, error)
	VerifyToken(tokenValue string) (*entity.AuthToken, *entity.AuthDetails, error)
}
