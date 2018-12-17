package file

import (
	"fmt"
	"os"
)

func walkFunc(path string, info os.FileInfo, err error) error {
	fmt.Printf("%s\n", path)
	return nil
}
