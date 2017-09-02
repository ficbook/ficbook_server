package server

import (
	cr "crypto/rand"
	"fmt"
	mr "math/rand"
	u "github.com/satori/go.uuid"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[mr.Intn(len(letterRunes))]
    }
    return string(b)
}

func GenerateToken() string {
	b := make([]byte, 16)
	cr.Read(b)
	return fmt.Sprintf("%x", b)
}

func NewUUID() string {
	return fmt.Sprint(u.NewV4())
}