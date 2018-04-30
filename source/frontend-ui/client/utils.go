package client

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
)

const nonceSourceSize = 1024

func RandomSha1() (string, error) {
	nonceSource := make([]byte, nonceSourceSize)
	_, err := rand.Read(nonceSource)
	if err != nil {
		return "", err
	}
	sha := sha1.Sum(nonceSource)
	return hex.EncodeToString(sha[:]), nil
}
