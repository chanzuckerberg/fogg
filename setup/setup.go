package setup

import (
	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/plan"
	"github.com/spf13/afero"
)

func Setup(fs afero.Fs, conf *config.Config) error {
	p, err := plan.Eval(conf, false)
	if err != nil {
		return err
	}
	return setupPlugins(fs, p)
}

func setupPlugins(fs afero.Fs, p *plan.Plan) error {
	// log.Debug("setting up plugins")
	// apply := func(name string, plugin *plugins.CustomPlugin) error {
	// 	log.Infof("Setting up plugin %s", name)
	// 	return errs.WrapUserf(plugin.Install(fs, name), "Error applying plugin %s", name)
	// }

	// for pluginName, plugin := range p.Plugins.CustomPlugins {
	// 	err := apply(pluginName, plugin)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// for providerName, provider := range p.Plugins.TerraformProviders {
	// 	err := apply(providerName, provider)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

// // Plugins contains a plan around plugins
// type Plugins struct {
// 	CustomPlugins      map[string]*plugins.CustomPlugin
// 	TerraformProviders map[string]*plugins.CustomPlugin
// }

// // SetCustomPluginsPlan determines the plan for customPlugins
// func (p *Plugins) SetCustomPluginsPlan(customPlugins map[string]*plugins.CustomPlugin) {
// 	p.CustomPlugins = customPlugins
// 	for _, plugin := range p.CustomPlugins {
// 		plugin.SetTargetPath(plugins.CustomPluginDir)
// 	}
// }

// // SetTerraformProvidersPlan determines the plan for customPlugins
// func (p *Plugins) SetTerraformProvidersPlan(terraformProviders map[string]*plugins.CustomPlugin) {
// 	p.TerraformProviders = terraformProviders
// 	for _, plugin := range p.TerraformProviders {
// 		plugin.SetTargetPath(plugins.TerraformCustomPluginCacheDir)
// 	}
// }
