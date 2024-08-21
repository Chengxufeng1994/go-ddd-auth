package service

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/entity"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/repository/facade"
)

type UserDomainService struct {
	repository facade.UserRepository
}

func NewUserDomainService(
	repository facade.UserRepository,
) *UserDomainService {
	return &UserDomainService{
		repository: repository,
	}
}

func (s *UserDomainService) CreateUser(ctx context.Context, user *entity.User) error {
	user.Create()
	return s.repository.SaveUser(ctx, user)
}

func (s *UserDomainService) UpdateUser(ctx context.Context, newUser *entity.User) error {
	return s.repository.UpdateUser(ctx, newUser.ID, func(oldUser *entity.User) (*entity.User, error) {
		newUser.ID = oldUser.ID
		if newUser.Username != "" {
			newUser.UpdateUsername(oldUser.Username)
		}
		if newUser.Password == "" {
			newUser.UpdatePassword(oldUser.Password)
		}
		if newUser.Role == nil {
			newUser.UpdateRole(oldUser.Role)
		}
		newUser.Update()

		return newUser, nil
	})
}

func (s *UserDomainService) DeleteUserByID(ctx context.Context, id int) error {
	return s.repository.DeleteUserByID(ctx, id)
}

func (s *UserDomainService) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	return s.repository.GetUserByID(ctx, id)
}

func (s *UserDomainService) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	return s.repository.GetUserByUsername(ctx, username)
}

func (s *UserDomainService) VerifyUser(ctx context.Context, username, password string) (*entity.User, error) {
	user, err := s.repository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if err := user.ValidatePassword(password); err != nil {
		return nil, err
	}
	return user, nil
}
