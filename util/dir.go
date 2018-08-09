package util

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// CreateDirIfNotExists will create a directory (and all intermediate directories) if they do not exist
func CreateDirIfNotExists(target string, dest afero.Fs) error {
	if _, err := dest.Stat(target); err != nil {
		if err := dest.MkdirAll(target, 0755); err != nil {
			return errors.Wrapf(err, "could not create directory at %s", target)
		}
	}
	return nil
}
