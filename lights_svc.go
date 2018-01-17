package main

import (
	"log"
	"time"
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

	// Constantly poll server in case someone changed the
	// values using the legacy web interface.
	go self.watchServer(dispatch)

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

/*
 Watch the server and dispatch events in case something changed
*/
func (self *LightsSvc) watchServer(dispatch Dispatch) {
	for {
		nextLights, err := self.Cgi.FetchLights()
		if err != nil {
			log.Println(
				"An error occured while fetching state from server:",
				err,
			)
			time.Sleep(1 * time.Second)

			continue
		}

		// Diff with current values and dispatch
		// updated event if required.
		for i, nextLight := range nextLights {
			if i >= len(self.Lights) {
				log.Println("Interessting! A new light was installed?")
				dispatch(LightValueUpdated(nextLight))
				continue
			}

			currentLight := self.Lights[i]
			if currentLight.Id != nextLight.Id {
				log.Println("Something is wrong. Skipping.")
				continue
			}

			if currentLight.Value != nextLight.Value {
				// The value has changed! Inform our fellow services.
				dispatch(LightValueUpdated(nextLight))
			}
		}

		// Update state
		self.Lights = nextLights

		// Repeat after some timeout
		time.Sleep(1 * time.Second)
	}
}
