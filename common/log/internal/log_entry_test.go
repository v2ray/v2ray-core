package internal_test

import (
	"testing"

	. "github.com/v2ray/v2ray-core/common/log/internal"
	"github.com/v2ray/v2ray-core/common/serial"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestAccessLog(t *testing.T) {
	v2testing.Current(t)

	entry := &AccessLog{
		From:   serial.StringT("test_from"),
		To:     serial.StringT("test_to"),
		Status: "Accepted",
		Reason: serial.StringT("test_reason"),
	}

	entryStr := entry.String()
	assert.StringLiteral(entryStr).Contains(serial.StringT("test_from"))
	assert.StringLiteral(entryStr).Contains(serial.StringT("test_to"))
	assert.StringLiteral(entryStr).Contains(serial.StringT("test_reason"))
	assert.StringLiteral(entryStr).Contains(serial.StringT("Accepted"))
}
