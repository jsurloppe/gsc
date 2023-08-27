package gsc

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// TypeFile represents a file in portage
const TypeFile = "obj"

// TypeDirectory represents a directory in portage
const TypeDirectory = "dir"

// TypeSymlink represents a symlink in portage
const TypeSymlink = "sym"

// LinkStateValid is the valid state of a symlink
const LinkStateValid = "valid"

// LinkStateDead is the invalid state of a symlink
const LinkStateDead = "dead"

// ErrDeadLink is the error in case of dead link
var ErrDeadLink = errors.New(LinkStateDead)

// CheckErr check error, log and exit if err is not nil
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// GetItemType returns the portage type from a real file (info)
func GetItemType(info fs.FileInfo) (string, error) {
	mode := info.Mode()
	switch {
	case mode.IsRegular():
		return TypeFile, nil
	case mode.IsDir():
		return TypeDirectory, nil
	case mode&os.ModeSymlink != 0:
		return TypeSymlink, nil
	default:
		return "", errors.New("invalid mode")
	}
}

// LoadIgnorePatterns loads patterns from a .gscignore file.
func LoadIgnorePatterns(filesystem fs.FS, filename string) []string {
	file, err := filesystem.Open(filename)
	CheckErr(err)
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pattern := scanner.Text()
		if pattern == "" || pattern[0] == '#' { // empty line or comment
			continue
		}
		patterns = append(patterns, pattern)
	}

	err = scanner.Err()
	CheckErr(err)

	return patterns
}

// WalkSymlink walks through a symlink
func WalkSymlink(filesystem ExFS, path string) (string, string, error) {
	var typ string
	dir := filepath.Dir(path)
	target, err := filesystem.Readlink(path)
	CheckErr(err)

	if !filepath.IsAbs(target) {
		target = filepath.Join(dir, target)
	}

	info, err := filesystem.Lstat(target)
	if err != nil && os.IsNotExist(err) {
		return target, "", ErrDeadLink
	}

	typ, err = GetItemType(info)
	CheckErr(err)

	return target, typ, nil
}

// BuildPackageMap builds an map with path as key and PkgItem as value
// of portage database content
func BuildPackageMap(filesystem fs.FS, dbPath string) (entries map[string]*PkgItem) {
	entries = make(map[string]*PkgItem)
	matches, err := fs.Glob(filesystem, filepath.Join(dbPath, "/*/*/CONTENTS"))
	CheckErr(err)
	for _, match := range matches {
		split := strings.Split(match, "/")
		cat := split[len(split)-3]
		pkg := split[len(split)-2]
		fd, err := filesystem.Open(match)
		CheckErr(err)
		defer fd.Close()

		scanner := bufio.NewScanner(fd)
		for scanner.Scan() {
			line := scanner.Text()
			entry, err := NewPkgItem(cat, pkg, line)
			CheckErr(err)
			entries[entry.Path] = entry
		}
	}
	return
}

// IsExcluded checks if a path is excluded from check
func IsExcluded(ignorePatterns []string, path string) bool {
	for _, pattern := range ignorePatterns {
		match, err := filepath.Match(filepath.Clean(pattern), path)
		CheckErr(err)
		if match || strings.HasPrefix(path, filepath.Join(pattern, string(filepath.Separator))) {
			return true
		}
	}
	return false
}
