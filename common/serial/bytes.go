package serial

type Bytes interface {
	Bytes() []byte
}

type BytesLiteral []byte

func (this BytesLiteral) Value() []byte {
	return []byte(this)
}

func (this BytesLiteral) Int64Value() int64 {
	value := this.Value()
	return int64(value[0])<<56 +
		int64(value[1])<<48 +
		int64(value[2])<<40 +
		int64(value[3])<<32 +
		int64(value[4])<<24 +
		int64(value[5])<<16 +
		int64(value[6])<<8 +
		int64(value[7])
}

// String returns a string presentation of this ByteLiteral
func (this BytesLiteral) String() string {
	return string(this.Value())
}

// All returns true if all bytes in the ByteLiteral are the same as given value.
func (this BytesLiteral) All(v byte) bool {
  for _, b := range this {
    if b != v {
      return false
    }
  }
  return true
}