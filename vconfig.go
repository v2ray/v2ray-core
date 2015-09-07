package core

// User account that is used for connection to a VPoint
type VUser struct {
	// The ID of this VUser. This ID is served as an access token.
	// It is not necessary to be permanent.
	id VID
}

// The next VPoint server in the connection chain.
type VNext struct {
	// Address of VNext server, in the form of "IP:Port"
	ServerAddress string
	// User accounts for accessing VNext.
	User []VUser
}

// The config for VPoint server.
type VConfig struct {
	// Port of this VPoint server.
	Port           uint16
	AllowedClients []VUser
	ClientProtocol string
	VNextList      []VNext
}

type VConfigMarshaller interface {
	Marshal(config VConfig) ([]byte, error)
}

type VConfigUnmarshaller interface {
	Unmarshal(data []byte) (VConfig, error)
}
