// +build noasm !amd64,!arm64

package p751

import (
	. "v2ray.com/core/external/github.com/cloudflare/sidh/internal/arith"
	. "v2ray.com/core/external/github.com/cloudflare/sidh/internal/isogeny"
)

// Compute z = x + y (mod p).
func fp751AddReduced(z, x, y *FpElement) {
	var carry uint64

	// z=x+y % p751
	for i := 0; i < NumWords; i++ {
		z[i], carry = Addc64(carry, x[i], y[i])
	}

	// z = z - p751x2
	carry = 0
	for i := 0; i < NumWords; i++ {
		z[i], carry = Subc64(carry, z[i], p751x2[i])
	}

	// z = z + p751x2
	mask := uint64(0 - carry)
	carry = 0
	for i := 0; i < NumWords; i++ {
		z[i], carry = Addc64(carry, z[i], p751x2[i]&mask)
	}
}

// Compute z = x - y (mod p).
func fp751SubReduced(z, x, y *FpElement) {
	var borrow uint64

	for i := 0; i < NumWords; i++ {
		z[i], borrow = Subc64(borrow, x[i], y[i])
	}

	mask := uint64(0 - borrow)
	borrow = 0

	for i := 0; i < NumWords; i++ {
		z[i], borrow = Addc64(borrow, z[i], p751x2[i]&mask)
	}
}

// Conditionally swaps bits in x and y in constant time.
// mask indicates bits to be swaped (set bits are swapped)
// For details see "Hackers Delight, 2.20"
//
// Implementation doesn't actually depend on a prime field.
func fp751ConditionalSwap(x, y *FpElement, mask uint8) {
	var tmp, mask64 uint64

	mask64 = 0 - uint64(mask)
	for i := 0; i < len(x); i++ {
		tmp = mask64 & (x[i] ^ y[i])
		x[i] = tmp ^ x[i]
		y[i] = tmp ^ y[i]
	}
}

// Perform Montgomery reduction: set z = x R^{-1} (mod 2*p)
// with R=2^768. Destroys the input value.
func fp751MontgomeryReduce(z *FpElement, x *FpElementX2) {
	var carry, t, u, v uint64
	var uv Uint128
	var count int

	count = 5 // number of 0 digits in the least significat part of p751 + 1

	for i := 0; i < NumWords; i++ {
		for j := 0; j < i; j++ {
			if j < (i - count + 1) {
				uv = Mul64(z[j], p751p1[i-j])
				v, carry = Addc64(0, uv.L, v)
				u, carry = Addc64(carry, uv.H, u)
				t += carry
			}
		}
		v, carry = Addc64(0, v, x[i])
		u, carry = Addc64(carry, u, 0)
		t += carry

		z[i] = v
		v = u
		u = t
		t = 0
	}

	for i := NumWords; i < 2*NumWords-1; i++ {
		if count > 0 {
			count--
		}
		for j := i - NumWords + 1; j < NumWords; j++ {
			if j < (NumWords - count) {
				uv = Mul64(z[j], p751p1[i-j])
				v, carry = Addc64(0, uv.L, v)
				u, carry = Addc64(carry, uv.H, u)
				t += carry
			}
		}
		v, carry = Addc64(0, v, x[i])
		u, carry = Addc64(carry, u, 0)

		t += carry
		z[i-NumWords] = v
		v = u
		u = t
		t = 0
	}
	v, carry = Addc64(0, v, x[2*NumWords-1])
	z[NumWords-1] = v
}

// Compute z = x * y.
func fp751Mul(z *FpElementX2, x, y *FpElement) {
	var u, v, t uint64
	var carry uint64
	var uv Uint128

	for i := uint64(0); i < NumWords; i++ {
		for j := uint64(0); j <= i; j++ {
			uv = Mul64(x[j], y[i-j])
			v, carry = Addc64(0, uv.L, v)
			u, carry = Addc64(carry, uv.H, u)
			t += carry
		}
		z[i] = v
		v = u
		u = t
		t = 0
	}

	for i := NumWords; i < (2*NumWords)-1; i++ {
		for j := i - NumWords + 1; j < NumWords; j++ {
			uv = Mul64(x[j], y[i-j])
			v, carry = Addc64(0, uv.L, v)
			u, carry = Addc64(carry, uv.H, u)
			t += carry
		}
		z[i] = v
		v = u
		u = t
		t = 0
	}
	z[2*NumWords-1] = v
}

// Compute z = x + y, without reducing mod p.
func fp751AddLazy(z, x, y *FpElement) {
	var carry uint64
	for i := 0; i < NumWords; i++ {
		z[i], carry = Addc64(carry, x[i], y[i])
	}
}

// Compute z = x + y, without reducing mod p.
func fp751X2AddLazy(z, x, y *FpElementX2) {
	var carry uint64
	for i := 0; i < 2*NumWords; i++ {
		z[i], carry = Addc64(carry, x[i], y[i])
	}
}

// Reduce a field element in [0, 2*p) to one in [0,p).
func fp751StrongReduce(x *FpElement) {
	var borrow, mask uint64
	for i := 0; i < NumWords; i++ {
		x[i], borrow = Subc64(borrow, x[i], p751[i])
	}

	// Sets all bits if borrow = 1
	mask = 0 - borrow
	borrow = 0
	for i := 0; i < NumWords; i++ {
		x[i], borrow = Addc64(borrow, x[i], p751[i]&mask)
	}
}

// Compute z = x - y, without reducing mod p.
func fp751X2SubLazy(z, x, y *FpElementX2) {
	var borrow, mask uint64
	for i := 0; i < len(z); i++ {
		z[i], borrow = Subc64(borrow, x[i], y[i])
	}

	// Sets all bits if borrow = 1
	mask = 0 - borrow
	borrow = 0
	for i := NumWords; i < len(z); i++ {
		z[i], borrow = Addc64(borrow, z[i], p751[i-NumWords]&mask)
	}
}
