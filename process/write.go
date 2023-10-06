package process

import "os"

func Write(filename string, contents []byte) error {
	return os.WriteFile(filename, contents, os.ModePerm)
}
