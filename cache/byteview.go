package cache

type ByteView struct {
	bytes []byte
}

func NewByteView(b []byte) ByteView {
	return ByteView{bytes: cloneBytes(b)}
}

func (v ByteView) Len() int {
	return len(v.bytes)
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.bytes)
}

func (v ByteView) String() string {
	return string(v.bytes)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
