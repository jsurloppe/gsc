package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jsurloppe/gsc/pkg/gsc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

func init() {
}

var rootPath string
var ignoreFile = flag.StringP("ignore-file", "i", "configs/gscignore", "file containing pattern to ignore")
var versionFlag = flag.BoolP("version", "V", false, "show version and exit")

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

		log.SetLevel(log.DebugLevel)
		rootPath = args[0]

		ignorePatterns := gsc.LoadIgnorePatterns(*ignoreFile)

		pkgItems := gsc.BuildPackageMap()
		log.Debug("db built")

		files := make(map[string]*gsc.PkgItem)

		err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if gsc.IsExcluded(ignorePatterns, path) {
				switch {
				case d.IsDir():
					return fs.SkipDir
				default:
					return nil
				}
			}

			info, err := d.Info()
			gsc.CheckErr(err)
			fsItem := gsc.NewFsItem(path, info)

			var pkgItem *gsc.PkgItem
			var ok bool

			if pkgItem, ok = pkgItems[fsItem.Path]; ok {
				files[pkgItem.Path] = pkgItem
			}
			// } else if pkgItem, ok = pkgItems[fsItem.target]; ok {
			// 	files[pkgItem.target] = pkgItem
			// }

			switch {
			case fsItem.Typ == gsc.TYPE_SYMLINK && fsItem.LinkState == gsc.LINK_STATE_DEAD:
				log.WithFields(log.Fields{"file": fsItem.Path}).Error("dead link")
			case !ok:
				log.WithFields(log.Fields{"file": fsItem.Path, "type": fsItem.Typ, "target": fsItem.Target}).Warn("file not in db")
			case fsItem.Typ != pkgItem.Typ:
				log.WithFields(log.Fields{"file": fsItem.Path, "fileType": fsItem.Typ, "pkgType": pkgItem.Typ, "target": fsItem.Target}).Error("different file type")
			case fsItem.Typ == gsc.TYPE_FILE && !bytes.Equal(fsItem.Md5, pkgItem.Md5):
				log.WithFields(log.Fields{"file": fsItem.Path}).Warn("modified file")
				// default:
				// 	log.WithFields(log.Fields{"file": fsItem.path}).Info("ok")
			}
			return nil
		})
		gsc.CheckErr(err)

		for _, pkgItem := range pkgItems {
			if strings.HasPrefix(pkgItem.Path, rootPath) && !gsc.IsExcluded(ignorePatterns, pkgItem.Path) {
				if _, ok := files[pkgItem.Path]; !ok {
					log.WithFields(log.Fields{"file": pkgItem.Path, "cat": pkgItem.Cat, "pkg": pkgItem.Pkg}).Error("missing file")
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
