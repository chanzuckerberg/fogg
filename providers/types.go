package providers

// TypeProviderFormat is the provider format such as binary, zip, tar
type TypeProviderFormat string

const (
	// TypeProviderFormatBinary is a binary provider format
	TypeProviderFormatBinary TypeProviderFormat = "none"
	// TypeProviderFormatTar is a tar archived provider
	TypeProviderFormatTar TypeProviderFormat = "tar"
)

// CustomProvider is a custom terraform provider
type CustomProvider struct {
	URL    string             `json:"url,omitempty" validate:"required"`
	Format TypeProviderFormat `json:"format,omitempty" validate:"required"`
}
