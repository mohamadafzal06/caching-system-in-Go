package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/mohamadafzal06/cache-in-go/cache"
	"github.com/mohamadafzal06/cache-in-go/client"
)

func main() {
	var (
		listenAddr string
		leaderAddr string
	)

	flag.StringVar(&listenAddr, "listenaddr", ":3000", "listen address of the server")
	flag.StringVar(&leaderAddr, "leaderaddr", "", "listen address of the leader")
	flag.Parse()

	opts := ServerOpts{
		ListenAddr: listenAddr,
		IsLeader:   len(leaderAddr) == 0,
		LeaderAddr: leaderAddr,
	}

	go func() {
		time.Sleep(2 * time.Second)
		client, err := client.New(":3000", client.Options{})
		if err != nil {
			log.Fatalf("cannot connect to the server: %v", err)
		}
		for i := 0; i < 10; i++ {
			sendCommand(client)
			time.Sleep(200 * time.Millisecond)
		}

		client.Close()
	}()

	server := NewServer(opts, cache.New())

	server.Start()
}

func sendCommand(c *client.Client) {
	err := c.Set(context.Background(), []byte("fii"), []byte("mii"), 50000)
	if err != nil {
		log.Fatalf("cannot connect to the server: %v", err)
	}

}
