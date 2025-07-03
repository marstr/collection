package collection

import (
	"context"
	"os"
	"path/filepath"
)

// DirectoryOptions is a means of configuring a `Directory` instance to including various children in its enumeration without
// supplying a `Where` clause later.
type DirectoryOptions uint

// These constants define all of the supported options for configuring a `Directory`
const (
	DirectoryOptionsExcludeFiles = 1 << iota
	DirectoryOptionsExcludeDirectories
	DirectoryOptionsRecursive
)

// Directory treats a filesystem path as a collection of filesystem entries, specifically a collection of directories and files.
type Directory struct {
	Location string
	Options  DirectoryOptions
}

func (d Directory) applyOptions(loc string, info os.FileInfo) bool {
	if info.IsDir() && (d.Options&DirectoryOptionsExcludeDirectories) != 0 {
		return false
	}

	if !info.IsDir() && d.Options&DirectoryOptionsExcludeFiles != 0 {
		return false
	}

	return true
}

// Enumerate lists the items in a `Directory`
func (d Directory) Enumerate(ctx context.Context) Enumerator[string] {
	results := make(chan string)

	go func() {
		defer close(results)

		filepath.Walk(d.Location, func(currentLocation string, info os.FileInfo, openErr error) (err error) {
			if openErr != nil {
				err = openErr
				return
			}

			if d.Location == currentLocation {
				return
			}

			if info.IsDir() && d.Options&DirectoryOptionsRecursive == 0 {
				err = filepath.SkipDir
			}

			if d.applyOptions(currentLocation, info) {
				select {
				case results <- currentLocation:
					// Intentionally Left Blank
				case <-ctx.Done():
					err = ctx.Err()
				}
			}

			return
		})
	}()

	return results
}
