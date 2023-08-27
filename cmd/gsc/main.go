/*
GSC check Gentoo system consistency

Gentoo System Check (GSC) use your local package database to check part of your system
It reports files that were added/modified/deleted that differ from the database.

It's normal to have files modified or added, but sometimes you ends with orphan files,
dead symlinks, manually installed things that you forgot, ect...

This tool helps to have a quick overview of your system and helps you to keep it clean.

Usage:

	gsc [flags] [path]

The flags are:

	    --json
	        Print the logs as json, easier for parsing.
	    -i
			Use an ignore file to ignore some files or path.
		-V
			Print the version and exit
*/
package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/jsurloppe/gsc/pkg/gsc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

func init() {
}

var rootPath string
var ignoreFile = flag.StringP("ignore-file", "i", "", "file containing pattern to ignore")
var versionFlag = flag.BoolP("version", "V", false, "show version and exit")
var jsonFlag = flag.Bool("json", false, "show version and exit")

func skipOnPermError(err error) bool {
	if os.IsPermission(err) {
		log.Error(err)
		return true
	}
	gsc.CheckErr(err)
	return false
}

var rootCmd = &cobra.Command{
	Use:   "gsc [path]",
	Short: "Run Gentoo System Check on [path]",
	Run: func(cmd *cobra.Command, args []string) {
		if *versionFlag {
			fmt.Println("gsc: v0.1 -- HEAD")
			os.Exit(0)
		} else if len(args) != 1 {
			fmt.Println("gsc: A path is required")
			os.Exit(1)
		}

		if *jsonFlag {
			log.SetFormatter(&log.JSONFormatter{})
		} else {
			log.SetFormatter(&log.TextFormatter{})
		}
		log.SetLevel(log.DebugLevel)
		rootPath = args[0]

		filesystem := gsc.ExDirFs(rootPath)

		var ignorePatterns []string
		if *ignoreFile != "" {
			ignorePatterns = gsc.LoadIgnorePatterns(filesystem, *ignoreFile)
		}

		pkgItems := gsc.BuildPackageMap(filesystem, "/var/db/pkg/")
		log.Debug("db built")

		files := make(map[string]*gsc.PkgItem)

		err := fs.WalkDir(filesystem, rootPath, func(path string, d fs.DirEntry, err error) error {
			if gsc.IsExcluded(ignorePatterns, path) {
				switch {
				case d.IsDir():
					return fs.SkipDir
				default:
					return nil
				}
			}
			if skipOnPermError(err) {
				return nil
			}

			info, err := d.Info()
			gsc.CheckErr(err)

			fsItem, err := gsc.NewFsItem(filesystem, path, info)
			if skipOnPermError(err) {
				return nil
			}

			var pkgItem *gsc.PkgItem
			var ok bool

			if pkgItem, ok = pkgItems[fsItem.Path]; ok {
				files[pkgItem.Path] = pkgItem
			}

			switch {
			case fsItem.Typ == gsc.TypeSymlink && fsItem.LinkState == gsc.LinkStateDead:
				log.WithFields(log.Fields{"file": fsItem.Path}).Error("dead link")
			case !ok:
				log.WithFields(log.Fields{"file": fsItem.Path, "type": fsItem.Typ, "target": fsItem.Target}).Warn("file not in db")
			case fsItem.Typ != pkgItem.Typ:
				log.WithFields(log.Fields{"file": fsItem.Path, "fileType": fsItem.Typ, "pkgType": pkgItem.Typ, "target": fsItem.Target}).Error("different file type")
			case fsItem.Typ == gsc.TypeFile && !bytes.Equal(fsItem.Md5, pkgItem.Md5):
				log.WithFields(log.Fields{"file": fsItem.Path}).Warn("modified file")
			}
			return nil
		})
		gsc.CheckErr(err)

		for _, pkgItem := range pkgItems {
			if strings.HasPrefix(pkgItem.Path, rootPath) && !gsc.IsExcluded(ignorePatterns, pkgItem.Path) {
				if _, ok := files[pkgItem.Path]; !ok {
					log.WithFields(log.Fields{"file": pkgItem.Path, "cat": pkgItem.Cat, "pkg": pkgItem.Pkg}).Error("missing file or access problem")
				}
			}
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
