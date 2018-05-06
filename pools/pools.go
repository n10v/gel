package pools

import (
	"bytes"
	"sync"
)

// bbPool is a pool of *bytes.Buffer.
var bbPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

// GetBytesBuffer returns empty *bytes.Buffer.
func GetBytesBuffer() *bytes.Buffer {
	return bbPool.Get().(*bytes.Buffer)
}

// PutBytesBuffer resets buf and puts to pool.
func PutBytesBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bbPool.Put(buf)
}

// bsPool is a pool of byte slices.
var bsPool = sync.Pool{
	New: func() interface{} { return nil },
}

// GetByteSlice returns []byte with len == size.
func GetByteSlice(size int) []byte {
	fromPool := bsPool.Get()
	if fromPool == nil {
		return make([]byte, size)
	}
	bs := fromPool.([]byte)
	if cap(bs) < size {
		bs = make([]byte, size)
	}
	return bs[0:size]
}

// PutByteSlice puts b to pool.
func PutByteSlice(b []byte) {
	bsPool.Put(b)
}
