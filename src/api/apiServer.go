package api

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"slices"
	"strings"

	conf "github.com/OpenFarLands/TheStoneProxy/src/config"
	"github.com/OpenFarLands/TheStoneProxy/src/server"
	"github.com/OpenFarLands/TheStoneProxy/src/utils/syncmap"
)

type ApiServer struct {
	Users *syncmap.Map[net.Conn, *server.Client]
	Addr  string
}

type apiResponse struct {
	Body  any    `json:"body"`
	Error string `json:"error"`
}

var config *conf.Config

func Setup(paramConfig *conf.Config, users *syncmap.Map[net.Conn, *server.Client]) error {
	config = paramConfig

	log.Printf("Starting api server on %v.", config.ApiServerAddress)
	
	if len(config.ApiWhitelist) == 0 {
		log.Print("Api whitelist is empty, so anyone could use it. Disable api server by setting Use_api_server to false in config.toml if you don't need it.")
	}

	serv := &ApiServer{
		Users: users,
		Addr:  config.ApiServerAddress,
	}

	http.HandleFunc("/port2ip", serv.port2ip)
	http.HandleFunc("/online", serv.online)

	go func ()  {
		err := http.ListenAndServe(config.ApiServerAddress, nil)
		if err != nil {
			log.Panic(err)
		}
	}()
	
	return nil
}

func (s *ApiServer) online(w http.ResponseWriter, r *http.Request) {
	if !isAllowed(addrStringToArray(r.RemoteAddr)[0]){
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var resp apiResponse
	w.Header().Set("Content-Type", "application/json")

	online := 0
	s.Users.Range(func(key net.Conn, value *server.Client) bool {
		online++
		return true
	})

	resp = apiResponse{
		Body: struct {
			Online int `json:"online"`
		}{
			Online: online,
		},
		Error: "",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (s *ApiServer) port2ip(w http.ResponseWriter, r *http.Request) {
	if !isAllowed(addrStringToArray(r.RemoteAddr)[0]) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var resp apiResponse
	w.Header().Set("Content-Type", "application/json")

	port := r.URL.Query().Get("port")
	if port == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{
			Body:  "",
			Error: "Port is empty",
		})
		return
	}

	s.Users.Range(func(key net.Conn, value *server.Client) bool {
		proxyPort := addrStringToArray(value.Addr.LocalAddr().String())[1]
		originIp := addrStringToArray(value.Addr.RemoteAddr().String())[0]

		if proxyPort == port {
			resp = apiResponse{
				Body: struct {
					Ip   string `json:"ip"`
					Port string `json:"port"`
				}{
					Ip:   originIp,
					Port: port,
				},
				Error: "",
			}
			return false
		}
		return true
	})

	if resp == (apiResponse{}) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiResponse{
			Body:  nil,
			Error: "This port isn't online",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func isAllowed(address string) bool {
	return slices.Contains(config.ApiWhitelist, address) || len(config.ApiWhitelist) == 0
}

func addrStringToArray(str string) []string {
	sepIndex := strings.LastIndex(str, ":")
	port := str[sepIndex+1:]
	ip := str[0:sepIndex]

	return []string{ip, port}
}