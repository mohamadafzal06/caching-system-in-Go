package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/mohamadafzal06/cache-in-go/cache"
	"github.com/mohamadafzal06/cache-in-go/mproto"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOpts
	//	followers map[net.Conn]struct{}
	cache cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		//	followers:  make(map[net.Conn]struct{}),
		cache: c,
	}

}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("starting the servert failed: %w", err)
	}

	log.Printf("server starting on port [%s]\n", s.ListenAddr)

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

	fmt.Println("connection made", conn.RemoteAddr())

	for {
		cmd, err := mproto.ParseCommand(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("parse command error: ", err)

			break
		}

		fmt.Println(cmd)

		go s.handleCommand(conn, cmd)
	}

	fmt.Println("connection closed:", conn.RemoteAddr())
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *mproto.CommandSet:
		s.handleSetCommand(conn, v)
	case *mproto.CommandGet:
	}
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *mproto.CommandSet) error {
	log.Printf("SET %s to %s\n", cmd.Key, cmd.Value)
	if err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL)); err != nil {
		return fmt.Errorf("cannot handle this set command: %w", err)
	}

	return nil
}
