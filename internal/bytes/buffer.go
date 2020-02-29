package bytes

import (
	"strconv"
	"sync"
	"time"
)

// Pool is a Pool of Buffers.
type Pool struct {
	p *sync.Pool
}

// NewPool creates a new instance of Pool.
func NewPool(size int) Pool {
	return Pool{p: &sync.Pool{
		New: func() interface{} {
			return &Buffer{b: make([]byte, 0, size)}
		},
	}}
}

// Get retrieves a Buffer from the Pool, creating one if necessary.
func (p Pool) Get() *Buffer {
	buf := p.p.Get().(*Buffer)
	buf.Reset()
	return buf
}

// Put adds a Buffer to the Pool.
func (p Pool) Put(buf *Buffer) {
	p.p.Put(buf)
}

// Buffer wraps a byte slice, providing continence functions.
type Buffer struct {
	b []byte
}

// AppendInt appends an integer to the underlying Buffer.
func (b *Buffer) AppendInt(i int64) {
	b.b = strconv.AppendInt(b.b, i, 10)
}

// AppendUint appends an unsigned integer to the underlying Buffer.
func (b *Buffer) AppendUint(i uint64) {
	b.b = strconv.AppendUint(b.b, i, 10)
}

// AppendFloat appends a float to the underlying Buffer.
func (b *Buffer) AppendFloat(f float64, fmt byte, prec, bitSize int) {
	b.b = strconv.AppendFloat(b.b, f, fmt, prec, bitSize)
}

// AppendBool appends a bool to the underlying Buffer.
func (b *Buffer) AppendBool(v bool) {
	b.b = strconv.AppendBool(b.b, v)
}

// AppendTime appends a time to the underlying Buffer, in the given layout.
func (b *Buffer) AppendTime(t time.Time, layout string) {
	b.b = t.AppendFormat(b.b, layout)
}

// WriteByte writes a single byte to the Buffer.
func (b *Buffer) WriteByte(v byte) {
	b.b = append(b.b, v)
}

// WriteString writes a string to the Buffer.
func (b *Buffer) WriteString(s string) {
	b.b = append(b.b, s...)
}

// Write implements io.Writer.
func (b *Buffer) Write(bs []byte) (int, error) {
	b.b = append(b.b, bs...)

	return len(bs), nil
}

// Len returns the length of the underlying byte slice.
func (b *Buffer) Len() int {
	return len(b.b)
}

// Cap returns the capacity of the underlying byte slice.
func (b *Buffer) Cap() int {
	return cap(b.b)
}

// Bytes returns a mutable reference to the underlying byte slice.
func (b *Buffer) Bytes() []byte {
	return b.b
}

// Reset resets the underlying byte slice. Subsequent writes re-use the slice's
// backing array.
func (b *Buffer) Reset() {
	b.b = b.b[:0]
}
