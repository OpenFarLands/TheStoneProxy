package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/OpenFarLands/TheStoneProxy/src/api"
	"github.com/OpenFarLands/TheStoneProxy/src/api/metrics"
	"github.com/OpenFarLands/TheStoneProxy/src/config"
	"github.com/OpenFarLands/TheStoneProxy/src/server"
)

var wg sync.WaitGroup

func main() {
	conf, err := config.New("./config.toml")
	if err != nil {
		log.Panic(err)
	}

	serv, err := server.New(conf.Network.LocalAddress, conf.Network.RemoteAddress, conf)
	if err != nil {
		log.Panic(err)
	}
	go serv.StartHandle()

	go func() {
		if conf.Api.UseApiServer {
			err = api.Setup(conf, &serv.Clients)
			if err != nil {
				log.Panicf("Api server error: %v", err)
			}
		}
	
		if conf.Metrics.UsePrometheus {
			err = metrics.Setup(conf)
			if err != nil {
				log.Panicf("Prometheus server error: %v", err)
			}
		}
	}()

	wg.Add(1)
	c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    func() {
        <-c
		log.Print("Stopping the server...")
		serv.StopHandle()
		wg.Done()
        os.Exit(1)
    }()
	wg.Wait()
}