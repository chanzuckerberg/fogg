package util

import (
	"reflect"
	"testing"
)

func Test_ParseRegistryUrl(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    *CodeArtifactRepository
		wantErr bool
	}{
		{
			name: "Valid CodeArtifact Repository URL",
			url:  "https://vincenthsh-123456789.d.codeartifact.ap-southeast-1.amazonaws.com/npm/npm-releases/",
			want: &CodeArtifactRepository{
				Domain:    "vincenthsh",
				AccountId: "123456789",
				Region:    "ap-southeast-1",
				Name:      "npm-releases",
			},
			wantErr: false,
		},
		{
			name:    "Invalid CodeArtifact URL",
			url:     "https://example.com",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRegistryUrl(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRegistryUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseRegistryUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
