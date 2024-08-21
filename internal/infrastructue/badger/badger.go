package badger

import (
	"github.com/dgraph-io/badger/v4"
)

type BadgerOption func(badger.Options) badger.Options

func NewBadger(path string, badgerOpts ...BadgerOption) (*badger.DB, error) {
	opt := badger.DefaultOptions(path)

	for _, badgerOpt := range badgerOpts {
		opt = badgerOpt(opt)
	}

	return badger.Open(opt)
}

func WithInMemory(opt badger.Options) badger.Options {
	return opt.WithInMemory(true)
}

func WithPath(path string) BadgerOption {
	return func(opt badger.Options) badger.Options {
		return opt.WithDir(path).WithValueDir(path)
	}
}
