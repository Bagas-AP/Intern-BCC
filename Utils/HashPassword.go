package Utils

import (
	"crypto/sha512"
	"encoding/hex"
)

func Hash(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))
	pw := hex.EncodeToString(hash.Sum(nil))
	return pw
}
