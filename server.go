package main

import (
	"fmt"
	"log"
	"net"

	"github.com/mohamadafzal06/cache-in-go/cache"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
}

type Server struct {
	ServerOpts
	cache cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
	}

}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("starting the servert failed: %w", err)
	}

	log.Printf("server starting on port[%s]\n", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("error while accepting request: %v", err)
			continue
		}
		go s.handleConn(conn)

	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("error while reading data from connection: %v", err)
			break
		}

		msg := buf[:n]
		fmt.Println(string(msg))
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	msg, err := parseMessage(rawCmd)
	if err != nil {
		fmt.Println("failed to parse command")
		return
	}

	switch msg.Cmd {
	case CMDGet:
		s.handleGetCmd(conn, msg)
	case CMDSet:
		s.handleSetCmd(conn, msg)
	}
}

func (s *Server) handleSetCmd(conn net.Conn, msg *Message) error {
	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
		return fmt.Errorf("cannot set this key-value: %w", err)
	}

	return nil
}
func (s *Server) handleGetCmd(conn net.Conn, msg *Message) ([]byte, error) {
	value, err := s.cache.Get(msg.Key)
	if err != nil {
		return []byte{}, fmt.Errorf("cannot set this key-value: %w", err)
	}

	return value, nil
}
