package config

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/mrtdeh/centor/pkg/service"
)

var serviceMap = map[string]service.Service{}

type config struct {
	Name       *string
	Host       *string
	Port       *string
	IsServer   *bool
	ServerAddr *string

	Service *ServiceConfig `hcl:"service,block"`
}

type ServiceConfig struct {
	Id   string `hcl:"id"`
	Name string `hcl:"name"`
	Port uint   `hcl:"port"`
}

func LoadConfiguration() *config {
	var c *config = &config{}

	// configs, err := loadConfigsFromDir("/etc/centor.d/")
	// if err != nil {
	// 	log.Fatalf("Failed to load configuration: %s", err)
	// }
	// parseServices(configs)

	c.Name = flag.String("n", "", "server host")
	c.Host = flag.String("h", "0.0.0.0", "server host")
	c.Port = flag.String("p", "10000", "server port")
	c.IsServer = flag.Bool("server", false, "is server")
	c.ServerAddr = flag.String("server-addr", "", "server address for dialing")
	flag.Parse()

	if c.Name == nil {
		log.Fatal("do not set your name")
	}

	return c
}

func loadConfigsFromDir(directory string) ([]config, error) {
	var configs []config

	files, err := filepath.Glob(filepath.Join(directory, "*.hcl"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		var config config
		err := hclsimple.DecodeFile(file, nil, &config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func parseServices(configs []config) (services []service.Service) {
	for _, c := range configs {
		if c.Service != nil {
			serviceMap[c.Service.Id] = service.Service{
				Name: c.Service.Name,
				Id:   c.Service.Id,
				Port: c.Service.Port,
			}
		}
	}
	return
}
