package server

import (
	"time"

	"github.com/OpenFarLands/go-raknet"
)

type Client struct {
	Addr *raknet.Conn
}

func (c *Client) GetLatency() time.Duration {
	return c.Addr.Latency()
}

func (c *Client) Close() error {
	return c.Addr.Close()
}

func (c *Client) SetDeadline(time time.Time) error {
	return c.Addr.SetDeadline(time)
}
