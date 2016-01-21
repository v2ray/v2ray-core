package dialer_test

import (
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
	. "github.com/v2ray/v2ray-core/transport/dialer"
)

func TestDialDomain(t *testing.T) {
	v2testing.Current(t)

	conn, err := Dial(v2net.TCPDestination(v2net.DomainAddress("google.com"), 443))
	assert.Error(err).IsNil()
	assert.StringLiteral(conn.RemoteAddr().Network()).Equals("tcp")
	conn.Close()
}
