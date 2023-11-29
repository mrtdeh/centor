package pluginManager

import (
	echo_plugin "github.com/mrtdeh/centor/plugins/echo"
	PluginKits "github.com/mrtdeh/centor/plugins/kits"
	packageupdater_plugin "github.com/mrtdeh/centor/plugins/packageUpdater"
)

type Config struct {
	PluginKits.Config
}

func Bootstrap(cnf Config) error {
	return PluginKits.Loader(cnf.Config, func(pms *PluginKits.PluginManagerService) {

		pms.AddPlugin(&echo_plugin.PluginProvider{
			PluginProps: PluginKits.PluginProps{
				Name: "echo",
			},
		})

		pms.AddPlugin(&packageupdater_plugin.PluginProvider{
			PluginProps: PluginKits.PluginProps{
				Name: "offline-update",
			},
		})

	})
}
