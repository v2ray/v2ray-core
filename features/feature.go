package features

import "v2ray.com/core/common"

//go:generate errorgen

// Feature is the interface for V2Ray features. All features must implement this interface.
// All existing features have an implementation in app directory. These features can be replaced by third-party ones.
type Feature interface {
	common.HasType
	common.Runnable
}

// PrintDeprecatedFeatureWarning prints a warning for deprecated feature.
func PrintDeprecatedFeatureWarning(feature string) {
	newError("You are using a deprecated feature: " + feature + ". Please update your config file with latest configuration format, or update your client software.").WriteToLog()
}
