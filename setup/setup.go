package setup

import (
	"github.com/apex/log"
	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/plan"
	"github.com/chanzuckerberg/fogg/plugins"
	"github.com/spf13/afero"
)

func Setup(fs afero.Fs, conf *config.Config) error {
	p, err := plan.Eval(conf, false, false)
	if err != nil {
		return err
	}
	return setupPlugins(fs, p)
}

func setupPlugins(fs afero.Fs, p *plan.Plan) error {
	log.Debug("setting up plugins")
	apply := func(name string, plugin *plugins.CustomPlugin) error {
		log.Infof("Setting up plugin %s", name)
		return errs.WrapUserf(plugin.Install(fs, name), "Error applying plugin %s", name)
	}

	for pluginName, plugin := range p.Plugins.CustomPlugins {
		err := apply(pluginName, plugin)
		if err != nil {
			return err
		}
	}

	for providerName, provider := range p.Plugins.TerraformProviders {
		err := apply(providerName, provider)
		if err != nil {
			return err
		}
	}
	return nil
}
