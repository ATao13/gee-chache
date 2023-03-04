package geecache

// PeerGetter 的 Get() 方法用于从对应 group 查找缓存值
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}
