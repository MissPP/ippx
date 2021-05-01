package conn

import (
	"crypto/rc4"
)

var key = "ippxippxippxippx"

func Rc4(s []byte) []byte {
	key := []byte(key)
	c, _ := rc4.NewCipher(key)
	c.XORKeyStream(s, s)
	return s
}
