package internal

import (
	"crypto/sha256"
	"fmt"
	"time"
)

func GenerateLinkHash(now func() time.Time, str string) string {
	strWithDate := str + now().String()
	hash := sha256.Sum256([]byte(strWithDate))
	result := fmt.Sprintf("%x", hash)
	return result[:12]
}
