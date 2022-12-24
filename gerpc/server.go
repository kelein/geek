package gerpc

import (
	"encoding/json"
	"io"
	"net"
	"reflect"
	"sync"

	"github.com/pkg/errors"

	"geek/gerpc/codec"
	"geek/glog"
)

// OriginMagicNum for protocol
const OriginMagicNum = 0x3bef5c

var invalidRequest = struct{}{}

// DefaultOption for custom server settings
var DefaultOption = &Option{
	MagicNumber: OriginMagicNum,
	CodecKind:   codec.GOB,
}

var defaultServer = NewServer()

// Option for custom server settings
type Option struct {
	MagicNumber int
	CodecKind   codec.Kind
}

type request struct {
	header *codec.Header
	argv   reflect.Value
	replyv reflect.Value
}

// Server stands for RPC server
type Server struct{}

// NewServer create a Server instance
func NewServer() *Server {
	return &Server{}
}

// Accept holds connections and serves requests with default server
func Accept(l net.Listener) {
	defaultServer.Accept(l)
}

// Accept holds connections and serves requests
func (s *Server) Accept(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			glog.Errorf("rpc server accept error: %v", err)
			return
		}
		go s.ServeConn(conn)
	}
}

// ServeConn run server on the given connection
func (s *Server) ServeConn(conn net.Conn) {
	defer conn.Close()

	opt := &Option{}
	if err := json.NewDecoder(conn).Decode(opt); err != nil {
		glog.Errorf("rpc server options error: %v", err)
		return
	}

	if opt.MagicNumber != OriginMagicNum {
		glog.Errorf("rpc server invalid magic number: %x", opt.MagicNumber)
		return
	}

	fn := codec.NewCodecFuncMap[opt.CodecKind]
	if fn == nil {
		glog.Errorf("rpc server invalid codec: %v", opt.CodecKind)
		return
	}
	s.serveCodec(fn(conn))
}

func (s *Server) serveCodec(c codec.Codec) {
	sending := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	for {
		req, err := s.readRequest(c)
		if err != nil {
			if req == nil {
				break
			}
			req.header.Error = err.Error()
			s.sendResponse(c, req.header, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go s.handleRequest(c, req, sending, wg)
	}

	wg.Wait()
	c.Close()
}

func (s *Server) readRequestHeader(c codec.Codec) (*codec.Header, error) {
	h := &codec.Header{}
	if err := c.ReadHeader(h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			glog.Errorf("rpc server read header error: %v", err)
		}
		return nil, err
	}
	return h, nil
}

func (s *Server) readRequest(c codec.Codec) (*request, error) {
	h, err := s.readRequestHeader(c)
	if err != nil {
		return nil, errors.Wrap(err, "read header error")
	}
	req := &request{header: h}

	// TODO: parse request args
	req.argv = reflect.New(reflect.TypeOf(""))

	if err := c.ReadBody(req.argv.Interface()); err != nil {
		glog.Errorf("rpc server read argv error: %v", err)
		return nil, errors.Wrap(err, "read argv error")
	}
	return req, nil
}

func (s *Server) sendResponse(c codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := c.Write(h, body); err != nil {
		glog.Errorf("rpc server write body error: %v", err)
	}
}

func (s *Server) handleRequest(c codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	glog.Info(req.header, req.argv.Elem())
	// TODO: parse request reply
	req.replyv = reflect.ValueOf(req.header.Seq)
	body := req.replyv.Interface()
	s.sendResponse(c, req.header, body, sending)
}
