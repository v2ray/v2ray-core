// +build amd64,!noasm arm64,!noasm

package p751

import (
	. "v2ray.com/core/external/github.com/cloudflare/sidh/internal/isogeny"
)

// If choice = 0, leave x,y unchanged. If choice = 1, set x,y = y,x.
// If choice is neither 0 nor 1 then behaviour is undefined.
// This function executes in constant time.
//go:noescape
func fp751ConditionalSwap(x, y *FpElement, choice uint8)

// Compute z = x + y (mod p).
//go:noescape
func fp751AddReduced(z, x, y *FpElement)

// Compute z = x - y (mod p).
//go:noescape
func fp751SubReduced(z, x, y *FpElement)

// Compute z = x + y, without reducing mod p.
//go:noescape
func fp751AddLazy(z, x, y *FpElement)

// Compute z = x + y, without reducing mod p.
//go:noescape
func fp751X2AddLazy(z, x, y *FpElementX2)

// Compute z = x - y, without reducing mod p.
//go:noescape
func fp751X2SubLazy(z, x, y *FpElementX2)

// Compute z = x * y.
//go:noescape
func fp751Mul(z *FpElementX2, x, y *FpElement)

// Compute Montgomery reduction: set z = x * R^{-1} (mod 2*p).
// It may destroy the input value.
//go:noescape
func fp751MontgomeryReduce(z *FpElement, x *FpElementX2)

// Reduce a field element in [0, 2*p) to one in [0,p).
//go:noescape
func fp751StrongReduce(x *FpElement)
