package utils

import (
	"crypto/md5"
	"encoding/hex"
)

//CalcMd5 ..
func CalcMd5(text string) string {
	bytes := md5.Sum([]byte(text))
	return hex.EncodeToString(bytes[:])
}
