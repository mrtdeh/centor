package pluginManager

import (
	PluginKits "github.com/mrtdeh/centor/plugins/assets"
	echo_plugin "github.com/mrtdeh/centor/plugins/echo"
	installer_plugin "github.com/mrtdeh/centor/plugins/installer"
	system_plugin "github.com/mrtdeh/centor/plugins/system"
	timeSyncer_plugin "github.com/mrtdeh/centor/plugins/time_syncer"
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

		pms.AddPlugin(&installer_plugin.PluginProvider{
			PluginProps: PluginKits.PluginProps{
				Name: "installer",
			},
		})

		pms.AddPlugin(&timeSyncer_plugin.PluginProvider{
			PluginProps: PluginKits.PluginProps{
				Name: "time-syncer",
			},
		})

		pms.AddPlugin(&system_plugin.PluginProvider{
			PluginProps: PluginKits.PluginProps{
				Name: "system",
			},
		})

	})
}
