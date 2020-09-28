package state

import (
	"testing"

	"github.com/stretchr/testify/require"
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
			r := require.New(t)
			got, err := collectRemoteStateReferences(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("collectRemoteStateReferences() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			r.ElementsMatch(got, tt.want)
		})
	}
}
