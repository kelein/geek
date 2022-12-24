package codec

import "io"

// Header for Request
type Header struct {
	ServiceMethod string
	Seq           uint64
	Error         string
}

// Codec of abstract
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

// NewCodecFunc func to create Codec
type NewCodecFunc func(io.ReadWriteCloser) Codec

// Kind type of Codec
type Kind string

// Codec Kind
const (
	GOB  Kind = "application/gob"
	JSON Kind = "application/json"
)

// NewCodecFuncMap Codec Factory
var NewCodecFuncMap = make(map[Kind]NewCodecFunc)

func init() {
	NewCodecFuncMap[GOB] = NewGobCodec
}
