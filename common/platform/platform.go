// Package platform provides platform specific functionalities.
package platform

type environment interface {
	ExpandEnv(s string) string
}

func ExpandEnv(s string) string {
	return environmentInstance.ExpandEnv(s)
}
