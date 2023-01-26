package main

import (
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsersha256(t *testing.T) {
	message := "vitwit.com"
	userhash := string(Usersha256([]byte(message)))
	h := sha256.New()
	h.Write([]byte(message))
	inbuilthash := string(h.Sum(nil))
	assert.Equal(t, userhash, inbuilthash)
}
