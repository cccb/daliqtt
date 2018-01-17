package main

// MQTT Messages
//
// -- Action Types
const SET_LIGHT_VALUE_REQUEST = "SET_LIGHT_VALUE_REQUEST"
const SET_LIGHT_VALUE_SUCCESS = "SET_LIGHT_VALUE_SUCCESS"
const SET_LIGHT_VALUE_ERROR = "SET_LIGHT_VALUE_ERROR"

const GET_LIGHT_VALUES_REQUEST = "GET_LIGHT_VALUES_REQUEST"
const GET_LIGHT_VALUES_SUCCESS = "GET_LIGHT_VALUES_SUCCESS"
const GET_LIGHT_VALUES_ERROR = "GET_LIGHT_VALUES_ERROR"

// -- Events
const LIGHT_VALUE_UPDATED = "LIGHT_VALUE_UPDATE"

// Actions
type Action struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Payloads
type LightValuePayload struct {
	Id    int `json:"id"`
	Value int `json:"value"`
}

type LightValuesPayload []LightValuePayload

type ErrorPayload string

// Action creators
func SetLightValueRequest(id, value, int) Action {
	return Action{
		Type: SET_LIGHT_VALUE_REQUEST,
		Payload: LightValuePayload{
			Id:    id,
			Value: value,
		},
	}
}

func SetLightValueSuccess(id, value int) Action {
	return Action{
		Type: SET_LIGHT_VALUE_SUCCESS,
		Payload: LightValuePayload{
			Id:    id,
			Value: value,
		},
	}
}

func SetLightValueError(err error) Action {
	return Action{
		Type:    SET_LIGHT_VALUE_ERROR,
		Payload: err.Error(),
	}
}

func GetLightValuesRequest() Action {
	return Action{
		Type:    GET_LIGHT_VALUES_REQUEST,
		Payload: nil,
	}
}

func GetLightValuesSuccess(lights []Light) Action {
	payload := LightValuesPayload{}
	for _, light := range lights {
		payload = append(payload, LightValuePayload{
			Id:    light.Id,
			Value: light.Value,
		})
	}

	return Action{
		Type:    GET_LIGHT_VALUES_SUCCESS,
		Payload: payload,
	}
}

func GetLightValuesError(err error) Action {
	return Action{
		Type:    GET_LIGHT_VALUES_ERROR,
		Payload: err.Error(),
	}
}
