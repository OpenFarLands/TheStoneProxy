//go:build !api

package api

import (
	"log"
	"net"

	conf "github.com/OpenFarLands/TheStoneProxy/src/config"
	"github.com/OpenFarLands/TheStoneProxy/src/server"
	"github.com/OpenFarLands/TheStoneProxy/src/utils/syncmap"
)

type ApiServer struct {
	Clients *syncmap.Map[net.Conn, *server.Client]
	Addr    string
}

func Setup(paramConfig *conf.Config, users *syncmap.Map[net.Conn, *server.Client]) error {
	log.Println("This is the build without api module! All configurations related to the api will be ignored.")
	return nil 
}
