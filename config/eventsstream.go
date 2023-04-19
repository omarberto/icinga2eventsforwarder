package config

// import (
//     "io/ioutil"
//     "os"

//     "github.com/influxdata/toml"
// )

type EVENTSSTREAM struct {
	UrlRequest string
	InsecureSkipVerify  bool

	Username string
	Password string
}
