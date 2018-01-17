package main

import (
	"log"
)

type LightsSvc struct {
	Lights []Light
	Cgi    *LichtCgi
}

func NewLightsSvc(lichtCgiBase string) *LightsSvc {
	log.Println("Using licht.cgi @", lichtCgiBase)
	// Load initial state from server
	cgi := NewLichtCgi(lichtCgiBase)

	lights, err := cgi.FetchLights()
	if err != nil {
		log.Println("Could not fetch lights from server:", err)
	}

	// Initial State
	svc := &LightsSvc{
		Lights: lights,
		Cgi:    cgi,
	}

	return svc
}

/*
 Service main: React to incoming actions and dispatch
 responses.
*/
func (self *LightsSvc) Handle(actions chan Action, dispatch Dispatch) {
	// Hanlde actions
	for action := range actions {
		switch action.Type {
		case SET_LIGHT_VALUE_REQUEST:
			dispatch(SetLightValueSuccess(23, 42))
		case GET_LIGHT_VALUES_REQUEST:
			dispatch(GetLightValuesSuccess(self.Lights))
		}
	}
}
