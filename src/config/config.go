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
		OfflinePongMessage string
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
	var conf Config

	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		f, err := os.Create("config.toml")
		if err != nil {
			log.Fatalf("error creating config: %v", err)
		}
		data, err := toml.Marshal(conf)
		if err != nil {
			log.Fatalf("error encoding default config: %v", err)
		}
		if _, err := f.Write(data); err != nil {
			log.Fatalf("error writing encoded default config: %v", err)
		}
		_ = f.Close()
	}

	data, err := os.ReadFile("config.toml")
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &conf); err != nil {
		log.Fatalf("error decoding config: %v", err)
	}

	return &conf, nil
}
