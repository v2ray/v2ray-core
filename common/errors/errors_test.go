package errors_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/common/errors"
	"github.com/v2ray/v2ray-core/testing/unit"
)

type MockError struct {
	errors.ErrorCode
}

func (err MockError) Error() string {
	return "This is a fake error."
}

func TestHasCode(t *testing.T) {
	assert := unit.Assert(t)

	err := MockError{ErrorCode: 101}
	assert.Error(err).HasCode(101)
}
