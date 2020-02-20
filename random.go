package simpic

import (
	"crypto/rand"
	"fmt"
)

func randomBytes(len int) []byte {
	res := make([]byte, len)
	n, err := rand.Read(res)

	if n < len || err != nil {
		panic(fmt.Sprintf("Unable to generate random bytes. Wanted: %d, got: %d, err: %s", len, n, err))
	}

	return res
}
