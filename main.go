package main

import (
	"flag"
	"log"

	"github.com/cameliot/alpaca"
)

var version = "unknown"

type Config struct {
	Mqtt         *MqttConfig
	LichtCgiBase string
}

func parseFlags() *Config {
	host := flag.String("host", "localhost:1883", "MQTT broker host")
	user := flag.String("user", "", "MQTT broker host")
	password := flag.String("password", "", "MQTT broker host")
	baseTopic := flag.String("topic", "dali", "MQTT base topic")

	lichtCgiBase := flag.String("lichtcgi", "http://dali", "licht.cgi server")

	flag.Parse()

	mqttConfig := &MqttConfig{
		Host:     *host,
		User:     *user,
		Password: *password,

		BaseTopic: *baseTopic,
	}

	config := &Config{
		Mqtt:         mqttConfig,
		LichtCgiBase: *lichtCgiBase,
	}

	return config
}

func main() {
	log.Println("Starting dali to mqtt bridge v.", version)

	// Initialize configuration
	config := parseFlags()

	// Initialize MQTT connection
	actions, dispatch := alpaca.DialMqtt(
		config.Mqtt.BrokerUri(),
		alpaca.Routes{
			"lights": config.Mqtt.BaseTopic,
			"meta":   "v1/_meta",
		},
	)

	lightsActions := make(alpaca.Actions)
	metaActions := make(alpaca.Actions)

	// So far so good. Let now the lights service
	// handle actions
	lightsSvc := NewLightsSvc(config.LichtCgiBase)
	go lightsSvc.Handle(lightsActions, dispatch)

	// Hanlde meta actions for service discovery
	metaSvc := NewMetaSvc(
		"daliqtt@mainhall",
		"daliqtt",
		version,
		"Bridge Dali / Lights to MQTT",
	)
	go metaSvc.Handle(metaActions, dispatch)

	for action := range actions {
		lightsActions <- action
		metaActions <- action
	}

}
