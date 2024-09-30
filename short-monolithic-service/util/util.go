package util

import (
	"crypto/md5"
	"encoding/hex"
)

func GenerateShortCode(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[0:4])
}
