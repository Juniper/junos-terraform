package netconf

type Client interface {
	Close() error
	SendCommit() error
	MarshalConfig(obj interface{}) error
	SendDirectTransaction(obj interface{}, commit bool) error
	SendUpdate(id string, diff string, commit bool) error
}
