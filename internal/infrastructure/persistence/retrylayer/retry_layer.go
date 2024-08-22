package retrylayer

import (
	"context"
	"time"

	rolerepository "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/repository/facade"
	"github.com/pkg/errors"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/entity"
	userrepository "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/repository/facade"
)

type RetryLayer struct {
	userrepository.UserRepository
	rolerepository.RoleRepository
}

func NewRetryLayer(
	userRepository userrepository.UserRepository,
) *RetryLayer {
	return &RetryLayer{
		UserRepository: userRepository,
	}
}

func (s *RetryLayer) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	tries := 0
	for {
		result, err := s.UserRepository.GetUserByID(ctx, id)
		if err == nil {
			return result, nil
		}
		tries++
		if tries >= 3 {
			err = errors.Wrap(err, "giving up after 3 consecutive repeatable transaction failures")
			return result, err
		}
		time.Sleep(100 * time.Millisecond)
	}
}
