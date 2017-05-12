package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func RandStr(size int) string {
	key := make([]byte, size)
	rand.Read(key)
	return hex.EncodeToString(key)
}

func RandNumber(size int) string {
	key := make([]byte, size)
	rand.Read(key)
	for i := range key {
		key[i] = key[i]%10 + '0'
	}
	return string(key)
}
