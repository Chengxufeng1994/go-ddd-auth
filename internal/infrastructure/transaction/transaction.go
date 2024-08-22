package transaction

import (
	"context"
	"fmt"
	"sync"

	"gorm.io/gorm"
)

type DBFactory interface {
	GetDB(ctx context.Context) *gorm.DB
}

type defaultDBFactory struct {
	db *gorm.DB
	mu sync.Mutex
}

func NewDefaultDBFactory(db *gorm.DB) *defaultDBFactory {
	return &defaultDBFactory{db: db}
}

func (f *defaultDBFactory) GetDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(ctxWithTrx).(*gorm.DB)
	if ok {
		return tx
	}

	return f.db
}

type TransactionPropagation int

const (
	PropagationRequired TransactionPropagation = iota
	PropagationRequiresNew
	PropagationNested
	PropagationNever
)

func defaultPropagation() TransactionPropagation {
	return PropagationRequired
}

func isPropagationSupport(propagation TransactionPropagation) bool {
	return propagation >= PropagationRequired && propagation <= PropagationNever
}

type TransactionManager interface {
	Do(ctx context.Context, bizFn func(txCtx context.Context) error, propagations ...TransactionPropagation) error
}

type trxKey string

var ctxWithTrx = trxKey("tx")

type GormTransactionManager struct {
	factory DBFactory
}

var _ TransactionManager = (*GormTransactionManager)(nil)

func NewGormTransactionManager(factory DBFactory) *GormTransactionManager {
	return &GormTransactionManager{factory: factory}
}

func (m *GormTransactionManager) Do(ctx context.Context, bizFn func(txCtx context.Context) error, propagations ...TransactionPropagation) error {
	propagation := defaultPropagation()
	if len(propagations) > 0 && isPropagationSupport(propagations[0]) {
		propagation = propagations[0]
	}

	switch propagation {
	case PropagationNever:
		return m.withNeverPropagation(ctx, bizFn)
	case PropagationNested:
		return m.withNestedPropagation(ctx, bizFn)
	case PropagationRequired:
		return m.withRequiredPropagation(ctx, bizFn)
	case PropagationRequiresNew:
		return m.withRequiresNewPropagation(ctx, bizFn)
	}
	panic("not support propagation")
}

func (m *GormTransactionManager) withRequiredPropagation(ctx context.Context, bizFn func(txCtx context.Context) error) error {
	db := m.factory.GetDB(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, ctxWithTrx, tx)
		return bizFn(txCtx)
	})
}

func (m *GormTransactionManager) withRequiresNewPropagation(ctx context.Context, bizFn func(txCtx context.Context) error) error {
	// Suspend the existing transaction if any
	db := m.factory.GetDB(ctx).Session(&gorm.Session{NewDB: true})
	return db.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, ctxWithTrx, tx)
		return bizFn(txCtx)
	})
}

func (m *GormTransactionManager) withNestedPropagation(ctx context.Context, bizFn func(txCtx context.Context) error) error {
	db := m.factory.GetDB(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		// Create a savepoint
		if err := tx.Exec("SAVEPOINT sp").Error; err != nil {
			return err
		}
		txCtx := context.WithValue(ctx, ctxWithTrx, tx)
		if err := bizFn(txCtx); err != nil {
			// Rollback to the savepoint in case of error
			tx.Exec("ROLLBACK TO SAVEPOINT sp")
			return err
		}
		// Release the savepoint
		return tx.Exec("RELEASE SAVEPOINT sp").Error
	})
}

func (m *GormTransactionManager) withNeverPropagation(ctx context.Context, bizFn func(txCtx context.Context) error) error {
	tx := ctx.Value(ctxWithTrx)
	if tx != nil {
		return fmt.Errorf("transaction not allowed")
	}
	return bizFn(ctx)
}
