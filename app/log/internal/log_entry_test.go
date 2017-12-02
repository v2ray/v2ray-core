package internal_test

import (
	"testing"

	. "v2ray.com/core/app/log/internal"
	. "v2ray.com/ext/assert"
)

func TestAccessLog(t *testing.T) {
	assert := With(t)

	entry := &AccessLog{
		From:   "test_from",
		To:     "test_to",
		Status: "Accepted",
		Reason: "test_reason",
	}

	entryStr := entry.String()
	assert(entryStr, HasSubstring, "test_from")
	assert(entryStr, HasSubstring, "test_to")
	assert(entryStr, HasSubstring, "test_reason")
	assert(entryStr, HasSubstring, "Accepted")
}
