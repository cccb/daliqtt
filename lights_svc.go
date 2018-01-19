package main

import (
	"fmt"
	"log"
	"time"
)

type LightsSvc struct {
	Lights []Light
	Cgi    *LichtCgi

	updateBuffer map[int]Light
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

		updateBuffer: map[int]Light{},
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

	// Apply updates over HTTP at a constant rate
	go self.applyUpdatesProc(dispatch)

	// Hanlde actions
	for action := range actions {
		switch action.Type {
		case SET_LIGHT_VALUE_REQUEST:
			self.handleSetLightValue(action, dispatch)
		case GET_LIGHT_VALUES_REQUEST:
			dispatch(GetLightValuesSuccess(self.Lights))
		}
	}
}

func (self *LightsSvc) handleSetLightValue(action Action, dispatch Dispatch) {
	// Create new light update from
	payload := action.Payload.(LightValuePayload)

	// Update state
	if payload.Id >= len(self.Lights) {
		// Huh. This should not happen.
		dispatch(SetLightValueError(501, fmt.Errorf(
			"Set light id > registered lights",
		)))
		return
	}

	light := Light{payload.Id, payload.Value}
	self.Lights[payload.Id] = light

	// Queue update
	self.updateBuffer[light.Id] = light

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
				dispatch(SetLightValueSuccess(nextLight.Id, nextLight.Value))
				continue
			}

			currentLight := self.Lights[i]
			if currentLight.Id != nextLight.Id {
				log.Println("Something is wrong. Skipping.")
				continue
			}

			if currentLight.Value != nextLight.Value {
				// The value has changed! Inform our fellow services.
				dispatch(SetLightValueSuccess(nextLight.Id, nextLight.Value))
			}
		}

		// Update state
		self.Lights = nextLights

		// Repeat after some timeout
		time.Sleep(5 * time.Second)
	}
}

/*
 Apply updates with a constant rate
*/
func (self *LightsSvc) applyUpdatesProc(dispatch Dispatch) {

	for {
		for id, light := range self.updateBuffer {

			// Set light value on server
			err := self.Cgi.Update(light.Id, light.Value)
			if err != nil {
				dispatch(SetLightValueError(500, err))
				return
			}

			delete(self.updateBuffer, id)

			// OK
			dispatch(SetLightValueSuccess(light.Id, light.Value))

			time.Sleep(time.Second / 15) // Limit FPS
		}

		time.Sleep(time.Second / 15) // Limit FPS
	}
}
