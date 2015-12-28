package flib

import (
	"crypto/md5"
	"io"
	"os"
)

func GetFileHashCode(fPath string) ([]byte, error) {
	if _, err := os.Stat(fPath); err == nil {
		file, inerr := os.Open(fPath)
		if inerr == nil {
			coder := md5.New()
			io.Copy(coder, file)
			sum := coder.Sum(nil)
			return sum, nil
		} else {
			return nil, inerr
		}
	} else {
		return nil, err
	}
}
