package bitmask

// Byte is a bitmask in byte.
type Byte byte

// Has returns true if this bitmask contains another bitmask.
func (b Byte) Has(bb Byte) bool {
	return (b & bb) != 0
}

func (b *Byte) Add(bb Byte) {
	*b |= bb
}

func (b *Byte) Clear(bb Byte) {
	*b &= ^bb
}

func (b *Byte) Toggle(bb Byte) {
	*b ^= bb
}
