package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/mohamadafzal06/cache-in-go/cache"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	//LeaderAddr string
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
		go s.handleCommand(conn, msg)
		fmt.Println(string(msg))
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	msg, err := parseMessage(rawCmd)
	if err != nil {
		fmt.Println("failed to parse command")
		return
	}

	fmt.Printf("recieve command from %s\n", msg.Cmd)

	switch msg.Cmd {
	case CMDGet:
		err = s.handleGetCmd(conn, msg)
		if err != nil {
			log.Println("failed to parse command", err)
			conn.Write([]byte(err.Error()))
			return
		}
	case CMDSet:
		if bytes.Equal(msg.Value, []byte{}) {
			log.Println("there is no value to set")
			conn.Write([]byte(fmt.Sprintln("there is no value to set")))
			return
		}

		err = s.handleSetCmd(conn, msg)
		if err != nil {
			log.Println("failed to parse command", err)
			conn.Write([]byte(err.Error()))
			return
		}
	}

	//if err != nil {
	//	log.Println("failed to parse command", err)
	//	conn.Write([]byte(err.Error()))
	//}
}

func (s *Server) handleSetCmd(conn net.Conn, msg *Message) error {
	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
		return fmt.Errorf("cannot set this key-value: %w", err)
	}

	go s.sendToFollowers(context.TODO(), msg)

	return nil
}

func (s *Server) handleGetCmd(conn net.Conn, msg *Message) error {
	value, err := s.cache.Get(msg.Key)
	if err != nil {
		return fmt.Errorf("cannot set this key[%s]-value[%s]: %w",
			string(msg.Key),
			string(msg.Value), err)
	}

	_, err = conn.Write(value)
	if err != nil {
		return fmt.Errorf("cannot get value of this key[%s]: %w", string(msg.Key), err)
	}
	return nil
}

func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
	return nil
}
