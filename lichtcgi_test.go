package main

import (
	"testing"
)

func TestFetchLights(t *testing.T) {
	t.Log("Fetching lights")

	cgi := NewLichtCgi("http://dali")
	lights, err := cgi.FetchLights()
	if err != nil {
		t.Error(err)
	}

	t.Log("Res:", lights)
}
