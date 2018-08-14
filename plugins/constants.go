package plugins

const (
	// TerraformPluginCacheDir is the directory where terraform caches tf approved providers
	// See https://www.terraform.io/docs/configuration/providers.html#provider-plugin-cache
	TerraformPluginCacheDir = ".terraform.d/plugin-cache"

	// CustomPluginDir where we place custom binaries
	CustomPluginDir = ".bin"
)
