package config

import (
    "wp/tls"
)

type NATS struct {
	Server string
	Subject string

	Secure bool
	//TlsCa       string        `toml:"tls_ca"`
	//TlsKey      string        `toml:"tls_key"`
	//InsecureSkipVerify  bool
	tls.ClientConfig
}
