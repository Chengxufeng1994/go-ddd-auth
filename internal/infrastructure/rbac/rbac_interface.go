package rbac

type RBACAdapter interface {
	VerifyPermission(subject string, object string, action string) (bool, error)
}
