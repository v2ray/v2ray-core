package core

// User account that is used for connection to a VPoint
type VUser struct {
	id VID // The ID of this VUser.
}

// The next VPoint server in the connection chain.
type VNext struct {
	ServerAddress string  // Address of VNext server, in the form of "IP:Port"
	User          []VUser // User accounts for accessing VNext.
}

// The config for VPoint server.
type VConfig struct {
	Port           uint16 // Port of this VPoint server.
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
