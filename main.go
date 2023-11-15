package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/OpenFarLands/TheStoneProxy/src/api"
	"github.com/OpenFarLands/TheStoneProxy/src/config"
	"github.com/OpenFarLands/TheStoneProxy/src/server"
)

var wg sync.WaitGroup

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
	

	wg.Add(1)
	c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
		log.Print("Stopping the server...")
		serv.StopHandle()
		wg.Done()
        os.Exit(1)
    }()
	wg.Wait()
}
