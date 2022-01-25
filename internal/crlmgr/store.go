package crlmgr

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Stat the modification time from the target file.
func Stat(target string) (time.Time, error) {
	i, err := os.Stat(target)
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to get target file info: %w", err)
	}

	return i.ModTime(), nil
}

// Store the contents of body in the target file.
func Store(target string, ts time.Time, body io.Reader) (err error) {
	f, drr := os.Create(target)
	if drr != nil {
		return fmt.Errorf("unable to create target file: %w", drr)
	}

	defer func() { //nolint:gosec // I'm handling it silly linter.
		if drr := f.Close(); drr != nil {
			err = fmt.Errorf("unable to close file: %w", drr)
		}
	}()

	if _, drr := io.Copy(f, body); drr != nil {
		return fmt.Errorf("unable to write target contents: %w", drr)
	}

	if drr := os.Chtimes(target, ts, ts); drr != nil {
		return fmt.Errorf("unable to set modification time: %w", drr)
	}

	return err
}
