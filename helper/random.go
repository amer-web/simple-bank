package helper

import (
	"fmt"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	return stringWithCharset(length, charset)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EGP", "EUR"}
	length := len(currencies)
	return currencies[rand.Intn(length)]
}
func RandomEmail() string {
	return fmt.Sprintf("%s@example.com", RandomString(8))
}
