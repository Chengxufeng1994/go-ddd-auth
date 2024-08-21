package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/token/entity"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/token/repository/facade"
	"github.com/dgraph-io/badger/v4"
)

type BadgerTokenRepository struct {
	db *badger.DB
}

var _ facade.TokenRepository = (*BadgerTokenRepository)(nil)

func NewTokenRepository(db *badger.DB) *BadgerTokenRepository {
	return &BadgerTokenRepository{
		db: db,
	}
}

func buildAccessTokenKey(userID int) []byte {
	key := fmt.Sprintf("access_token:user_id:%d", userID)
	return []byte(key)
}

func buildRefreshTokenKey(userID int) []byte {
	key := fmt.Sprintf("refresh_token:user_id:%d", userID)
	return []byte(key)
}

func (r *BadgerTokenRepository) SaveAccessToken(ctx context.Context, authDetails *entity.AuthDetails, accessToken *entity.AuthToken, ttl time.Duration) error {
	byt, _ := json.Marshal(accessToken)
	return r.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(buildAccessTokenKey(authDetails.User.UserID), byt).WithTTL(ttl)
		return txn.SetEntry(e)
	})
}

func (r *BadgerTokenRepository) GetAccessToken(ctx context.Context, authDetails *entity.AuthDetails) (*entity.AuthToken, error) {
	var valCopy []byte
	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(buildAccessTokenKey(authDetails.User.UserID))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	var authToken entity.AuthToken
	err = json.Unmarshal(valCopy, &authToken)
	if err != nil {
		return nil, err
	}

	return &authToken, nil
}

func (r *BadgerTokenRepository) SaveRefreshToken(ctx context.Context, authDetails *entity.AuthDetails, refreshToken *entity.AuthToken, ttl time.Duration) error {
	byt, _ := json.Marshal(refreshToken)
	return r.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(buildRefreshTokenKey(authDetails.User.UserID), byt).WithTTL(ttl)
		return txn.SetEntry(e)
	})
}
