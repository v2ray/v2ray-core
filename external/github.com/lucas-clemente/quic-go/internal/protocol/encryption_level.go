package protocol

// EncryptionLevel is the encryption level
// Default value is Unencrypted
type EncryptionLevel int

const (
	// EncryptionUnspecified is a not specified encryption level
	EncryptionUnspecified EncryptionLevel = iota
	// EncryptionInitial is the Initial encryption level
	EncryptionInitial
	// EncryptionHandshake is the Handshake encryption level
	EncryptionHandshake
	// Encryption1RTT is the 1-RTT encryption level
	Encryption1RTT
)

func (e EncryptionLevel) String() string {
	switch e {
	case EncryptionInitial:
		return "Initial"
	case EncryptionHandshake:
		return "Handshake"
	case Encryption1RTT:
		return "1-RTT"
	}
	return "unknown"
}
