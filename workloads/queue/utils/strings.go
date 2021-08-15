package utils

import (
	"math/rand"
	"time"
	"strconv"
	"fmt"
)

const kLetterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const TimestampStrLen = 20

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = kLetterBytes[rand.Intn(len(kLetterBytes))]
	}
	return string(b)
}

func FormatTime(t time.Time) string {
	return fmt.Sprintf("%020d", t.UnixNano())
}

func ParseTime(payload string) time.Time {
	timeStr := payload[0:TimestampStrLen]
	if s, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
		return time.Unix(0, s)
	} else {
		panic(err)
	}
}
