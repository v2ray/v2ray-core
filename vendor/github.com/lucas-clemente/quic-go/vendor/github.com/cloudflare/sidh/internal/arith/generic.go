// +build noasm !amd64

package internal

// helper used for Uint128 representation
type Uint128 struct {
	H, L uint64
}

// Adds 2 64bit digits in constant time.
// Returns result and carry (1 or 0)
func Addc64(cin, a, b uint64) (ret, cout uint64) {
	t := a + cin
	ret = b + t
	cout = ((a & b) | ((a | b) & (^ret))) >> 63
	return
}

// Substracts 2 64bit digits in constant time.
// Returns result and borrow (1 or 0)
func Subc64(bIn, a, b uint64) (ret, bOut uint64) {
	var tmp1 = a - b
	// Set bOut if bIn!=0 and tmp1==0 in constant time
	bOut = bIn & (1 ^ ((tmp1 | uint64(0-tmp1)) >> 63))
	// Constant time check if x<y
	bOut |= (a ^ ((a ^ b) | (uint64(a-b) ^ b))) >> 63
	ret = tmp1 - bIn
	return
}

// Multiplies 2 64bit digits in constant time
func Mul64(a, b uint64) (res Uint128) {
	var al, bl, ah, bh, albl, albh, ahbl, ahbh uint64
	var res1, res2, res3 uint64
	var carry, maskL, maskH, temp uint64

	maskL = (^maskL) >> 32
	maskH = ^maskL

	al = a & maskL
	ah = a >> 32
	bl = b & maskL
	bh = b >> 32

	albl = al * bl
	albh = al * bh
	ahbl = ah * bl
	ahbh = ah * bh
	res.L = albl & maskL

	res1 = albl >> 32
	res2 = ahbl & maskL
	res3 = albh & maskL
	temp = res1 + res2 + res3
	carry = temp >> 32
	res.L ^= temp << 32

	res1 = ahbl >> 32
	res2 = albh >> 32
	res3 = ahbh & maskL
	temp = res1 + res2 + res3 + carry
	res.H = temp & maskL
	carry = temp & maskH
	res.H ^= (ahbh & maskH) + carry
	return
}
