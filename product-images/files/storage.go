package files

import "io"

// Storage defines the behaviour for file operations
// Implementations may be of the time local disk, or cloud storage, etc
type Storage interface {
	Save(path string, file io.Reader) error
}
