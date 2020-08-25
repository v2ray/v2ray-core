// +build amd64,!noasm

// Sets capabilities flags for x86 according to information received from
// CPUID. It was written in accordance with
// "Intel® 64 and IA-32 Architectures Developer's Manual: Vol. 2A".
// https://www.intel.com/content/www/us/en/architecture-and-technology/64-ia-32-architectures-software-developer-vol-2a-manual.html

package utils

// Performs CPUID and returns values of registers
// go:nosplit
func cpuid(eaxArg, ecxArg uint32) (eax, ebx, ecx, edx uint32)

// Returns true in case bit 'n' in 'bits' is set, otherwise false
func bitn(bits uint32, n uint8) bool {
	return (bits>>n)&1 == 1
}

func init() {
	// CPUID returns max possible input that can be requested
	max, _, _, _ := cpuid(0, 0)
	if max < 7 {
		return
	}

	_, ebx, _, _ := cpuid(7, 0)
	X86.HasBMI2 = bitn(ebx, 8)
	X86.HasADX = bitn(ebx, 19)
}
