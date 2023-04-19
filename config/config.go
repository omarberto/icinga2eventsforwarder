package config

import (
    "io/ioutil"
    "os"

    "github.com/influxdata/toml"
)


type Config struct {
	Agent *AGENT

	EventsStream *EVENTSSTREAM
	Nats     *NATS
	Director *MYSQL
	Icinga *MYSQL
}

type AGENT struct {
	Host string
}

func LoadConfiguration(configFilePath string) (*Config) {
	var c Config
	
	f, err := os.Open(configFilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	if err := toml.Unmarshal(buf, &c); err != nil {
		panic(err)
	}

	return &c
}
