package server

import (
	"errors"
	"log"
	"net"
	"sync"

	"context"
	"fmt"
	"strings"
	"time"

	conf "github.com/OpenFarLands/TheStoneProxy/src/config"
	"github.com/OpenFarLands/TheStoneProxy/src/utils/syncmap"
	"github.com/OpenFarLands/go-raknet"
)

const defaultTimeout = 8

var config *conf.Config

type Server struct {
	Clients      syncmap.Map[net.Conn, *Client]
	ProxyAddr    *net.UDPAddr
	UpstreamAddr string
	Timeout      int
	Listener     *raknet.Listener
}

func New(proxyAddr, upstreamAddr string, paramConfig *conf.Config) (*Server, error) {
	config = paramConfig

	timeout := config.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	proxyUdpAddr, err := net.ResolveUDPAddr("udp4", proxyAddr)
	if err != nil {
		return nil, err
	}

	return &Server{
		ProxyAddr:    proxyUdpAddr,
		UpstreamAddr: upstreamAddr,
		Timeout:      timeout,
	}, nil
}

func (s *Server) handleConnection(conn net.Conn) {
	log.Printf("Сlient connected: %v", conn.RemoteAddr().String())

	server, err := raknet.Dial(s.UpstreamAddr)
	if err != nil {
		server, err = raknet.Dial(s.UpstreamAddr)
		if err != nil {
			conn.Close()
			log.Printf("Error connecting to the server: %v\n", err)
			return
		}
	}

	raknetConn, ok := conn.(*raknet.Conn)
	if !ok {
		log.Print("Error: failed to use net.Conn as raknet.Conn")
		return
	}
	s.Clients.Store(server, &Client{Addr: raknetConn})

	addOnline(1)
	defer func() {
		addOnline(-1)
		server.Close()
		conn.Close()
		s.Clients.Delete(server)
		log.Printf("Client disconnected: %v", conn.RemoteAddr().String())
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	// Read from client, proxy to server
	go func() {
		b := make([]byte, 1024*1024*4)
		for {
			conn.SetDeadline(time.Now().Add(time.Duration(s.Timeout) * time.Second))

			n, err := conn.Read(b)
			if err != nil {
				if !raknet.ErrConnectionClosed(err) && !errors.Is(err, context.DeadlineExceeded) {
					log.Printf("Error reading from client connection: %v\n", err)
				}
				wg.Done()
				return
			}

			packet := b[:n]

			_, err = server.Write(packet)
			if err != nil {
				if !raknet.ErrConnectionClosed(err) {
					log.Printf("Error writing to server connection: %v\n", err)
				}
				wg.Done()
				return
			}
		}
	}()

	// Read from server, relay to client
	go func() {
		b := make([]byte, 1024*1024*4)
		for {
			server.SetDeadline(time.Now().Add(time.Duration(s.Timeout) * time.Second))
			n, err := server.Read(b)
			if err != nil {
				if !raknet.ErrConnectionClosed(err) && !errors.Is(err, context.DeadlineExceeded) {
					log.Printf("Error reading from server connection: %v\n", err)
				}
				wg.Done()
				return
			}

			packet := b[:n]

			_, err = conn.Write(packet)
			if err != nil {
				if !raknet.ErrConnectionClosed(err) {
					log.Printf("Error writing to client connection: %v\n", err)
				}
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
}

func (s *Server) StopHandle() {
	s.Clients.Range(func(key net.Conn, value *Client) bool {
		value.Addr.Close()
		key.Close()
		s.Clients.Delete(key)
		return true
	})
}

func (s *Server) StartHandle() {
	log.Printf("Starting listening on %v, proxying to %v.", s.ProxyAddr.String(), s.UpstreamAddr)

	listener, err := raknet.Listen(s.ProxyAddr.String())
	if err != nil {
		log.Panic(err)
	}
	defer listener.Close()

	// Get motd from upstream
	go func() {
		for {
			motd, err := raknet.PingTimeout(s.UpstreamAddr, time.Second)
			if err != nil {
				continue
			}

			arrayMotd := strings.Split(string(motd), ";")
			arrayMotd[6] = fmt.Sprint(listener.ID())
			arrayMotd[10] = fmt.Sprint(s.ProxyAddr.Port)
			arrayMotd[11] = fmt.Sprint(s.ProxyAddr.Port)
			stringMotd := strings.Join(arrayMotd, ";")

			listener.PongData([]byte(stringMotd))

			time.Sleep(time.Duration(config.MotdGetInterval) * time.Second) 
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("Failed to accept connection from client")
			continue
		}

		go s.handleConnection(conn)
	}
}
