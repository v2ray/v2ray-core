package net

import (
  "net"
)

type VMessHandler struct {
  
}

func (*VMessHandler) Listen(port uint8) error {
  listener, err := net.Listen("tcp", ":" + string(port))
  if err != nil {
    return err
  }
  
  
  return nil
}
