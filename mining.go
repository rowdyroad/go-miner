package mining

import (
	"math"
	"math/rand"
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
	x.rate = x.getRate()
	return &x
}

//GetHashAndNonce get result hash and nonce
func (c *Miner) GetHashAndNonce(duration time.Duration, data []byte) ([]byte, []byte) {
	rate := duration.Nanoseconds() / c.rateDuration.Nanoseconds() * c.rate
	bits := int(math.Floor(math.Log2(float64(rate))))
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

func (c *Miner) getRate() int64 {
	b := make([]byte, c.dataLength)
	var count int64
	q := false
	go func() {
		for !q {
			c.rand.Read(b)
			c.hashFunc(b)
			count++
		}
	}()
	time.Sleep(c.rateDuration)
	q = true
	return count
}
