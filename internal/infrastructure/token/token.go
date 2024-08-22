package token

import "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/entity"

type TokenEnhancer interface {
	SignToken(authToken *entity.AuthToken, authDetails *entity.AuthDetails) (*entity.AuthToken, error)
	VerifyToken(tokenValue string) (*entity.AuthToken, *entity.AuthDetails, error)
}
