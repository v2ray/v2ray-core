package utils

type x86 struct {
	// Signals support for MULX which is in BMI2
	HasBMI2 bool

	// Signals support for ADX
	HasADX bool
}

var X86 x86
