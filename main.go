package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/mohamadafzal06/cache-in-go/cache"
)

func main() {
	var (
		listenAddr string
		leaderAddr string
	)

	flag.StringVar(&listenAddr, "listenaddr", ":4000", "listen address of the server")
	flag.StringVar(&leaderAddr, "leaderaddr", ":3000", "listen address of the leader")
	flag.Parse()

	opts := ServerOpts{
		//ListenAddr: ":4000",
		ListenAddr: listenAddr,
		IsLeader:   len(leaderAddr) == 0,
		LeaderAddr: leaderAddr,
	}

	//conn, err := net.Dial("tcp", ":4000")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//_, err = conn.Write([]byte("SET foo bar 5_000_000"))
	//if err != nil {
	//	log.Fatal(err)
	//}

	go func() {
		time.Sleep(2 * time.Second)
		conn, err := net.Dial("tcp", ":4000")
		if err != nil {
			log.Fatal(err)
		}

		conn.Write([]byte("SET foo bar 250000000000000000"))

		time.Sleep(2 * time.Second)

		conn.Write([]byte("GET foo"))
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		fmt.Println(string(buf[:n]))
	}()

	server := NewServer(opts, cache.New())
	err := server.Start()
	if err != nil {
		panic("server cannot start")
	}

}
