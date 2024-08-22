package service

import (
	"fmt"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/assert"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/util"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 64
)

type PasswordDomainService struct {
}

func NewPasswordDomainService() *PasswordDomainService {
	return &PasswordDomainService{}
}

func (svc *PasswordDomainService) EncryptPassword(rawPassword string) (string, error) {
	err := svc.ValidatePassword(rawPassword)
	if err != nil {
		return "", err
	}
	return util.HashPassword(rawPassword)
}

func (svc *PasswordDomainService) VerifyPassword(rawPassword string, hashedPassword string) error {
	return util.ComparePassword(hashedPassword, rawPassword)
}

// validatePassword check password is not empty and length between 8 and 64
func (svc *PasswordDomainService) ValidatePassword(rawPassword string) error {
	if err := assert.AssertArgumentNotEmpty(rawPassword); err != nil {
		return err
	}
	if err := assert.AssertArgumentLength(rawPassword, MinPasswordLength, MaxPasswordLength); err != nil {
		return err
	}
	return nil
}

func (svc *PasswordDomainService) ChangedPassword(currentPassword string, changedPassword string) (string, error) {
	err := svc.VerifyPassword(changedPassword, currentPassword)
	if err == nil {
		return "", fmt.Errorf("password cannot be the same as current password")
	}

	err = svc.ValidatePassword(changedPassword)
	if err != nil {
		return "", fmt.Errorf("password validate error: %w", err)
	}

	return util.HashPassword(changedPassword)
}
