package kcp

// xorfwd performs XOR forwards in words, x[i] ^= x[i-4], i from 0 to len
func xorfwd(x []byte)

// xorbkd performs XOR backwords in words, x[i] ^= x[i-4], i from len to 0
func xorbkd(x []byte)
