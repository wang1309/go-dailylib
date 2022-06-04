package simple_bp

type ByteBuffer struct {
	B []byte
}

// Reset makes ByteBuffer.B empty.
func (b *ByteBuffer) Reset() {
	b.B = b.B[:0]
}

