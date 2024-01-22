package common

import (
	"crypto/rand"
	"math/big"
)

var Letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		max := *big.NewInt(int64(len(Letters)))
		r, err := rand.Int(rand.Reader, &max)
		if err != nil {
			b[i] = Letters[0]
		} else {
			b[i] = Letters[r.Int64()]
		}
	}
	return string(b)
}
