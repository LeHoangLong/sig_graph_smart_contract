package service

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
)

type hashGenerator struct {
}

func NewHashGenerator() *hashGenerator {
	return &hashGenerator{}
}

func (g *hashGenerator) New(id string, secret string) string {
	secretId := fmt.Sprintf("%s%s", id, secret)
	hashBytes := sha512.Sum512([]byte(secretId))
	currentHash := base64.StdEncoding.EncodeToString(hashBytes[:])
	return currentHash
}
