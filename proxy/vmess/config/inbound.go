package config

type Inbound interface {
	AllowedUsers() []User
}
