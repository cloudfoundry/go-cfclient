package ios

import (
	"io"
	"os"
)

func Close(c io.Closer) {
	if c != nil {
		_ = c.Close()
	}
}

func CleanupTempFile(f *os.File) {
	if f != nil {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}
}
