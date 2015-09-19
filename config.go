package core

// User is the user account that is used for connection to a Point
type User struct {
	Id ID `json:"id"` // The ID of this User.
}

type ConnectionConfig interface {
  Protocol() string
  Content() []byte
}

type PointConfig interface {
  Port() uint16
  InboundConfig() ConnectionConfig
  OutboundConfig() ConnectionConfig
}
