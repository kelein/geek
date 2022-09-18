package cache

// ByteView stands for immutable bytes
type ByteView struct {
	b []byte
}

// Len returns byte view length
func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

// ByteSlice return a data copy of byte slice
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
