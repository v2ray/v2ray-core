package stack

// TwoBytes is a [8]byte which is always allocated on stack.
//
//go:notinheap
type TwoBytes [2]byte

// EightBytes is a [8]byte which is always allocated on stack.
//
//go:notinheap
type EightBytes [8]byte
