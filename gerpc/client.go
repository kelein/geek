package gerpc

import (
	"errors"
	"io"
	"sync"

	"geek/gerpc/codec"
)

// Call stands for an active RPC
type Call struct {
	Seq           uint64
	ServiceMethod string
	Args          interface{}
	Reply         interface{}
	Error         error
	Done          chan *Call
}

func (call *Call) done() {
	call.Done <- call
}

// Client stands for a RPC client
type Client struct {
	cc       codec.Codec
	opt      *Option
	sending  sync.Mutex
	header   codec.Header
	mu       sync.Mutex
	seq      uint64
	pending  map[uint64]*Call
	closing  bool
	shutdown bool
}

var _ io.Closer = (*Client)(nil)

// ErrShutdown client shutdown error
var ErrShutdown = errors.New("connection on closing")

// Close the connection
func (cli *Client) Close() error {
	cli.mu.Lock()
	defer cli.mu.Unlock()
	if cli.closing {
		return ErrShutdown
	}
	cli.closing = true
	return cli.cc.Close()
}

// IsAvailable checks whether the client works
func (cli *Client) IsAvailable() bool {
	cli.mu.Lock()
	defer cli.mu.Unlock()
	return !cli.shutdown && !cli.closing
}
