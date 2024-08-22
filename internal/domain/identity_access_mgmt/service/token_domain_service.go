package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/entity"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/repository/facade"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/token"

	"golang.org/x/net/context"
)

var (
	accessTokenDuration  time.Duration = 12 * time.Hour
	refreshTokenDuration time.Duration = 24 * time.Hour
)

type TokenDomainService struct {
	tokenEnhancer   token.TokenEnhancer
	tokenRepository facade.TokenRepository
}

func NewTokenDomainService(
	tokenEnhancer token.TokenEnhancer,
	tokenRepository facade.TokenRepository,
) *TokenDomainService {
	return &TokenDomainService{
		tokenEnhancer:   tokenEnhancer,
		tokenRepository: tokenRepository,
	}
}

func (svc *TokenDomainService) CreateAccessToken(ctx context.Context, userDetails *entity.UserDetails) (*entity.AuthToken, *entity.AuthToken, error) {
	authDetails := &entity.AuthDetails{
		User: userDetails,
	}

	accessToken, err := svc.createAccessToken(ctx, authDetails)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := svc.createRefreshToken(ctx, authDetails)
	if err != nil {
		return nil, nil, err
	}

	err = svc.tokenRepository.SaveAccessToken(ctx, authDetails, accessToken, accessTokenDuration)
	if err != nil {
		return nil, nil, err
	}
	err = svc.tokenRepository.SaveRefreshToken(ctx, authDetails, refreshToken, accessTokenDuration)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (svc *TokenDomainService) createAccessToken(_ context.Context, authDetails *entity.AuthDetails) (*entity.AuthToken, error) {
	expiresTime := time.Now().Add(accessTokenDuration)
	accessToken := &entity.AuthToken{
		TokenValue:  "",
		ExpiresTime: &expiresTime,
	}

	return svc.tokenEnhancer.SignToken(accessToken, authDetails)
}

func (svc *TokenDomainService) createRefreshToken(_ context.Context, authDetails *entity.AuthDetails) (*entity.AuthToken, error) {
	expiresTime := time.Now().Add(refreshTokenDuration)
	refreshToken := &entity.AuthToken{
		TokenValue:  "",
		ExpiresTime: &expiresTime,
	}

	return svc.tokenEnhancer.SignToken(refreshToken, authDetails)
}

func (svc *TokenDomainService) CreateRefreshToken(ctx context.Context, userDetails *entity.UserDetails) (*entity.AuthToken, error) {
	authDetails := &entity.AuthDetails{
		User: userDetails,
	}
	expiresTime := time.Now().Add(refreshTokenDuration)
	refreshToken := &entity.AuthToken{
		TokenValue:  "",
		ExpiresTime: &expiresTime,
	}

	var err error
	refreshToken, err = svc.tokenEnhancer.SignToken(refreshToken, authDetails)
	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (svc *TokenDomainService) RemoveAccessToken(ctx context.Context) {}

func (svc *TokenDomainService) RemoveRefreshToken(ctx context.Context) {}

func (svc *TokenDomainService) VerifyToken(ctx context.Context, tokenValue string) (*entity.AuthToken, *entity.AuthDetails, error) {
	authToken, authDetails, err := svc.tokenEnhancer.VerifyToken(tokenValue)
	if err != nil {
		return nil, nil, err
	}

	existingAuthToken, err := svc.tokenRepository.GetAccessToken(ctx, authDetails)
	if existingAuthToken == nil {
		return nil, nil, fmt.Errorf("invalid token")
	}
	if !strings.EqualFold(existingAuthToken.TokenValue, authToken.TokenValue) {
		return nil, nil, fmt.Errorf("invalid token")
	}
	return authToken, authDetails, err
}
