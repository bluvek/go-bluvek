package bvutils

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5Encode md5处理
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
