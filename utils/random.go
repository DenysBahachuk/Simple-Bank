package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alhabet = "abcdefghijklmnopqrstuvwxyz"

// Create a new random source
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomInt(min, max int64) int64 {
	return min + r.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alhabet)

	for i := 0; i < n; i++ {
		c := alhabet[r.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomAmount() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies)

	return currencies[rand.Intn(n)]
}
