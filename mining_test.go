package miner

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkMining(t *testing.B) {
	c := NewMiner(1024, time.Second, func(s []byte) []byte {
		hash := sha256.Sum256(s)
		return hash[:]
	})
	t.Log("Rate:", c.GetRate())
	t.Log("Bits:", c.GetBits(time.Second))
	t.ResetTimer()
	for j := 0; j < 10; j++ {
		t.StartTimer()
		hash, nonce := c.GetHashAndNonce(time.Second, []byte("hello"))
		t.StopTimer()
		t.Log(time.Now(), hex.EncodeToString(hash))
		checkHash := sha256.Sum256(append([]byte("hello"), nonce...))
		assert.Equal(t, hash, checkHash[:], checkHash)

	}

}
