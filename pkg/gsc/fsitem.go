package gsc

import (
	"crypto/md5"
	"io"
	"io/fs"
	"os"
)

// FsItem represents a file on the filesystem
type FsItem struct {
	Typ       string      // object type
	Path      string      // real path
	Target    string      // target, only if typ is symlink
	Endtarget string      // last no (valid) symlink recursive target
	Md5       []byte      // checksum of the file if typ is obj (file)
	Info      fs.FileInfo // file info
	LinkState string      // state of the link if symlink
}

// NewFsItem create new FsItem from path and info
func NewFsItem(filesystem ExFS, path string, info fs.FileInfo) (*FsItem, error) {
	typ, err := GetItemType(info)
	CheckErr(err)

	var md5b []byte
	var target, targetType, linkState, endTarget string

	if typ == TypeSymlink {
		linkState = LinkStateValid
		target, targetType, err = WalkSymlink(filesystem, path)
		if err == ErrDeadLink {
			linkState = LinkStateDead
		} else if err != nil {
			CheckErr(err)
		}

		for linkState == LinkStateValid && targetType == TypeSymlink {
			target, targetType, err = WalkSymlink(filesystem, target)
			if err == ErrDeadLink {
				linkState = LinkStateDead
			}
		}
	}
	if typ == TypeFile {
		file, err := os.Open(path)
		if os.IsPermission(err) {
			return nil, err
		}
		defer file.Close()

		hash := md5.New()
		_, err = io.Copy(hash, file)
		CheckErr(err)
		md5b = hash.Sum(nil)
	}
	return &FsItem{
		Typ:       typ,
		Path:      path,
		Target:    target,
		Endtarget: endTarget,
		Md5:       md5b,
		Info:      info,
		LinkState: linkState,
	}, nil
}
