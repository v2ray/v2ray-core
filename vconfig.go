package core

import (
	v2net "github.com/v2ray/v2ray-core/net"
)

// VUser is the user account that is used for connection to a VPoint
type VUser struct {
	Id VID // The ID of this VUser.
}

// VNext is the next VPoint server in the connection chain.
type VNext struct {
	Address v2net.VAddress // Address of VNext server
	Users   []VUser        // User accounts for accessing VNext.
}

// VConfig is the config for VPoint server.
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
