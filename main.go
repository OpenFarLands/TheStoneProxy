package main

import (
	"log"

	"github.com/OpenFarLands/TheStoneProxy/src/api"
	"github.com/OpenFarLands/TheStoneProxy/src/config"
	"github.com/OpenFarLands/TheStoneProxy/src/server"
)

func main() {

	conf, err := config.New("./config.toml")
	if err != nil {
		log.Panic(err)
	}

	serv, err := server.New(conf.LocalAddress, conf.RemoteAddress, conf)
	if err != nil {
		log.Panic(err)
	}

	go serv.StartHandle()

	err = api.Setup(conf, &serv.Users)
	if err != nil {
		log.Panicf("Failed to setup api server: %v", err)
	}
}
