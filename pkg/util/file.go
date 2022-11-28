package util

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// GetFileNameSize get file name and size info
func GetFileNameSize(path string) (string, int64, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return "", 0, err
	}

	statInfo, err := f.Stat()
	if err != nil {
		return "", 0, err
	}

	return statInfo.Name(), statInfo.Size(), nil
}

// SHA512File sha512 a file
func SHA512File(path string) (string, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return "", err
	}

	fullSHA512 := sha512.New()
	for {
		data := make([]byte, 1024*1024)
		n, err := file.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		if n <= 0 {
			break
		}

		fullSHA512.Write(data[0:n])
	}
	sha512Result := fullSHA512.Sum(nil)
	return fmt.Sprintf("%x", sha512Result), nil
}

// FileIsExists file is exists
func FileIsExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}

	return true
}

// ReadChunk read chunk from file
func ReadChunk(path string, index int64, size int64) ([]byte, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	if _, err := f.Seek(index, 0); err != nil {
		return nil, err
	}

	data := make([]byte, size)
	readSize, err := f.Read(data)
	if err != nil {
		return nil, err
	}

	if int64(readSize) != size {
		return nil, errors.New("system error")
	}

	return data, nil
}
