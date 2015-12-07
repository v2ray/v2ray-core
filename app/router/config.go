package router

type Config interface {
	Strategy() string
	Settings() interface{}
}
