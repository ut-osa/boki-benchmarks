package utils

import (
	"math/rand"
)

const kLetterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = kLetterBytes[rand.Intn(len(kLetterBytes))]
	}
	return string(b)
}
