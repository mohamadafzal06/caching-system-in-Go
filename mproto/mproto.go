package mproto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Command byte

const (
	CMDNonce Command = iota
	CMDSet
	CMDGet
	CmdDel
)

type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int
}

func (c *CommandSet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CMDSet)

	binary.Write(buf, binary.LittleEndian, int32(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)

	binary.Write(buf, binary.LittleEndian, int32(len(c.Value)))
	binary.Write(buf, binary.LittleEndian, c.Value)

	binary.Write(buf, binary.LittleEndian, int32(c.TTL))

	return buf.Bytes()
}

type CommandGet struct {
	Key []byte
}

func (c *CommandGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CMDGet)

	binary.Write(buf, binary.LittleEndian, int32(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)

	return buf.Bytes()
}

func ParseCommand(r io.Reader) (any, error) {
	var cmd Command
	if err := binary.Read(r, binary.LittleEndian, &cmd); err != nil {
		return nil, err
	}

	switch cmd {
	case CMDSet:
		return parseSetCommand(r), nil
	case CMDGet:
		return parseGetCommand(r), nil
	case CmdDel:
	default:
		return nil, fmt.Errorf("invalid command")
	}

	return nil, nil
}

func parseSetCommand(r io.Reader) *CommandSet {
	cmd := &CommandSet{}

	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &cmd.Key)

	var valueLen int32
	binary.Read(r, binary.LittleEndian, &valueLen)
	cmd.Value = make([]byte, valueLen)
	binary.Read(r, binary.LittleEndian, &cmd.Value)

	var ttl int32
	binary.Read(r, binary.LittleEndian, &ttl)
	cmd.TTL = int(ttl)
	return cmd
}

func parseGetCommand(r io.Reader) *CommandGet {
	cmd := &CommandGet{}

	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &cmd.Key)

	return cmd
}
