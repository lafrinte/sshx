package fs

import (
	"io"
	"os"
)

func ReadFileByte(p string) ([]byte, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func ReadFile(p string) (string, error) {
	buf, err := ReadFileByte(p)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
