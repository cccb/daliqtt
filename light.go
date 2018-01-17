package main

import (
	"fmt"
)

// Handle light state
type Light struct {
	Id    int
	Value int
}

type LightsState map[string]Light

/*
 Initialize lights state:
 Assign handles to light value
*/
func NewLightsState() *LightsState {
	state := &LightsState{
		"entry":     Light{0, 0},
		"foh":       Light{1, 0},
		"desk_wall": Light{2, 0},
		"desk_bar":  Light{3, 0},
	}

	return state
}

func (self *LightsState) Set(handle string, value int) {
	light := (*self)[handle]
	light.Value = value
	(*self)[handle] = light

	// TODO: Talk to dali
}

func (self *LightsState) Read(handle string) int {
	// TODO: Talk to dali

	return (*self)[handle].Value
}

func (self *LightsState) Refresh() {
	for handle, _ := range *self {
		fmt.Println("Refreshing state of", handle)
		self.Read(handle)
	}
}
