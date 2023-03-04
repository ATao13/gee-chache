package geecache

type PeerGetter interface {
	Get(group string, peer string) ([]byte, error)
}

type PeerPicker interface {
	PickPeer(string) (PeerGetter, bool)
}
