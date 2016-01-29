package crypto

type Authenticator interface {
	AuthBytes() int
	Authenticate(auth []byte, data []byte) []byte
}
