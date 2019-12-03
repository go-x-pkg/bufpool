package bufpool

import (
	"bytes"
	"sync"
)

var pool *sync.Pool = nil

type Buf struct {
	bytes.Buffer
	pool *sync.Pool
}

/* impl Unrefer for Buf */
func (b *Buf) Unref() {
	b.Release()
}

func (b *Buf) Release() {
	b.pool.Put(b)
}

func newBuf(pool *sync.Pool) *Buf {
	b := new(Buf)
	b.pool = pool

	return b
}

func NewBuf() *Buf {
	b := pool.Get().(*Buf)
	b.Reset()

	return b
}

func init() {
	pool = &sync.Pool{
		New: func() interface{} {
			buf := newBuf(pool)
			return buf
		},
	}
}
