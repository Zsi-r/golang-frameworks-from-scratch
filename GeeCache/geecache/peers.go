package geecache

// PeerPicker chooses a `PeerGetter` accordig to the input key
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter get cache value from corresponding group, i.e. HTTP client
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
