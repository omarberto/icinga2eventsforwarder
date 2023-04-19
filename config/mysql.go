package config

//import "time"

type MYSQL struct {
	Connection string        `toml:"connection"`
	Username string        `toml:"username"`
	Password string        `toml:"password"`
	RefreshInterval string        `toml:"refreshinterval"`
}
