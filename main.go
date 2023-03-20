package main

import (
	"log"
	"net"
	"time"

	"github.com/mohamadafzal06/cache-in-go/cache"
)

func main() {
	opts := ServerOpts{
		ListenAddr: ":4000",
		IsLeader:   true,
	}

	go func() {
		time.Sleep(2 * time.Second)
		conn, err := net.Dial("tcp", ":4000")
		if err != nil {
			log.Fatal(err)
		}

		conn.Write([]byte("set foo bar"))
	}()

	server := NewServer(opts, cache.New())
	err := server.Start()
	if err != nil {
		panic("server cannot start")
	}

}
