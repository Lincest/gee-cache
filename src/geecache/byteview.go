package geecache

/**
    gee_cache
    @author: roccoshi
    @desc: ByteView: readonly structure to save cache value
**/

// ByteView holds an immutable view of bytes
type ByteView struct {
	b []byte
}

// Len view's size
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice copy of data
// in order to let origin ByteView immutable
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

// String makes bytes data as string
func (v ByteView) String() string {
	return string(v.b)
}
