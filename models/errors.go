package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"slices"
)

var (
	ErrEmailTaken = errors.New("models: email address is already in use")
	ErrNotFound   = errors.New("models: resource could not be found")
)

type FileError struct {
	Issue string
}

func (fe FileError) Error() string {
	return fmt.Sprintf("invalid file: %v", fe.Issue)
}

func checkContentType(r io.Reader, allowedTypes []string) ([]byte, error) {
	testBytes := make([]byte, 512)
	n, err := r.Read(testBytes)
	if err != nil {
		return nil, fmt.Errorf("checking content type: %w", err)
	}

	contentType := http.DetectContentType(testBytes)
	if !slices.Contains(allowedTypes, contentType) {
		return nil, FileError{
			Issue: fmt.Sprintf("invalid content type: %v", contentType),
		}
	}

	return testBytes[:n], nil
}

func checkExtension(filename string, allowedExtensions []string) error {
	if !hasExtension(filename, allowedExtensions) {
		return FileError{
			Issue: fmt.Sprintf("invalid extension: %v", filepath.Ext(filename)),
		}
	}

	return nil
}
