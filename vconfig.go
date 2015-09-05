package core

type VUser struct {
  id VID
}

type VConfig struct {
  RunAs VUser
  Port uint16
  AllowedClients []VUser
  AllowedProtocol string
}

type VConfigMarshaller interface {
  Marshal(config VConfig) ([]byte, error)
}

type VConfigUnmarshaller interface {
  Unmarshal(data []byte) (VConfig, error)
}
