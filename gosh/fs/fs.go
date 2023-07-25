package fs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
)

func DirName(p string) string {
	return filepath.Dir(p)
}

func UserName() string {
	c, err := user.Current()
	if err != nil {
		return ""
	}

	return c.Username
}

func stat(p string) (os.FileInfo, error) {
	return os.Stat(p)
}

func IsExist(p string) bool {
	_, err := stat(p)
	return !errors.Is(err, os.ErrNotExist)
}

func IsFile(p string) bool {
	s, err := stat(p)
	if err != nil {
		return false
	}
	return s.Mode().IsRegular()
}

func IsDir(p string) bool {
	s, err := stat(p)
	if err != nil {
		return false
	}

	return s.Mode().IsDir()
}

func IsAbs(p string) bool {
	return IsExist(p) && filepath.IsAbs(p)
}

func TempFile(dir string, pattern string) (*os.File, error) {
	return os.CreateTemp(dir, pattern)
}

func TempDir(dir string, pattern string) (string, error) {
	return os.MkdirTemp(dir, pattern)
}

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

func WriteFileByte(path string, buf []byte, perm os.FileMode) error {
	dir := DirName(path)
	if !IsExist(dir) {
		if err := os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("dir %s does not exist and has no permission to create it", dir)
		}
	}

	if IsExist(dir) && !IsDir(dir) {
		return fmt.Errorf("dir %s does not a valid directory", dir)
	}

	if !IsExist(path) {
		if _, err := os.Create(path); err != nil {
			return fmt.Errorf("%s does not exist and has not permission to create it", path)
		}
	}

	if !IsFile(path) {
		return fmt.Errorf("%s is not a valid directory", path)
	}

	return os.WriteFile(path, buf, perm)
}

func WriteFile(path string, s string, perm os.FileMode) error {
	return WriteFileByte(path, []byte(s), perm)
}

func WriteTempFile(dir string, pattern string, s string) (string, error) {
	f, err := TempFile(dir, pattern)
	if err != nil {
		return "", err
	}

	_, err = f.Write([]byte(s))
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

func WriteTempDir(dir string, dirPattern string, filePattern string, s string) (string, string, error) {
	dir, err := TempDir(dir, dirPattern)
	if err != nil {
		return "", "", err
	}

	path, err := WriteTempFile(dir, filePattern, s)
	if err != nil {
		return dir, "", err
	}

	return dir, path, nil
}

func Remove(path string, force bool) error {
	if force {
		return os.RemoveAll(path)
	}

	return os.Remove(path)
}
