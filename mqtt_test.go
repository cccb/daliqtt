package main

import (
	"testing"
)

func TestEncodeMessagePayload(t *testing.T) {
	action := SetLightValueRequest(2, 23)

	payload, err := encodeMessagePayload(action)
	if err != nil {
		t.Error(err)
	}

	if string(payload) != "{\"id\":2,\"value\":23}" {
		t.Error("Unexpected payload result:", string(payload))
	}

}
