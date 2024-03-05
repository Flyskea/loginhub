package random

import "math/rand"

const (
	Lowercase    = "abcdefghijklmnopqrstuvwxyz"
	Upper        = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Symbol       = "~!@#$%^&*()_+{}|:\"<>?[];',./"
	numbers      = "0123456789"
	Alphanumeric = Lowercase + Upper + numbers
	All          = Alphanumeric + Lowercase + Upper + Symbol
)

func RandomString(set string, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = set[rand.Intn(len(set))]
	}
	return string(b)
}
