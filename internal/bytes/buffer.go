// Package bytes implements performant bytes implementations.
package bytes

import (
	"strconv"
	"unicode/utf8"
)

// Buffer wraps a byte slice, providing continence functions.
type Buffer struct {
	b []byte
}

// NewBuffer returns a buffer.
func NewBuffer(size int) *Buffer {
	return &Buffer{b: make([]byte, 0, size)}
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

// WriteByte writes a single byte to the Buffer.
func (b *Buffer) WriteByte(v byte) {
	b.b = append(b.b, v)
}

// WriteString writes a string to the Buffer.
func (b *Buffer) WriteString(s string) {
	b.b = append(b.b, s...)
}

// WriteRune writes a rune to the Buffer.
func (b *Buffer) WriteRune(r rune) {
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	b.b = append(b.b, buf[:n]...)
}

// Write implements io.Writer.
func (b *Buffer) Write(bs []byte) {
	b.b = append(b.b, bs...)
}

// Bytes returns a mutable reference to the underlying byte slice.
func (b *Buffer) Bytes() []byte {
	return b.b
}

// Len returns the length of the buffer.
func (b *Buffer) Len() int {
	return len(b.b)
}

// Reset resets the underlying byte slice. Subsequent writes re-use the slice's
// backing array.
func (b *Buffer) Reset() {
	b.b = b.b[:0]
}
