package rbac

import (
	"time"

	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

func NewCasbinEnforcer(confPath string, policyPath string) (*casbin.SyncedEnforcer, error) {
	adapter := fileadapter.NewAdapter(policyPath)
	e, err := casbin.NewSyncedEnforcer(confPath, adapter)
	if err != nil {
		return nil, err
	}
	// e.EnableLog(true)
	e.StartAutoLoadPolicy(5 * time.Second)

	return e, nil
}

type CasbinAdapter struct {
	casbin *casbin.SyncedEnforcer
}

var _ RBACAdapter = (*CasbinAdapter)(nil)

func NewCasbinAdapter(casbin *casbin.SyncedEnforcer) *CasbinAdapter {
	return &CasbinAdapter{
		casbin: casbin,
	}
}

func (a *CasbinAdapter) VerifyPermission(subject string, object string, action string) (bool, error) {
	return a.casbin.Enforce(subject, object, action)
}
