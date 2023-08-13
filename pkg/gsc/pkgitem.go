package gsc

import (
	"encoding/hex"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type PkgItem struct {
	Typ    string
	Path   string
	Target string
	Cat    string
	Pkg    string
	Md5    []byte
}

func NewPkgItem(cat, pkg, line string) (*PkgItem, error) {
	split := strings.Split(line, " ")
	typ := split[0]
	switch typ {
	case TYPE_FILE:
		path := strings.Join(split[1:len(split)-2], " ")
		md5, err := hex.DecodeString(split[len(split)-2])
		if err != nil {
			log.Fatal(err)
		}
		return &PkgItem{
			Typ:  typ,
			Path: path,
			Cat:  cat,
			Pkg:  pkg,
			Md5:  md5,
		}, nil
	case TYPE_DIRECTORY:
		path := strings.Join(split[1:], " ")
		return &PkgItem{
			Typ:  typ,
			Path: path,
			Cat:  cat,
			Pkg:  pkg,
		}, nil
	case TYPE_SYMLINK:
		path := strings.Join(split[1:len(split)-1], " ")
		symSlice := strings.Split(path, " -> ")
		path = symSlice[0]
		target := symSlice[1]

		return &PkgItem{
			Typ:    typ,
			Path:   path,
			Target: target,
			Cat:    cat,
			Pkg:    pkg,
		}, nil
	default:
		return nil, fmt.Errorf("invalid type %s/%s: %s", cat, pkg, line)
	}
}
