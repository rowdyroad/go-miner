package miner

import (
	"math"
	"math/rand"
	"runtime"
	"time"
)

//Miner struct
type Miner struct {
	dataLength   int
	rate         int64
	rateDuration time.Duration
	hashFunc     func([]byte) []byte
	rand         *rand.Rand
}

//NewMiner constructor for Miner struct
func NewMiner(dataLength int, rateDuration time.Duration, f func([]byte) []byte) *Miner {
	x := Miner{
		dataLength:   dataLength,
		rateDuration: rateDuration,
		hashFunc:     f,
		rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return &x
}

//GetHashAndNonce get result hash and nonce
func (c *Miner) GetHashAndNonce(duration time.Duration, data []byte) ([]byte, []byte) {
	bits := c.GetBits(duration)
	nonce := make([]byte, c.dataLength-len(data))
	for {
		c.rand.Read(nonce)
		ret := c.hashFunc(append(data, nonce...))
		ok := true
		for j := 0; j < bits/8; j++ {
			if ret[j] != 0 {
				ok = false
				break
			}
		}
		if ok {
			for j := 0; j < bits%8; j++ {
				if bit := ret[bits/8+1] >> byte(j) & 0x01; bit == 1 {
					ok = false
					break
				}
			}
		}
		if ok {
			return ret, nonce
		}
	}
	return nil, nil
}

func (c *Miner) GetBits(duration time.Duration) int {
	rate := duration.Nanoseconds() / c.rateDuration.Nanoseconds() * c.getRate()
	return int(math.Floor(math.Log2(float64(rate))))
}

func (c *Miner) GetRate() int64 {
	if c.rate == 0 {
		c.rate = c.getRate()
	}
	return c.rate
}

func (c *Miner) getRate() int64 {
	runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(runtime.NumCPU())
	b := make([]byte, c.dataLength)

	var count int64
	end := time.Now().Add(c.rateDuration)
	for time.Now().Before(end) {
		c.rand.Read(b)
		c.hashFunc(b)
		count++
	}
	return count
}
