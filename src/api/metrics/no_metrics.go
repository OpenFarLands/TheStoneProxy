//go:build !metrics

package metrics

import (
	"log"

	conf "github.com/OpenFarLands/TheStoneProxy/src/config"
)

func Setup(paramConfig *conf.Config) error {
	log.Println("This is the build without metrics module! All configurations related to the prometheus metrics will be ignored.")
	return nil
}