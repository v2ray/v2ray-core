package socks

import (
	// "bytes"
	"testing"

	// "github.com/v2ray/v2ray-core"
	// "github.com/v2ray/v2ray-core/testing/mocks"
	// "github.com/v2ray/v2ray-core/testing/unit"
)

func TestSocksTcpConnect(t *testing.T) {
	t.Skip("Not ready yet.")
  /*
	assert := unit.Assert(t)

	port := uint16(12384)

	uuid := "2418d087-648d-4990-86e8-19dca1d006d3"
	vid, err := core.UUIDToVID(uuid)
	assert.Error(err).IsNil()

	config := core.VConfig{
		port,
		[]core.VUser{core.VUser{vid}},
		"",
		[]core.VNext{}}

	och := new(mocks.FakeOutboundConnectionHandler)
	och.Data2Send = bytes.NewBuffer(make([]byte, 1024))
	och.Data2Return = []byte("The data to be returned to socks server.")

	vpoint, err := core.NewVPoint(&config, SocksServerFactory{}, och)
	assert.Error(err).IsNil()

	err = vpoint.Start()
	assert.Error(err).IsNil()
  */
}
