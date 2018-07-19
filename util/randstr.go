package util

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// NewID ...
func NewID() (id string) {
	rnum := rand.New(src).Intn(10000)
	id = fmt.Sprintf("%s%04d", strings.Replace(time.Now().Truncate(time.Millisecond).Format("20060102150405.00"), ".", "", -1), rnum)
	return
}

// RandString : make Random String
func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// RandNum ...
func RandNum(max int) int {
	return rand.New(src).Intn(max)
}

// IntRayleighCDF ...
func IntRayleighCDF() int {
	rnum := rand.New(src).Float64()
	return int(math.Sqrt(-2 * math.Log(float64(1)-rnum)))
}
