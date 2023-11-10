package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	LocalAddress     string   `toml:"Local_address"`
	RemoteAddress    string   `toml:"Remote_address"`
	Timeout          int      `toml:"Timeout"`
	MotdGetInterval  int      `toml:"Motd_get_interval"`
	UseApiServer     bool     `toml:"Use_api_server"`
	ApiServerAddress string   `toml:"Api_server_address"`
	ApiWhitelist     []string `toml:"Api_whitelist"`
}

func New(filename string) (*Config, error) {
	var conf Config

	_, err := toml.DecodeFile(filename, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
