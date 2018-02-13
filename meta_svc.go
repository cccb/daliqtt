package main

import (
	"github.com/cameliot/alpaca"

	"log"
	"time"
)

// Actions
const PING = "@meta/PING"
const PONG = "@meta/PONG"

const WHOIS = "@meta/WHOIS"
const IAMA = "@meta/IAMA"

// Payloads
type PingPayload string
type WhoisPayload string

type PongPayload struct {
	Handle    string    `json:"handle"`
	Timestamp time.Time `json:"timestamp"`
}

type IamaPayload struct {
	Handle      string    `json:"handle"` // Unique Handle e.g. dali_mainhall
	Name        string    `json:"name"`   // Service Name e.g. dalimqtt
	Version     string    `json:"version"`
	Description string    `json:"description"`
	StartedAt   time.Time `json:"started_at"`
}

// Action Creators
func Pong(handle string) alpaca.Action {
	return alpaca.Action{
		Type: PONG,
		Payload: PongPayload{
			Handle:    handle,
			Timestamp: time.Now(),
		},
	}
}

// Huh. This hmm... well.
func Iama(payload IamaPayload) alpaca.Action {
	return alpaca.Action{
		Type:    IAMA,
		Payload: payload,
	}
}

// Service
type MetaSvc struct {
	iama      IamaPayload
	startedAt time.Time
}

func NewMetaSvc(
	handle string,
	name string,
	version string,
	description string,
) *MetaSvc {
	svc := &MetaSvc{
		iama: IamaPayload{
			Handle:      handle,
			Name:        name,
			Version:     version,
			Description: description,
			StartedAt:   time.Now(),
		},
	}

	return svc
}

func (self *MetaSvc) Handle(actions alpaca.Actions, dispatch alpaca.Dispatch) {
	log.Println("Processing meta actions")

	for action := range actions {
		switch action.Type {
		case PING:
			self.handlePing(action, dispatch)
			break
		case WHOIS:
			self.handleWhois(action, dispatch)
			break
		}
	}
}

/*
 Handle PING,
 Reply only if the requested service is a wildcard ("*") or
 identified by the service handler
*/
func (self *MetaSvc) handlePing(
	action alpaca.Action,
	dispatch alpaca.Dispatch,
) {
	payload := ""
	action.DecodePayload(&payload)

	// Are we pinged?
	if payload != "*" && payload != self.iama.Handle {
		return
	}

	// Reply with PONG
	dispatch(Pong(self.iama.Handle))
}

/*
 Handle WHOIS,

 Reply only if the requested service is a wildcard ("*") or
 identified by the service handler

 Provide information about this service
*/
func (self *MetaSvc) handleWhois(
	action alpaca.Action,
	dispatch alpaca.Dispatch,
) {
	payload := ""
	action.DecodePayload(&payload)

	// Are we pinged?
	if payload != "*" && payload != self.iama.Handle {
		return
	}

	// Reply with IAMA
	dispatch(Iama(self.iama))
}
