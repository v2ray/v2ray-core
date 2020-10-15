package trojan

import (
	"crypto/sha256"
	"encoding/hex"
	fmt "fmt"

	"v2ray.com/core/common"
	"v2ray.com/core/common/protocol"
)

// MemoryAccount is an account type converted from Account.
type MemoryAccount struct {
	Password string
	Key      []byte
}

// AsAccount implements protocol.AsAccount.
func (a *Account) AsAccount() (protocol.Account, error) {
	password := a.GetPassword()
	key := hexSha224(password)
	return &MemoryAccount{
		Password: password,
		Key:      key,
	}, nil
}

// Equals implements protocol.Account.Equals().
func (a *MemoryAccount) Equals(another protocol.Account) bool {
	if account, ok := another.(*MemoryAccount); ok {
		return a.Password == account.Password
	}
	return false
}

func hexSha224(password string) []byte {
	buf := make([]byte, 56)
	hash := sha256.New224()
	common.Must2(hash.Write([]byte(password)))
	hex.Encode(buf, hash.Sum(nil))
	return buf
}

func hexString(data []byte) string {
	str := ""
	for _, v := range data {
		str += fmt.Sprintf("%02x", v)
	}
	return str
}
