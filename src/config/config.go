package config

import (
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Network struct {
		LocalAddress       string
		RemoteAddress      string
		Timeout            int
		MotdGetInterval    int
	}
	OfflineMotd struct {
		Motd            string
		ProtocolVersion int
		VersionName     string
		LevelName       string
	}
	Api struct {
		UseApiServer     bool
		ApiServerAddress string
		ApiWhitelist     []string
	}
	Metrics struct {
		UsePrometheus             bool
		PrometheusAddress         string
		PrometheusBearerAuthToken string
	}
}

func New(filename string) (*Config, error) {
	conf := Config{
		Network: struct {
			LocalAddress       string
			RemoteAddress      string
			Timeout            int
			MotdGetInterval    int
		}{
			LocalAddress:       "0.0.0.0:19132",
			RemoteAddress:      "0.0.0.0:19133",
			Timeout:            60,
			MotdGetInterval:    10,
		},
		OfflineMotd: struct {
			Motd            string
			ProtocolVersion int
			VersionName     string
			LevelName       string
		}{
			Motd: "§c§lOffline",
			ProtocolVersion: 1,
			VersionName: "1.0.0",
			LevelName: "Powered by TheStoneProxy",
		},
	}

	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		f, err := os.Create("config.toml")
		if err != nil {
			log.Fatalf("Error creating config: %v", err)
		}
		data, err := toml.Marshal(conf)
		if err != nil {
			log.Fatalf("Error encoding default config: %v", err)
		}
		if _, err := f.Write(data); err != nil {
			log.Fatalf("Error writing encoded default config: %v", err)
		}
		f.Close()
	}

	data, err := os.ReadFile("config.toml")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &conf); err != nil {
		log.Fatalf("Error decoding config: %v", err)
	}

	if conf.Network.LocalAddress == "" {
		log.Fatal("Config error: There is a config, but Network.LocalAddress is empty. Please specify the address where proxy should listen.")
	}
	if conf.Network.RemoteAddress == "" {
		log.Fatal("Config error: There is a config, but Network.RemoteAddress is empty. Please specify the address where proxy should proxy to.")
	}
	if conf.Network.Timeout <= 0 {
		log.Fatal("Config error: There is a config, but Network.Timeout is less that or equal to zero. Please specify the time in seconds before proxy disconnects inactive clients.")
	}
	if conf.Network.MotdGetInterval <= 0 {
		log.Fatal("Config error: There is a config, but Network.MotdGetInterval is less that or equal to zero. Please specify the interval for proxy to fetch and update motd from upstream server.")
	}

	return &conf, nil
}
