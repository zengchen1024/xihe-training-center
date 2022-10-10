package utils

import (
	"crypto/md5"
	"fmt"
	"time"
)

func Retry(f func() error) (err error) {
	if err = f(); err == nil {
		return
	}

	m := 100 * time.Millisecond
	t := m

	for i := 1; i < 10; i++ {
		time.Sleep(t)
		t += m

		if err = f(); err == nil {
			return
		}
	}

	return
}

func GenMD5(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}
