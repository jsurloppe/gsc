package gsc

import (
	"io/fs"
	"os"
)

type PkgFS interface {
	fs.FS
}

type ExFS interface {
	fs.FS
	Readlink(name string) (string, error)
	Lstat(name string) (os.FileInfo, error)
}

type ExDirFs string

func DirFs(dir string) ExFS {
	return ExDirFs(dir)
}

func (f ExDirFs) Open(name string) (fs.File, error) {
	return os.Open(name)
}

func (f ExDirFs) Readlink(name string) (string, error) {
	return os.Readlink(name)
}

func (f ExDirFs) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(name)
}
