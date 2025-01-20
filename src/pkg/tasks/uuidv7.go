package tasks

import (
	"crypto/rand"
	"math/big"
	"time"
)

func Uuid() [16]byte {
	var value [16]byte
	_, err := rand.Read(value[:])
	if err == nil {
		timestamp := big.NewInt(time.Now().UnixMilli())
		timestamp.FillBytes(value[0:6])

		value[6] = (value[6] & 0x0F) | 0x70
		value[8] = (value[8] & 0x3F) | 0x80
	}

	return value
}
