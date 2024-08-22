package searchengine

import (
	"context"
	"time"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/aggregate"
)

type SearchEngineInterface interface {
	Start() error
	Stop() error
	GetFullVersion() string
	GetVersion() int
	GetPlugins() []string
	GetName() string
	// IsEnabled returns a boolean indicating whether the engine is enabled in the settings
	IsEnabled() bool
	IsActive() bool
	IsIndexingEnabled() bool
	IsSearchEnabled() bool
	IsAutocompletionEnabled() bool
	IsIndexingSync() bool
	// IndexChannel indexes a given channel. The userIDs are only populated
	// for private channels.
	IndexUser(ctx context.Context, user *aggregate.User) error
	DeleteUser(user *aggregate.User) error
	PurgeIndexes(ctx context.Context) error
	RefreshIndexes(ctx context.Context) error
	DataRetentionDeleteIndexes(ctx context.Context, cutoff time.Time) error
}
