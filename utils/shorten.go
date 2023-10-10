package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
)

func Shorten(longURL string, numOfChar int) (string, error) {
	h := md5.New()
	_, err := h.Write([]byte(longURL))
	if err != nil {
		return "", err
	}

	d := h.Sum(nil)

	var b bytes.Buffer
	_, err = base64.NewEncoder(base64.RawURLEncoding, &b).Write(d)
	if err != nil {
		return "", err
	}

	return b.String()[:numOfChar], nil
}
