package domain

import "io"

// UploadedFile represents an uploaded file with its metadata.
type UploadedFile struct {
	Filename string
	Size     int64
	Reader   io.Reader
}
