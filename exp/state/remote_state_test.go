package state

import (
	"reflect"
	"testing"
)

func Test_collectRemoteStateReferences(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"main", args{"testdata"}, []string{"one", "two", "three", "four", "five", "six"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := collectRemoteStateReferences(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("collectRemoteStateReferences() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collectRemoteStateReferences() = %v, want %v", got, tt.want)
			}
		})
	}
}
