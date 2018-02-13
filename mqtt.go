package main

type MqttConfig struct {
	Host     string
	User     string
	Password string

	BaseTopic string
}

func (config MqttConfig) BrokerUri() string {
	uri := "tcp://"
	if config.User != "" {
		uri += config.User

		if config.Password != "" {
			uri += ":" + config.Password
		}

		uri += "@"
	}

	uri += config.Host

	return uri
}
