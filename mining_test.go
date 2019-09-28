package miner

import (
	"crypto/sha256"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	c := NewMiner(1024, time.Second, func(s []byte) []byte {
		hash := sha256.Sum256(s)
		return hash[:]
	})
	tt := time.Now()
	hash, nonce := c.GetHashAndNonce(10*time.Second, []byte("hello"))
	assert.Equal(t, time.Now().Unix()-tt.Unix() > 1, true, time.Now().Unix(), tt.Unix())
	assert.Equal(t, time.Now().Unix()-tt.Unix() < 30, true, time.Now().Unix(), tt.Unix())
	checkHash := sha256.Sum256(append([]byte("hello"), nonce...))
	assert.Equal(t, hash, checkHash[:], checkHash)

}
