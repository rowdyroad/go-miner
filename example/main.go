package main

import (
	"crypto/sha256"
	"log"
	"time"

	"github.com/rowdyroad/go-miner"
)

func main() {
	c := miner.NewMiner(1024, time.Second, func(s []byte) []byte {
		hash := sha256.Sum256(s)
		return hash[:]
	})

	for {
		tt := time.Now()
		hash, nonce := c.GetHashAndNonce(10*time.Minute, []byte("hello"))
		checkHash := sha256.Sum256(append([]byte("hello"), nonce...))
		log.Println("OK", time.Now().Sub(tt), hash, checkHash)
	}
}
