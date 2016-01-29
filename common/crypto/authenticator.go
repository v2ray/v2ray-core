package crypto

type Authenticator interface {
	AuthSize() int
	Authenticate(auth []byte, data []byte) []byte
}
