package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

func GetSign(key, secret string) string {
	timestamp := time.Now().UnixMilli()
	data := fmt.Sprintf("%d\n%s\n%s", timestamp, secret, key)

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	sign := h.Sum(nil)

	return fmt.Sprintf("%d%s", timestamp, base64.StdEncoding.EncodeToString(sign))
}
