// +build generate

package main

import (
	"fmt"
	"log"
	"os"
)

func writeQuarterRound(file *os.File, a, b, c, d int) {
	add := "x%d+=x%d\n"
	xor := "x=x%d^x%d\n"
	rotate := "x%d=(x << %d) | (x >> (32 - %d))\n"

	fmt.Fprintf(file, add, a, b)
	fmt.Fprintf(file, xor, d, a)
	fmt.Fprintf(file, rotate, d, 16, 16)

	fmt.Fprintf(file, add, c, d)
	fmt.Fprintf(file, xor, b, c)
	fmt.Fprintf(file, rotate, b, 12, 12)

	fmt.Fprintf(file, add, a, b)
	fmt.Fprintf(file, xor, d, a)
	fmt.Fprintf(file, rotate, d, 8, 8)

	fmt.Fprintf(file, add, c, d)
	fmt.Fprintf(file, xor, b, c)
	fmt.Fprintf(file, rotate, b, 7, 7)
}

func writeChacha20Block(file *os.File) {
	fmt.Fprintln(file, `
func ChaCha20Block(s *[16]uint32, out []byte, rounds int) {
  var x0,x1,x2,x3,x4,x5,x6,x7,x8,x9,x10,x11,x12,x13,x14,x15 = s[0],s[1],s[2],s[3],s[4],s[5],s[6],s[7],s[8],s[9],s[10],s[11],s[12],s[13],s[14],s[15]
	for i := 0; i < rounds; i+=2 {
    var x uint32
    `)

	writeQuarterRound(file, 0, 4, 8, 12)
	writeQuarterRound(file, 1, 5, 9, 13)
	writeQuarterRound(file, 2, 6, 10, 14)
	writeQuarterRound(file, 3, 7, 11, 15)
	writeQuarterRound(file, 0, 5, 10, 15)
	writeQuarterRound(file, 1, 6, 11, 12)
	writeQuarterRound(file, 2, 7, 8, 13)
	writeQuarterRound(file, 3, 4, 9, 14)
	fmt.Fprintln(file, "}")
	for i := 0; i < 16; i++ {
		fmt.Fprintf(file, "binary.LittleEndian.PutUint32(out[%d:%d], s[%d]+x%d)\n", i*4, i*4+4, i, i)
	}
	fmt.Fprintln(file, "}")
	fmt.Fprintln(file)
}

func main() {
	file, err := os.OpenFile("chacha_core.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to generate chacha_core.go: %v", err)
	}
	defer file.Close()

	fmt.Fprintln(file, "// GENERATED CODE. DO NOT MODIFY!")
	fmt.Fprintln(file, "package internal")
	fmt.Fprintln(file)
	fmt.Fprintln(file, "import \"encoding/binary\"")
	fmt.Fprintln(file)
	writeChacha20Block(file)
}
