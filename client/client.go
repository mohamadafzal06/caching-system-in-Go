package client

import (
	"context"
	"fmt"
	"net"

	"github.com/mohamadafzal06/cache-in-go/mproto"
)

type Options struct{}

type Client struct {
	conn net.Conn
}

func New(endpoint string, opts Options) (*Client, error) {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to server: %w", err)
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("cannot close the connection: %w", err)
	}

	return nil
}

func (c *Client) Set(context context.Context, key []byte, value []byte, ttl int) error {
	cmd := &mproto.CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}
	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return fmt.Errorf("cannot send set command: %w", err)
	}
	return nil
}
