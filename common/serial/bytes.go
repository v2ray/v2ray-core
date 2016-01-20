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

func (this BytesLiteral) String() string {
	return string(this.Value())
}

func (this BytesLiteral) AllZero() bool {
  for _, b := range this {
    if b != 0 {
      return false
    }
  }
  return true
}