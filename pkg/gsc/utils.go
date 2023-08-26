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

const TYPE_FILE = "obj"
const TYPE_DIRECTORY = "dir"
const TYPE_SYMLINK = "sym"
const LINK_STATE_VALID = "valid"
const LINK_STATE_DEAD = "dead"

var ErrDeadLink = errors.New(LINK_STATE_DEAD)

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetItemType(info fs.FileInfo) (string, error) {
	mode := info.Mode()
	switch {
	case mode.IsRegular():
		return TYPE_FILE, nil
	case mode.IsDir():
		return TYPE_DIRECTORY, nil
	case mode&os.ModeSymlink != 0:
		return TYPE_SYMLINK, nil
	default:
		return "", errors.New("invalid mode")
	}
}

// loadIgnorePatterns loads patterns from a .gsyscheckignore file.
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

func IsExcluded(filesystem fs.FS, ignorePatterns []string, path string) bool {
	for _, pattern := range ignorePatterns {
		match, err := filepath.Match(filepath.Clean(pattern), path)
		CheckErr(err)
		if match || strings.HasPrefix(path, filepath.Join(pattern, string(filepath.Separator))) {
			return true
		}
	}
	return false
}
