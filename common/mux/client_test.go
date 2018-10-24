package mux_test

import (
	"testing"

	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/mux"
)

func TestIncrementalPickerFailure(t *testing.T) {
	picker := mux.IncrementalWorkerPicker{
		New: func() (*mux.ClientWorker, error) {
			return nil, errors.New("test")
		},
	}

	_, err := picker.PickAvailable()
	if err == nil {
		t.Error("expected error, but nil")
	}
}
