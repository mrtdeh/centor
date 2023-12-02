package pluginManager

import (
	PluginKits "github.com/mrtdeh/centor/plugins/assets"
	echo_plugin "github.com/mrtdeh/centor/plugins/echo"
	packageupdater_plugin "github.com/mrtdeh/centor/plugins/packageUpdater"
	time_plugin "github.com/mrtdeh/centor/plugins/time"
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

		pms.AddPlugin(&time_plugin.PluginProvider{
			PluginProps: PluginKits.PluginProps{
				Name: "time-manager",
			},
		})

	})
}
