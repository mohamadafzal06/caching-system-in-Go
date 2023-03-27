package mproto

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSetCommand(t *testing.T) {
	cmd := &CommandSet{
		Key:   []byte("Foo"),
		Value: []byte("Bar"),
		TTL:   2,
	}

	fmt.Println(string(cmd.Bytes()))

	r := bytes.NewReader(cmd.Bytes())

	parsedCmd, err := ParseCommand(r)
	assert.Nil(t, err)

	assert.Equal(t, cmd, parsedCmd)

}

func TestParseGetCommand(t *testing.T) {
	cmd := &CommandGet{
		Key: []byte("Foo"),
	}

	fmt.Println(string(cmd.Bytes()))

	r := bytes.NewReader(cmd.Bytes())

	parsedCmd, err := ParseCommand(r)
	assert.Nil(t, err)

	assert.Equal(t, cmd, parsedCmd)

}
