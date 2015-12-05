package log

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/v2ray/v2ray-core/common/serial"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestAccessLog(t *testing.T) {
	v2testing.Current(t)

	filename := "/tmp/test_access_log.log"
	InitAccessLogger(filename)
	_, err := os.Stat(filename)
	assert.Error(err).IsNil()

	Access("test_from", "test_to", AccessAccepted, "test_reason")
	<-time.After(2 * time.Second)

	accessLoggerInstance.(*fileLogWriter).close()
	accessLoggerInstance = &noOpLogWriter{}

	content, err := ioutil.ReadFile(filename)
	assert.Error(err).IsNil()

	contentStr := serial.StringLiteral(content)
	assert.String(contentStr).Contains(serial.StringLiteral("test_from"))
	assert.String(contentStr).Contains(serial.StringLiteral("test_to"))
	assert.String(contentStr).Contains(serial.StringLiteral("test_reason"))
	assert.String(contentStr).Contains(serial.StringLiteral("accepted"))
}
