package gsc

import (
	"encoding/hex"
	"io/fs"
	"os"
	"reflect"
	"testing"
	"testing/fstest"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestBuildPackageMap(t *testing.T) {
	type args struct {
		filesystem fs.FS
		dbPath     string
	}

	fsData := map[string]*fstest.MapFile{
		"var/db/pkg/sys-apps/openrc-0.47.1/CONTENTS": { // leading slash doesn't work with MapFS
			Data: []byte(`dir /lib64
obj /lib64/libeinfo.so.1 7697b9de8d325dfb30b14d747610c014 1691333557
sym /lib64/libeinfo.so -> libeinfo.so.1 1691333556
`),
		},
	}

	md5, _ := hex.DecodeString("7697b9de8d325dfb30b14d747610c014")

	tests := []struct {
		name string
		args args
		want map[string]*PkgItem
	}{
		{
			name: "ContentParse",
			args: args{
				filesystem: fstest.MapFS(fsData),
				dbPath:     "var/db/pkg/",
			},
			want: map[string]*PkgItem{
				"/lib64/libeinfo.so.1": {
					Typ:   TypeFile,
					Path:  "/lib64/libeinfo.so.1",
					Cat:   "sys-apps",
					Pkg:   "openrc-0.47.1",
					Md5:   md5,
					Mtime: time.Unix(1691333557, 0),
				},
				"/lib64/libeinfo.so": {
					Typ:    TypeSymlink,
					Path:   "/lib64/libeinfo.so",
					Target: "libeinfo.so.1",
					Cat:    "sys-apps",
					Pkg:    "openrc-0.47.1",
					Mtime:  time.Unix(1691333556, 0),
				},
				"/lib64": {
					Typ:  TypeDirectory,
					Path: "/lib64",
					Cat:  "sys-apps",
					Pkg:  "openrc-0.47.1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildPackageMap(tt.args.filesystem, tt.args.dbPath)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("BuildPackageMap() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestLoadIgnorePatterns(t *testing.T) {
	type args struct {
		filesystem fs.FS
		filename   string
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "SyntaxCheck",
			args: args{
				filesystem: os.DirFS("testdata"),
				filename:   "gscignore",
			},
			want: []string{
				"/etc/sysctl.conf",
				"/etc/portage/",
				"/usr/bin/GET",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadIgnorePatterns(tt.args.filesystem, tt.args.filename); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadIgnorePatterns() = %v, want %v", got, tt.want)
			}
		})
	}
}
