package log_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	. "v2ray.com/core/common/log"
	. "v2ray.com/ext/assert"
)

func TestFileLogger(t *testing.T) {
	assert := With(t)

	f, err := ioutil.TempFile("", "vtest")
	assert(err, IsNil)
	path := f.Name()
	common.Must(f.Close())

	creator, err := CreateFileLogWriter(path)
	assert(err, IsNil)

	handler := NewLogger(creator)
	handler.Handle(&GeneralMessage{Content: "Test Log"})
	time.Sleep(2 * time.Second)

	common.Must(common.Close(handler))

	f, err = os.Open(path)
	assert(err, IsNil)

	b, err := buf.ReadAllToBytes(f)
	assert(err, IsNil)
	assert(string(b), HasSubstring, "Test Log")

	common.Must(f.Close())
}
