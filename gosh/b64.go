package gosh

import "encoding/base64"

func B64Encrypt(text []byte) string {
    return base64.StdEncoding.EncodeToString(text)
}

func B64Decrypt(text string) ([]byte, error) {
    out, err := base64.StdEncoding.DecodeString(text)
    if err != nil {
        return []byte(""), err
    }

    return out, nil
}
