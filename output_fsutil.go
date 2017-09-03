package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// File and directory permission.
const (
	defaultFilePermissions      = 0600
	defaultDirectoryPermissions = 0760
	defaultBackupPermissions    = 0500
)

// Base struct for custom errors.
type baseError struct {
	message string
}

func (be baseError) Error() string {
	return be.message
}

type cannotOpenFileError struct {
	baseError
}

func newCannotOpenFileError(fname string) *cannotOpenFileError {
	return &cannotOpenFileError{baseError{message: "Cannot open file: " + fname}}
}

type notDirectoryError struct {
	baseError
}

func reportInternalError(err error) {
	fmt.Fprintf(os.Stderr, "log4g internal error: %s\n", err)
}

// fileFilter is a filtering criteria function for '*os.File'.
// Must return 'false' to set aside the given file.
type fileFilter func(os.FileInfo, *os.File) bool

// filePathFilter is a filtering creteria function for file path.
// Must return 'false' to set aside the given file.
type filePathFilter func(filePath string) bool

func isRegular(m os.FileMode) bool {
	return m&os.ModeType == 0
}

// getDirFilePaths return full paths of the files located in the directory.
// Remark: Ignores files for which fileFilter returns false.
func getDirFilePaths(dirPath string, fpFilter filePathFilter, pathIsName bool) ([]string, error) {
	dfi, err := os.Open(dirPath)
	if err != nil {
		return nil, newCannotOpenFileError("Cannot open directory " + dirPath)
	}
	defer dfi.Close()

	var absDirPath string
	if !filepath.IsAbs(dirPath) {
		absDirPath, err = filepath.Abs(dirPath)
		if err != nil {
			return nil, fmt.Errorf("cannot get absolute path of directory: %s", err.Error())
		}
	} else {
		absDirPath = dirPath
	}

	// TODO: check if dirPath is really directory.
	// Size of read buffer (i.e. chunk of items read at a time).
	rbs := 2 << 5
	filePaths := []string{}

	var fp string
L:
	for {
		// Read directory entities by reasonable chuncks
		// to prevent overflows on big number of files.
		fis, e := dfi.Readdir(rbs)
		switch e {
		// It's OK.
		case nil:
			// Do nothing, just continue cycle.
		case io.EOF:
			break L
			// Indicate that something went wrong.
		default:
			return nil, e
		}
		// THINK: Maybe, use async running.
		for _, fi := range fis {
			// NB: Should work on every Windows and non-Windows OS.
			if isRegular(fi.Mode()) {
				if pathIsName {
					fp = fi.Name()
				} else {
					// Build full path of a file.
					fp = filepath.Join(absDirPath, fi.Name())
				}
				// Check filter condition.
				if fpFilter != nil && !fpFilter(fp) {
					continue
				}
				filePaths = append(filePaths, fp)
			}
		}
	}
	return filePaths, nil
}

// tryRemoveFile gives a try removing the file
// only ignoring an error when the file does not exist.
func tryRemoveFile(filePath string) (err error) {
	err = os.Remove(filePath)
	if os.IsNotExist(err) {
		err = nil
		return
	}
	return
}
