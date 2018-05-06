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
