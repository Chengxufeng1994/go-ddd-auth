package dto

import "github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/util"

type SearchUserResponseDto struct {
	Users      []*UserDto       `json:"users"`
	Pagination *util.Pagination `join:"pagination"`
}
