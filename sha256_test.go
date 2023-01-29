package main

import (
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsersha256(t *testing.T) {
	message := []byte("vitwit.com")
	userhashByteArray := Usersha256(PaddedData(string(message)))
	userhash := string(userhashByteArray)
	assert.NotNil(t, userhash)
}

func TestUsersha256HashSize(t *testing.T) {
	message := []byte("saiteja")
	userhashByteArray := Usersha256(PaddedData(string(message)))
	assert.Equal(t, 8*len(userhashByteArray), 256)
}

func TestUsersha256ToInbuiltSha256(t *testing.T) {
	message := []byte("vitwitcom")
	userhash := string(Usersha256(PaddedData(string(message))))
	h := sha256.New()
	h.Write(message)
	inbuilthash := string(h.Sum(nil))
	assert.Equal(t, userhash, inbuilthash)
}
