package netconf

type Client interface {
	Close() error
	DeleteConfig(applyGroup string, commit bool) (string, error)
	SendCommit() error
	MarshalGroup(id string, obj interface{}) error
	MarshalConfig(obj interface{}) error
	SendTransaction(id string, obj interface{}, commit bool) error
	SendDirectTransaction(obj interface{}, commit bool) error
	SendUpdate(id string, diff string, commit bool) error
}
