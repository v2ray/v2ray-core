package main

import (
  "log"
  
  "github.com/v2ray/v2ray-core"
  "github.com/v2ray/v2ray-core/net/freedom"
  "github.com/v2ray/v2ray-core/net/socks"
)

func main() {
  port := uint16(8888)
  
  uuid := "2418d087-648d-4990-86e8-19dca1d006d3"
	vid, err := core.UUIDToVID(uuid)
  if err != nil {
    log.Fatal(err)
  }

	config := core.VConfig{
		port,
		[]core.VUser{core.VUser{vid}},
		"",
		[]core.VNext{}}

	vpoint, err := core.NewVPoint(&config, socks.SocksServerFactory{}, freedom.FreedomFactory{})
  if err != nil {
    log.Fatal(err)
  }
  err = vpoint.Start()
  if err != nil {
    log.Fatal(err)
  }
  
  finish := make(chan bool)
  <-finish
}
