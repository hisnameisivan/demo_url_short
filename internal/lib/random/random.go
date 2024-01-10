package random

import (
	"math/rand"
	"time"
)

func NewRandomString(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRTSUVWXYZ" +
		"abcdefghijklmnopqrtsuvwxyz" +
		"0123456789",
	)

	bResult := make([]rune, size)
	for i := range bResult {
		bResult[i] = chars[rnd.Intn(len(chars))]
	}

	return string(bResult)
}
