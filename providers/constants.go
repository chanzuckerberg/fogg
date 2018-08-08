package providers

const (
	// CustomPluginCacheDir is the directory used by terraform to search for custom providers
	// We default to linux_amd64 since we're running terraform inside of docker
	// We vendor providers here
	// See https://www.terraform.io/docs/configuration/providers.html#third-party-plugins
	CustomPluginCacheDir = "terraform.d/plugin-cache/linux_amd64"
	// PluginCacheDir is the directory where terraform cahes tf approved providers
	// See https://www.terraform.io/docs/configuration/providers.html#provider-plugin-cache
	PluginCacheDir = ".terraform.d/plugin-cache"
)
