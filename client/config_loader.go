package main

import (
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

func loadConfigsFromDir(directory string) ([]Config, error) {
	var configs []Config

	files, err := filepath.Glob(filepath.Join(directory, "*.hcl"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		var config Config
		err := hclsimple.DecodeFile(file, nil, &config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func parseServices(configs []Config) (services []Service) {
	for _, c := range configs {
		if c.Service != nil {
			serviceMap[c.Service.Id] = Service{
				Name: c.Service.Name,
				Id:   c.Service.Id,
				Port: c.Service.Port,
			}
		}
	}
	return
}
