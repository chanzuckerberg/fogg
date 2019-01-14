package plugins

const (
	// TerraformPluginCacheDir is the directory where terraform caches tf approved providers
	// See https://www.terraform.io/docs/configuration/providers.html#provider-plugin-cache
	// TerraformPluginCacheDir = ".terraform.d/plugin-cache"

	// TerraformCustomPluginCacheDir is the directory used by terraform to search for custom providers
	// We vendor providers here
	// See https://www.terraform.io/docs/configuration/providers.html#third-party-plugins
	TerraformCustomPluginCacheDir = "terraform.d/plugins/{{.OS}}_{{.Arch}}"
	// CustomPluginDir where we place custom binaries
	CustomPluginDir = ".fogg/bin"
)
