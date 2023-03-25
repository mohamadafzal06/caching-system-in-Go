package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Command string

const (
	CMDSet Command = "SET"
	CMDGet Command = "GET"
)

//type MSGSet struct {
//	Key   []byte
//	Value []byte
//	TTL   time.Duration
//}
//
//type MSGGet struct {
//	Key []byte
//}

type Message struct {
	Cmd   Command
	Key   []byte
	Value []byte
	TTL   time.Duration
}

func parseMessage(raw []byte) (*Message, error) {
	var (
		rawStr = string(raw)
		parts  = strings.Split(rawStr, " ")
	)

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid command")
	}
	msg := &Message{
		Cmd: Command(parts[0]),
		Key: []byte(parts[1]),
	}

	if msg.Cmd == CMDSet {
		if len(parts) < 4 {
			return nil, fmt.Errorf("invalid SET command")
		}
		msg.Value = []byte(parts[2])
		n, err := strconv.Atoi(parts[3])
		if err != nil {
			log.Println("invalid SET TTl")
		}
		msg.TTL = time.Duration(n)
	}

	if msg.Cmd == CMDGet {
		return msg, nil
	}

	return msg, nil
}

func (m *Message) ToBytes() []byte {
	switch m.Cmd {
	case CMDGet:
		ret := fmt.Sprintf("%s %s", m.Cmd, string(m.Key))
		return []byte(ret)
	case CMDSet:
		ret := fmt.Sprintf("%s %s %s %d", m.Cmd, string(m.Key), string(m.Value), m.TTL)
		return []byte(ret)
	}

	return []byte{}
}
