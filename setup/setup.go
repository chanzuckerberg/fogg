package setup

import (
	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/plugins"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// Setup will download custom plugins and providers
func Setup(fs afero.Fs, conf *v2.Config) error {
	logrus.Debug("setting up plugins")
	cacheDir, err := util.GetFoggCachePath()
	if err != nil {
		return err
	}

	apply := func(name string, plugin *plugins.CustomPlugin) error {
		logrus.Infof("Setting up plugin %s", name)
		return errs.WrapUserf(plugin.Install(fs, name), "Error applying plugin %s", name)
	}

	cache := plugins.GetPluginCache(cacheDir)
	for pluginName, plugin := range conf.Plugins.CustomPlugins {
		plugin.WithTargetPath(plugins.CustomPluginDir).WithCache(cache)
		err := apply(pluginName, plugin)
		if err != nil {
			return err
		}
	}

	for providerName, provider := range conf.Plugins.TerraformProviders {
		provider.WithTargetPath(plugins.TerraformCustomPluginCacheDir).WithCache(cache)
		err := apply(providerName, provider)
		if err != nil {
			return err
		}
	}
	return nil
}
