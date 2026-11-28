//go:build !linux && !windows

package localfile

import (
	"errors"
	"os"
)

func Lock(string) (*os.File, error) {
	return nil, errors.New("file locking requires Linux or Windows")
}
