package netconf

type Client interface {
	Close() error
	ReadGroup(applygroup string) (string, error)
	UpdateRawConfig(applygroup string, netconfcall string, commit bool) (string, error)
	DeleteConfig(applygroup string) (string, error)
	DeleteConfigNoCommit(applygroup string) (string, error)
	SendCommit() error
	MarshalGroup(id string, obj interface{}) error
	SendTransaction(id string, obj interface{}, commit bool) error
	SendRawConfig(netconfcall string, commit bool) (string, error)
	ReadRawGroup(applygroup string) (string, error)
}
