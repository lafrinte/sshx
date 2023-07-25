package b64

import "encoding/base64"

func Encrypt(text []byte) string {
	return base64.StdEncoding.EncodeToString(text)
}

func Decrypt(text string) ([]byte, error) {
	out, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return []byte(""), err
	}

	return out, nil
}
