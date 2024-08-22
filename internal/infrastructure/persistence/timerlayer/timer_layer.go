package timerlayer

import (
	"context"
	"fmt"
	"time"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/aggregate"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/repository/facade"
	identity_access_mgmt "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/repository/facade"
)

type TimerLayer struct {
	identity_access_mgmt.UserRepository
	identity_access_mgmt.RoleRepository
}

func NewTimerLayer(userRepository facade.UserRepository) *TimerLayer {
	return &TimerLayer{
		UserRepository: userRepository,
	}
}

func (s *TimerLayer) GetUserByID(ctx context.Context, id int) (*aggregate.User, error) {
	start := time.Now()
	user, err := s.UserRepository.GetUserByID(ctx, id)
	elapsed := float64(time.Since(start)) / float64(time.Second)
	fmt.Printf("GetUserByID time %f secons.\n", elapsed)
	return user, err
}
