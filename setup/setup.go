package setup

import (
	"github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/plugins"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func Setup(fs afero.Fs, conf *v1.Config) error {
	log.Debug("setting up plugins")
	apply := func(name string, plugin *plugins.CustomPlugin) error {
		log.Infof("Setting up plugin %s", name)
		return errs.WrapUserf(plugin.Install(fs, name), "Error applying plugin %s", name)
	}

	for pluginName, plugin := range conf.Plugins.CustomPlugins {
		plugin.SetTargetPath(plugins.CustomPluginDir)
		err := apply(pluginName, plugin)
		if err != nil {
			return err
		}
	}

	for providerName, provider := range conf.Plugins.TerraformProviders {
		provider.SetTargetPath(plugins.TerraformCustomPluginCacheDir)
		err := apply(providerName, provider)
		if err != nil {
			return err
		}
	}
	return nil
}
