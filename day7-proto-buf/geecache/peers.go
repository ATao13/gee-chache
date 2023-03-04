package geecache

import pb "github.com/ATao13/gee-chache/day7-proto-buf/geecache/geecachepb"

type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}

type PeerPicker interface {
	PickPeer(string) (PeerGetter, bool)
}
