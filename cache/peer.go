package cache

import "geek/cache/proto"

// PeerPicker of abstract
type PeerPicker interface {
	Pick(key string) (peer PeerGetter, ok bool)
}

// PeerGetter abstract Peer
// type PeerGetter interface {
// 	Get(group, key string) ([]byte, error)
// }

// PeerGetter abstract Peer with protobuf
type PeerGetter interface {
	Get(in *proto.Request, out *proto.Response) error
}
