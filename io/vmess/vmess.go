// Package vmess contains protocol definition, io lib for VMess.
package vmess

import (
  "fmt"
	"io"
)

// VMessInput implements the input message of VMess protocol. It only contains
// the header of a input message. The data part will be handled by conection
// handler directly, in favor of data streaming.
type VMessInput struct {
	version  byte
	userHash [16]byte
	randHash [256]byte
	respKey  [32]byte
	iv       [16]byte
	command  byte
	port     uint16
	target   [256]byte
}

func Read(reader io.Reader) (input *VMessInput, err error) {
  buffer := make([]byte, 17 /* version + user hash */)
  nBytes, err := reader.Read(buffer)
  if err != nil {
    return
  }
  if nBytes != len(buffer) {
    err = fmt.Errorf("Unexpected length of header %d", nBytes)
    return
  }
  
  return
}

type VMessOutput [4]byte

