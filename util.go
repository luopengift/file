package file

import (
	"fmt"
	"io"
	"os"
)

func walkFunc(path string, info os.FileInfo, err error) error {
	fmt.Printf("%s\n", path)
	return nil
}

// Append append file
func Append(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
