package log

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

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

	accessLoggerInstance.(*fileAccessLogger).close()
	accessLoggerInstance = &noOpAccessLogger{}

	content, err := ioutil.ReadFile(filename)
	assert.Error(err).IsNil()

	assert.Bool(strings.Contains(string(content), "test_from")).IsTrue()
	assert.Bool(strings.Contains(string(content), "test_to")).IsTrue()
	assert.Bool(strings.Contains(string(content), "test_reason")).IsTrue()
	assert.Bool(strings.Contains(string(content), "accepted")).IsTrue()
}
