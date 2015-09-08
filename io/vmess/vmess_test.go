package vmess

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/v2ray/v2ray-core"
)

func TestVMessSerialization(t *testing.T) {
	userId, err := core.UUIDToVID("2b2966ac-16aa-4fbf-8d81-c5f172a3da51")
	if err != nil {
		t.Fatal(err)
	}

	userSet := core.NewVUserSet()
	userSet.AddUser(core.VUser{userId})

	request := new(VMessRequest)
	request.SetVersion(byte(0x01))
	userHash := userId.Hash([]byte("ASK"))
	copy(request.UserHash(), userHash)

	_, err = rand.Read(request.RequestIV())
	if err != nil {
		t.Fatal(err)
	}

	_, err = rand.Read(request.RequestKey())
	if err != nil {
		t.Fatal(err)
	}

	_, err = rand.Read(request.ResponseHeader())
	if err != nil {
		t.Fatal(err)
	}

	request.SetCommand(byte(0x01))
	request.SetPort(80)
	request.SetDomain("v2ray.com")

	buffer := bytes.NewBuffer(make([]byte, 0, 300))
	requestWriter := NewVMessRequestWriter(userSet)
	err = requestWriter.Write(buffer, request)
	if err != nil {
		t.Fatal(err)
	}

	requestReader := NewVMessRequestReader(userSet)
	actualRequest, err := requestReader.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}

	if actualRequest.Version() != byte(0x01) {
		t.Errorf("Expected Version 1, but got %d", actualRequest.Version())
	}

	if !bytes.Equal(request.UserHash(), actualRequest.UserHash()) {
		t.Errorf("Expected user hash %v, but got %v", request.UserHash(), actualRequest.UserHash())
	}

	if !bytes.Equal(request.RequestIV(), actualRequest.RequestIV()) {
		t.Errorf("Expected request IV %v, but got %v", request.RequestIV(), actualRequest.RequestIV())
	}

	if !bytes.Equal(request.RequestKey(), actualRequest.RequestKey()) {
		t.Errorf("Expected request Key %v, but got %v", request.RequestKey(), actualRequest.RequestKey())
	}

	if !bytes.Equal(request.ResponseHeader(), actualRequest.ResponseHeader()) {
		t.Errorf("Expected response header %v, but got %v", request.ResponseHeader(), actualRequest.ResponseHeader())
	}

	if actualRequest.Command() != byte(0x01) {
		t.Errorf("Expected command 1, but got %d", actualRequest.Command())
	}

	if actualRequest.Port() != 80 {
		t.Errorf("Expected port 80, but got %d", actualRequest.Port())
	}

	if actualRequest.TargetAddress() != "v2ray.com" {
		t.Errorf("Expected target address v2ray.com, but got %s", actualRequest.TargetAddress())
	}
}
