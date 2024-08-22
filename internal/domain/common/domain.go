package common

type IDirty interface {
	Dirty()
	UnDirty()
	IsDirty() bool
}

type ISnapshot interface {
	Attach()
	Detach()
}
