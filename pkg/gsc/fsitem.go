package gsc

import (
	"crypto/md5"
	"io"
	"io/fs"
	"os"
)

type FsItem struct {
	Typ       string
	Path      string
	Target    string
	Endtarget string
	Md5       []byte
	Info      fs.FileInfo
	LinkState string
}

func NewFsItem(path string, info fs.FileInfo) *FsItem {
	typ, err := GetItemType(info)
	CheckErr(err)

	var md5b []byte
	var target, targetType, linkState, endTarget string

	if typ == TYPE_SYMLINK {
		linkState = LINK_STATE_VALID
		target, targetType, err = WalkSymlink(path)
		if err == ErrDeadLink {
			linkState = LINK_STATE_DEAD
		} else if err != nil {
			CheckErr(err)
		}

		for linkState == LINK_STATE_VALID && targetType == TYPE_SYMLINK {
			target, targetType, err = WalkSymlink(target)
			if err == ErrDeadLink {
				linkState = LINK_STATE_DEAD
			}
		}
	}
	if typ == TYPE_FILE {
		file, err := os.Open(path)
		CheckErr(err)
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
	}
}
