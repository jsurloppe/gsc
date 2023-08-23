package gsc

import (
	"reflect"
	"testing"
)

func TestNewPkgItem(t *testing.T) {
	type args struct {
		cat  string
		pkg  string
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    *PkgItem
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPkgItem(tt.args.cat, tt.args.pkg, tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPkgItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPkgItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
