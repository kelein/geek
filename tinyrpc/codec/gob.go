package codec

import (
	"bufio"
	"encoding/gob"
	"io"

	"geek/glog"

	"github.com/pkg/errors"
)

var _ Codec = (*GobCodec)(nil)

// GobCodec Gob Kind Codec
type GobCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

// NewGobCodec create a GobCodec instance
func NewGobCodec(conn io.ReadWriteCloser) Codec {
	return &GobCodec{
		conn: conn,
		buf:  bufio.NewWriter(conn),
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(conn),
	}
}

// ReadHeader read header by decoder
func (g *GobCodec) ReadHeader(h *Header) error {
	return g.dec.Decode(h)
}

// ReadBody read request body by decoder
func (g *GobCodec) ReadBody(body interface{}) error {
	return g.dec.Decode(body)
}

// Close closes the server connection
func (g *GobCodec) Close() error {
	return g.conn.Close()
}

func (g *GobCodec) Write(h *Header, body interface{}) error {
	defer func() {
		if err := g.buf.Flush(); err != nil {
			g.Close()
		}
	}()

	if err := g.enc.Encode(h); err != nil {
		glog.Errorf("gob codec encoding header error: %v", err)
		return errors.Wrap(err, "gob codec encoding header error")
	}

	if err := g.enc.Encode(body); err != nil {
		glog.Errorf("gob codec encoding body error: %v", err)
		return errors.Wrap(err, "gob codec encoding body error")
	}

	return nil
}
