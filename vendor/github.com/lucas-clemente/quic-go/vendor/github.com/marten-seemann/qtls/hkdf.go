// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qtls

// Mostly derived from golang.org/x/crypto/hkdf, but with an exposed
// Extract API.
//
// HKDF is a cryptographic key derivation function (KDF) with the goal of
// expanding limited input keying material into one or more cryptographically
// strong secret keys.
//
// RFC 5869: https://tools.ietf.org/html/rfc5869

import (
	"crypto"
	"crypto/hmac"
)

func hkdfExpand(hash crypto.Hash, prk, info []byte, l int) []byte {
	var (
		expander = hmac.New(hash.New, prk)
		res      = make([]byte, l)
		counter  = byte(1)
		prev     []byte
	)

	if l > 255*expander.Size() {
		panic("hkdf: requested too much output")
	}

	p := res
	for len(p) > 0 {
		expander.Reset()
		expander.Write(prev)
		expander.Write(info)
		expander.Write([]byte{counter})
		prev = expander.Sum(prev[:0])
		counter++
		n := copy(p, prev)
		p = p[n:]
	}

	return res
}

// HkdfExtract generates a pseudorandom key for use with Expand from an input secret and an optional independent salt.
func HkdfExtract(hash crypto.Hash, secret, salt []byte) []byte {
	if salt == nil {
		salt = make([]byte, hash.Size())
	}
	if secret == nil {
		secret = make([]byte, hash.Size())
	}
	extractor := hmac.New(hash.New, salt)
	extractor.Write(secret)
	return extractor.Sum(nil)
}
