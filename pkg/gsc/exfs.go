package gsc

import (
	"io/fs"
	"os"
)

// ExFS is an fs.FS extension with symlink operations
type ExFS interface {
	fs.FS
	Readlink(name string) (string, error)
	Lstat(name string) (os.FileInfo, error)
}

// ExDirFs is the base ExFS implementation
type ExDirFs string

// DirFs creates an ExDirFS
func DirFs(dir string) ExFS {
	return ExDirFs(dir)
}

// Open file by name
func (f ExDirFs) Open(name string) (fs.File, error) {
	return os.Open(name)
}

// Readlink by name
func (f ExDirFs) Readlink(name string) (string, error) {
	return os.Readlink(name)
}

// Lstat returns a fs.FileInfo for name
func (f ExDirFs) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(name)
}
