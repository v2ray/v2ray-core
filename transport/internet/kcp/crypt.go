package kcp

type BlockCrypt interface {
	// Encrypt encrypts the whole block in src into dst.
	// Dst and src may point at the same memory.
	Encrypt(dst, src []byte)

	// Decrypt decrypts the whole block in src into dst.
	// Dst and src may point at the same memory.
	Decrypt(dst, src []byte)
}

// None Encryption
type NoneBlockCrypt struct {
	xortbl []byte
}

func NewNoneBlockCrypt(key []byte) (BlockCrypt, error) {
	return new(NoneBlockCrypt), nil
}

func (c *NoneBlockCrypt) Encrypt(dst, src []byte) {}
func (c *NoneBlockCrypt) Decrypt(dst, src []byte) {}
