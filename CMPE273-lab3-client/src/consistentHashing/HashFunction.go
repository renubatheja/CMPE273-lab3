package consistentHashing

import (
	"crypto/md5"
)


func HashMD5(key string) []byte {
	m := md5.New()
	m.Write([]byte(key))
	return m.Sum(nil)
}