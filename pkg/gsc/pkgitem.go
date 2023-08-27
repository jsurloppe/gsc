package gsc

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// PkgItem represents as known as portage database
type PkgItem struct {
	Typ    string
	Path   string
	Target string
	Cat    string
	Pkg    string
	Md5    []byte
	Mtime  time.Time
}

// NewPkgItem creates a PkgItem from a raw line of portage database
func NewPkgItem(cat, pkg, line string) (pkgItem *PkgItem, err error) {
	split := strings.Split(line, " ")
	typ := split[0]

	var path, target string
	var mtime time.Time
	var md5 []byte

	switch typ {
	case TypeFile:
		path = strings.Join(split[1:len(split)-2], " ")
		md5, err = hex.DecodeString(split[len(split)-2])
		CheckErr(err)
		mtimeInt, err := strconv.ParseInt(split[len(split)-1], 10, 64)
		CheckErr(err)
		mtime = time.Unix(mtimeInt, 0)
	case TypeDirectory:
		path = strings.Join(split[1:], " ")
	case TypeSymlink:
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
