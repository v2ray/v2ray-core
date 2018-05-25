package internal

import "encoding/binary"

func ChaCha20Block(s *[16]uint32, out []byte, rounds int) {
	var x0, x1, x2, x3, x4, x5, x6, x7, x8, x9, x10, x11, x12, x13, x14, x15 = s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7], s[8], s[9], s[10], s[11], s[12], s[13], s[14], s[15]
	for i := 0; i < rounds; i += 2 {
		var x uint32

		x0 += x4
		x = x12 ^ x0
		x12 = (x << 16) | (x >> (32 - 16))
		x8 += x12
		x = x4 ^ x8
		x4 = (x << 12) | (x >> (32 - 12))
		x0 += x4
		x = x12 ^ x0
		x12 = (x << 8) | (x >> (32 - 8))
		x8 += x12
		x = x4 ^ x8
		x4 = (x << 7) | (x >> (32 - 7))
		x1 += x5
		x = x13 ^ x1
		x13 = (x << 16) | (x >> (32 - 16))
		x9 += x13
		x = x5 ^ x9
		x5 = (x << 12) | (x >> (32 - 12))
		x1 += x5
		x = x13 ^ x1
		x13 = (x << 8) | (x >> (32 - 8))
		x9 += x13
		x = x5 ^ x9
		x5 = (x << 7) | (x >> (32 - 7))
		x2 += x6
		x = x14 ^ x2
		x14 = (x << 16) | (x >> (32 - 16))
		x10 += x14
		x = x6 ^ x10
		x6 = (x << 12) | (x >> (32 - 12))
		x2 += x6
		x = x14 ^ x2
		x14 = (x << 8) | (x >> (32 - 8))
		x10 += x14
		x = x6 ^ x10
		x6 = (x << 7) | (x >> (32 - 7))
		x3 += x7
		x = x15 ^ x3
		x15 = (x << 16) | (x >> (32 - 16))
		x11 += x15
		x = x7 ^ x11
		x7 = (x << 12) | (x >> (32 - 12))
		x3 += x7
		x = x15 ^ x3
		x15 = (x << 8) | (x >> (32 - 8))
		x11 += x15
		x = x7 ^ x11
		x7 = (x << 7) | (x >> (32 - 7))
		x0 += x5
		x = x15 ^ x0
		x15 = (x << 16) | (x >> (32 - 16))
		x10 += x15
		x = x5 ^ x10
		x5 = (x << 12) | (x >> (32 - 12))
		x0 += x5
		x = x15 ^ x0
		x15 = (x << 8) | (x >> (32 - 8))
		x10 += x15
		x = x5 ^ x10
		x5 = (x << 7) | (x >> (32 - 7))
		x1 += x6
		x = x12 ^ x1
		x12 = (x << 16) | (x >> (32 - 16))
		x11 += x12
		x = x6 ^ x11
		x6 = (x << 12) | (x >> (32 - 12))
		x1 += x6
		x = x12 ^ x1
		x12 = (x << 8) | (x >> (32 - 8))
		x11 += x12
		x = x6 ^ x11
		x6 = (x << 7) | (x >> (32 - 7))
		x2 += x7
		x = x13 ^ x2
		x13 = (x << 16) | (x >> (32 - 16))
		x8 += x13
		x = x7 ^ x8
		x7 = (x << 12) | (x >> (32 - 12))
		x2 += x7
		x = x13 ^ x2
		x13 = (x << 8) | (x >> (32 - 8))
		x8 += x13
		x = x7 ^ x8
		x7 = (x << 7) | (x >> (32 - 7))
		x3 += x4
		x = x14 ^ x3
		x14 = (x << 16) | (x >> (32 - 16))
		x9 += x14
		x = x4 ^ x9
		x4 = (x << 12) | (x >> (32 - 12))
		x3 += x4
		x = x14 ^ x3
		x14 = (x << 8) | (x >> (32 - 8))
		x9 += x14
		x = x4 ^ x9
		x4 = (x << 7) | (x >> (32 - 7))
	}
	binary.LittleEndian.PutUint32(out[0:4], s[0]+x0)
	binary.LittleEndian.PutUint32(out[4:8], s[1]+x1)
	binary.LittleEndian.PutUint32(out[8:12], s[2]+x2)
	binary.LittleEndian.PutUint32(out[12:16], s[3]+x3)
	binary.LittleEndian.PutUint32(out[16:20], s[4]+x4)
	binary.LittleEndian.PutUint32(out[20:24], s[5]+x5)
	binary.LittleEndian.PutUint32(out[24:28], s[6]+x6)
	binary.LittleEndian.PutUint32(out[28:32], s[7]+x7)
	binary.LittleEndian.PutUint32(out[32:36], s[8]+x8)
	binary.LittleEndian.PutUint32(out[36:40], s[9]+x9)
	binary.LittleEndian.PutUint32(out[40:44], s[10]+x10)
	binary.LittleEndian.PutUint32(out[44:48], s[11]+x11)
	binary.LittleEndian.PutUint32(out[48:52], s[12]+x12)
	binary.LittleEndian.PutUint32(out[52:56], s[13]+x13)
	binary.LittleEndian.PutUint32(out[56:60], s[14]+x14)
	binary.LittleEndian.PutUint32(out[60:64], s[15]+x15)
}
