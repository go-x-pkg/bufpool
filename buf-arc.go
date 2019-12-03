package bufpool

import (
	"bytes"
	"sync"
)

var poolArc *sync.Pool = nil

// Atomically Reference Counted buffer
type BufArc struct {
	bytes.Buffer
	sync.Mutex

	rc   int64 // reference counter
	pool *sync.Pool
}

/* impl ReferDelta for BufArc */
func (b *BufArc) RefDelta(delta int64) {
	b.Lock()
	defer b.Unlock()
	b.rc += delta
}

func (b *BufArc) Ref() { b.RefDelta(1) }

/* impl Unrefer for BufArc */
func (b *BufArc) Unref() {
	b.Lock()
	defer b.Unlock()

	b.rc--

	if b.rc >= 1 {
		return
	}

	b.rc = 0
	b.release()
}

func (b *BufArc) release() {
	b.pool.Put(b)
}

func newBufArc(pool *sync.Pool) *BufArc {
	b := new(BufArc)
	b.pool = pool

	return b
}

func NewBufArc() *BufArc {
	b := poolArc.Get().(*BufArc)
	b.Reset()

	return b
}

func init() {
	poolArc = &sync.Pool{
		New: func() interface{} {
			buf := newBufArc(poolArc)
			return buf
		},
	}
}
