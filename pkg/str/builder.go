package go_json

import (
	"github.com/Eugene-Usachev/fastbytes"
	"sync"
)

type Builder struct {
	buf []byte
}

func newBuilder() *Builder {
	return &Builder{}
}

var builderPool = sync.Pool{New: func() any {
	return newBuilder()
}}

func GetBuilder() *Builder {
	return builderPool.Get().(*Builder)
}

func NewBuilderWithSting(str string) *Builder {
	b := newBuilder()
	b.buf = fastbytes.S2B(str)
	return b
}

func NewBuilderWithBytes(buf []byte) *Builder {
	b := newBuilder()
	b.buf = buf
	return b
}

func (b *Builder) String() string {
	return fastbytes.B2S(b.buf)
}

func (b *Builder) Buffer() []byte {
	return b.buf
}

func (b *Builder) Grow(n int) {
	buf := make([]byte, len(b.buf), cap(b.buf)+n)
	copy(buf, b.buf)
	b.buf = buf
}

func (b *Builder) GrowTo(n int) {
	l := len(b.buf)
	if n < l {
		b.buf = b.buf[:n]
	}
	buf := make([]byte, len(b.buf), n-l)
	copy(buf, b.buf)
	b.buf = buf
}

func (b *Builder) Reset() {
	b.buf = make([]byte, 0, len(b.buf))
}

func (b *Builder) ReleaseWithReset() {
	b.Reset()
	builderPool.Put(b)
}

func (b *Builder) Release() {
	builderPool.Put(b)
}

func (b *Builder) Len() int {
	return len(b.buf)
}

func (b *Builder) Cap() int {
	return cap(b.buf)
}
