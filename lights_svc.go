package main

import (
	"log"

	"github.com/eclipse/paho.mqtt.golang"
)

type LightsSvc struct {
	Lights []Light

	mqttClient mqtt.Client
	actions    chan Action
}
