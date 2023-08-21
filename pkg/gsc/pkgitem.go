package gsc

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type PkgItem struct {
	Typ    string
	Path   string
	Target string
	Cat    string
	Pkg    string
	Md5    []byte
	Mtime  time.Time
}

func NewPkgItem(cat, pkg, line string) (pkgItem *PkgItem, err error) {
	split := strings.Split(line, " ")
	typ := split[0]

	var path, target string
	var mtime time.Time
	var md5 []byte

	switch typ {
	case TYPE_FILE:
		path = strings.Join(split[1:len(split)-2], " ")
		md5, err = hex.DecodeString(split[len(split)-2])
		CheckErr(err)
		mtimeInt, err := strconv.ParseInt(split[len(split)-1], 10, 64)
		CheckErr(err)
		mtime = time.Unix(mtimeInt, 0)
	case TYPE_DIRECTORY:
		path = strings.Join(split[1:], " ")
	case TYPE_SYMLINK:
		path = strings.Join(split[1:len(split)-1], " ")
		symSlice := strings.Split(path, " -> ")
		path = symSlice[0]
		target = symSlice[1]
		mtimeInt, err := strconv.ParseInt(split[len(split)-1], 10, 64)
		CheckErr(err)
		mtime = time.Unix(mtimeInt, 0)
	default:
		return nil, fmt.Errorf("invalid type %s/%s: %s", cat, pkg, line)
	}

	return &PkgItem{
		Typ:    typ,
		Path:   path,
		Target: target,
		Md5:    md5,
		Cat:    cat,
		Pkg:    pkg,
		Mtime:  mtime,
	}, nil
}
