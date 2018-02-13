package main

import (
	"github.com/cameliot/alpaca"
)

// MQTT Messages
//
// -- Action Types
const SET_LIGHT_VALUE_REQUEST = "@lights/SET_LIGHT_VALUE_REQUEST"
const SET_LIGHT_VALUE_SUCCESS = "@lights/SET_LIGHT_VALUE_SUCCESS"
const SET_LIGHT_VALUE_ERROR = "@lights/SET_LIGHT_VALUE_ERROR"

const GET_LIGHT_VALUES_REQUEST = "@lights/GET_LIGHT_VALUES_REQUEST"
const GET_LIGHT_VALUES_SUCCESS = "@lights/GET_LIGHT_VALUES_SUCCESS"
const GET_LIGHT_VALUES_ERROR = "@lights/GET_LIGHT_VALUES_ERROR"

// Payloads
type LightValuePayload struct {
	Id    int `json:"id"`
	Value int `json:"value"`
}

type LightValuesPayload []LightValuePayload

type ErrorPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Action creators
func SetLightValueRequest(id, value int) alpaca.Action {
	return alpaca.Action{
		Type: SET_LIGHT_VALUE_REQUEST,
		Payload: LightValuePayload{
			Id:    id,
			Value: value,
		},
	}
}

func SetLightValueSuccess(id, value int) alpaca.Action {
	return alpaca.Action{
		Type: SET_LIGHT_VALUE_SUCCESS,
		Payload: LightValuePayload{
			Id:    id,
			Value: value,
		},
	}
}

func SetLightValueError(code int, err error) alpaca.Action {
	return alpaca.Action{
		Type: SET_LIGHT_VALUE_ERROR,
		Payload: ErrorPayload{
			Code:    code,
			Message: err.Error(),
		},
	}
}

func GetLightValuesRequest() alpaca.Action {
	return alpaca.Action{
		Type:    GET_LIGHT_VALUES_REQUEST,
		Payload: nil,
	}
}

func GetLightValuesSuccess(lights []Light) alpaca.Action {
	payload := LightValuesPayload{}
	for _, light := range lights {
		payload = append(payload, LightValuePayload{
			Id:    light.Id,
			Value: light.Value,
		})
	}

	return alpaca.Action{
		Type:    GET_LIGHT_VALUES_SUCCESS,
		Payload: payload,
	}
}

func GetLightValuesError(err error) alpaca.Action {
	return alpaca.Action{
		Type:    GET_LIGHT_VALUES_ERROR,
		Payload: err.Error(),
	}
}
