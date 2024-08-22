package bleveengine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/aggregate"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/searchengine"
	"github.com/blevesearch/bleve/analysis/analyzer/standard"
	"github.com/blevesearch/bleve/analysis/token/keyword"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

const (
	EngineName = "bleve"
	UserIndex  = "users"
)

type BleveEngine struct {
	UserIndex bleve.Index
	Mutex     sync.RWMutex
	ready     int32
	indexSync bool
}

var _ searchengine.SearchEngineInterface = (*BleveEngine)(nil)

var keywordMapping *mapping.FieldMapping
var standardMapping *mapping.FieldMapping
var dateMapping *mapping.FieldMapping

func init() {
	keywordMapping = bleve.NewTextFieldMapping()
	keywordMapping.Analyzer = keyword.Name

	standardMapping = bleve.NewTextFieldMapping()
	standardMapping.Analyzer = standard.Name

	dateMapping = bleve.NewNumericFieldMapping()
}

func getUserIndexMapping() *mapping.IndexMappingImpl {
	userMapping := bleve.NewDocumentMapping()
	userMapping.AddFieldMappingsAt("Id", keywordMapping)
	userMapping.AddFieldMappingsAt("SuggestionsWithFullname", keywordMapping)
	userMapping.AddFieldMappingsAt("SuggestionsWithoutFullname", keywordMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("_default", userMapping)

	return indexMapping
}

func NewBleveEngine() *BleveEngine {
	return &BleveEngine{}
}

func (b *BleveEngine) IndexUser(ctx context.Context, user *aggregate.User) error {
	panic("unimplemented")
}

func (b *BleveEngine) DeleteUser(user *aggregate.User) error {
	panic("unimplemented")
}

func (b *BleveEngine) getIndexDir(indexName string) string {
	return filepath.Join("data/bleve", indexName+".bleve")
}

func (b *BleveEngine) createOrOpenIndex(indexName string, mapping *mapping.IndexMappingImpl) (bleve.Index, error) {
	indexPath := b.getIndexDir(indexName)
	if index, err := bleve.Open(indexPath); err == nil {
		return index, nil
	}

	index, err := bleve.NewUsing(indexPath, mapping, "scorch", "scorch", map[string]any{
		"forceSegmentType":    "zap",
		"forceSegmentVersion": 15,
	})
	if err != nil {
		return nil, err
	}
	return index, nil
}

func (b *BleveEngine) openIndexes() error {
	if atomic.LoadInt32(&b.ready) != 0 {
		return fmt.Errorf("bleveengie already start")
	}

	var err error

	b.UserIndex, err = b.createOrOpenIndex(UserIndex, getUserIndexMapping())
	if err != nil {
		return fmt.Errorf("bleveengie create user index")
	}

	atomic.StoreInt32(&b.ready, 1)
	return nil
}

func (b *BleveEngine) closeIndexes() error {
	if err := b.UserIndex.Close(); err != nil {
		return err
	}
	atomic.StoreInt32(&b.ready, 0)
	return nil
}

func (b *BleveEngine) Start() error {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	fmt.Println("EXPERIMENTAL: Starting Bleve")

	return b.openIndexes()
}

// Stop implements searchengine.SearchEngineInterface.
func (b *BleveEngine) Stop() error {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	fmt.Println("EXPERIMENTAL: Stopping Bleve")

	return b.closeIndexes()
}

func (b *BleveEngine) RefreshIndexes(_ context.Context) error {
	return nil
}

func (b *BleveEngine) GetVersion() int {
	return 0
}

func (b *BleveEngine) GetFullVersion() string {
	return "0"
}

func (b *BleveEngine) GetPlugins() []string {
	return []string{}
}

func (b *BleveEngine) GetName() string {
	return EngineName
}

func (b *BleveEngine) deleteIndexes() error {
	if err := os.RemoveAll(b.getIndexDir(UserIndex)); err != nil {
		return fmt.Errorf("Bleveengine purge user index: %w", err)
	}
	return nil
}

func (b *BleveEngine) PurgeIndexes(ctx context.Context) error {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	fmt.Println("PurgeIndexes Bleve")
	if err := b.closeIndexes(); err != nil {
		return err
	}

	if err := b.deleteIndexes(); err != nil {
		return err
	}

	return b.openIndexes()
}

func (b *BleveEngine) DataRetentionDeleteIndexes(ctx context.Context, cutoff time.Time) error {
	return nil
}

func (b *BleveEngine) IsAutocompletionEnabled() bool {
	// TODO: configuration
	return true
}

func (b *BleveEngine) IsIndexingEnabled() bool {
	// TODO: configuration
	return true
}

func (b *BleveEngine) IsSearchEnabled() bool {
	// TODO: configuration
	return true
}

func (b *BleveEngine) IsActive() bool {
	return atomic.LoadInt32(&b.ready) == 1
}

func (b *BleveEngine) IsEnabled() bool {
	return b.IsIndexingSync()
}

func (b *BleveEngine) IsIndexingSync() bool {
	return b.indexSync
}
