package badger

import (
	"encoding/json"
	"time"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/cache"
	"github.com/dgraph-io/badger/v4"
)

type BadgerOption func(badger.Options) badger.Options

func WithInMemory(opt badger.Options) badger.Options {
	return opt.WithInMemory(true)
}

func WithPath(path string) BadgerOption {
	return func(opt badger.Options) badger.Options {
		return opt.WithDir(path).WithValueDir(path)
	}
}

type Badger struct {
	db            *badger.DB
	defaultExpiry time.Duration
}

var _ cache.Cache = (*Badger)(nil)

func NewBadger(path string, badgerOpts ...BadgerOption) (*Badger, error) {
	opt := badger.DefaultOptions(path)

	for _, badgerOpt := range badgerOpts {
		opt = badgerOpt(opt)
	}

	badgerDB, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	return &Badger{
		db:            badgerDB,
		defaultExpiry: 10 * time.Second}, nil

}

// Get implements cache.Cache.
func (b *Badger) Get(key string, value any) error {
	var valCopy []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
	})
	if err != nil {
		return err
	}

	err = json.Unmarshal(valCopy, &value)
	if err != nil {
		return err
	}

	return nil
}

// GetMulti implements cache.Cache.
func (b *Badger) GetMulti(keys []string, values []any) []error {
	panic("unimplemented")
}

// Name implements cache.Cache.
func (b *Badger) Name() string {
	return "badger"
}

// Purge implements cache.Cache.
func (b *Badger) Purge() error {
	panic("unimplemented")
}

// Remove implements cache.Cache.
func (b *Badger) Remove(key string) error {
	panic("unimplemented")
}

// RemoveMulti implements cache.Cache.
func (b *Badger) RemoveMulti(keys []string) error {
	panic("unimplemented")
}

// SetWithDefaultExpiry implements cache.Cache.
func (b *Badger) SetWithDefaultExpiry(key string, value any) error {
	byt, _ := json.Marshal(value)
	return b.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), byt).WithTTL(b.defaultExpiry)
		return txn.SetEntry(e)
	})
}

// SetWithExpiry implements cache.Cache.
func (b *Badger) SetWithExpiry(key string, value any, ttl time.Duration) error {
	byt, _ := json.Marshal(value)
	return b.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), byt).WithTTL(ttl)
		return txn.SetEntry(e)
	})
}

func (b *Badger) Close() error {
	return b.db.Close()
}
