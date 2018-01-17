package main

import (
	"testing"
)

func TestFetchLights(t *testing.T) {
	t.Log("Fetching lights")
	cgi := NewLichtCgi("http://localhost:2299")
	lights, err := cgi.FetchLights()
	if err != nil {
		t.Log("licht.cgi stub server running?")
		t.Error(err)
	}

	if len(lights) != 4 {
		t.Error("Expected 4 values")
	}

	t.Log("Success. Retrieved:", lights)
}

func TestSetLight(t *testing.T) {
	cgi := NewLichtCgi("http://localhost:2299")
	err := cgi.Update(Light{2, 127})
	if err != nil {
		t.Log("licht.cgi stub server running?")
		t.Error(err)
	}

	// Read, check update
	lights, err := cgi.FetchLights()
	if err != nil {
		t.Log("licht.cgi stub server running?")
		t.Error(err)
	}

	if lights[2].Value != 127 {
		t.Error("Expected light 2 to be set to 127")
	}

	t.Log("Success. Retrieved:", lights)

	err = cgi.Update(Light{2, 23})
	if err != nil {
		t.Log("licht.cgi stub server running?")
		t.Error(err)
	}

	// Read, check update
	lights, err = cgi.FetchLights()
	if err != nil {
		t.Log("licht.cgi stub server running?")
		t.Error(err)
	}

	if lights[2].Value != 23 {
		t.Error("Expected light 2 to be set to 23")
	}

	t.Log("Success. Retrieved:", lights)
}
