package jwt

import (
	"time"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/token/entity"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructue/token"

	"github.com/golang-jwt/jwt/v5"
)

type AuthTokenCustomClaims struct {
	UserDetails entity.UserDetails
	jwt.RegisteredClaims
}

type JwtTokenEnhancer struct {
	secretKey []byte
}

var _ token.TokenEnhancer = (*JwtTokenEnhancer)(nil)

func NewJwtTokenEnhancer(secretKey []byte) *JwtTokenEnhancer {
	return &JwtTokenEnhancer{
		secretKey: secretKey,
	}
}

func (enhancer *JwtTokenEnhancer) SignToken(authToken *entity.AuthToken, authDetails *entity.AuthDetails) (*entity.AuthToken, error) {
	expireTime := authToken.ExpiresTime
	userDetails := *authDetails.User

	claims := AuthTokenCustomClaims{
		UserDetails: userDetails,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userDetails.Username,
			Issuer:    "license-server",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(*expireTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenValue, err := token.SignedString(enhancer.secretKey)
	if err != nil {
		return nil, err
	}

	authToken.TokenValue = tokenValue
	authToken.TokenType = "jwt"

	return authToken, nil
}

func (enhancer *JwtTokenEnhancer) VerifyToken(tokenValue string) (*entity.AuthToken, *entity.AuthDetails, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return enhancer.secretKey, nil
	}
	jwtToken, err := jwt.ParseWithClaims(tokenValue, &AuthTokenCustomClaims{}, keyFunc)
	if err != nil {
		return nil, nil, err
	}

	claims := jwtToken.Claims.(*AuthTokenCustomClaims)

	expiresTime := time.Unix(claims.ExpiresAt.Unix(), 0)

	return &entity.AuthToken{
			TokenValue:  tokenValue,
			ExpiresTime: &expiresTime,
		}, &entity.AuthDetails{
			User: &claims.UserDetails,
		}, nil
}
