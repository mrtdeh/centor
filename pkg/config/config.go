package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"reflect"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

// var serviceMap = map[string]service.Service{}

type config struct {
	Name        *string `hcl:"name"`
	Host        *string `hcl:"host" json:"host,omitempty"`
	Port        *uint   `hcl:"port"`
	IsServer    *bool   `hcl:"is_server"`
	IsLeader    *bool   `hcl:"is_leader"`
	ServersAddr *string `hcl:"servers_address"`
	Service     *struct {
		Id   string `hcl:"id"`
		Name string `hcl:"name"`
		Port uint   `hcl:"port"`
	} `hcl:"service,block"`
}

type Service struct {
	Id   string
	Name string
	Port uint
}
type Config struct {
	Name        string
	Host        string
	Port        uint
	IsServer    bool
	IsLeader    bool
	ServersAddr string
	Services    []Service
}

func LoadConfiguration(path string) *Config {
	if path == "" {
		path = "/etc/centor.d/"
	}
	configs, err := loadConfigsFromDir(path)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	cnf, err := compile(configs)
	if err != nil {
		log.Fatalf("Failed to compile configuration: %s", err)
	}

	flag.StringVar(&cnf.Name, "name", cnf.Name, "")
	flag.StringVar(&cnf.Host, "host", cnf.Host, "")
	flag.UintVar(&cnf.Port, "port", cnf.Port, "")
	flag.StringVar(&cnf.ServersAddr, "servers-addr", cnf.ServersAddr, "")
	flag.BoolVar(&cnf.IsServer, "server", cnf.IsServer, "")
	flag.BoolVar(&cnf.IsLeader, "leader", cnf.IsLeader, "")
	flag.Parse()

	cb, _ := json.MarshalIndent(cnf, "", " ")
	fmt.Printf("configs : %s\n", cb)

	return cnf
}

func compile(configs []config) (*Config, error) {
	cnf := &Config{}
	for _, c := range configs {

		cnf.Name = check(c.Name, cnf.Name).(string)
		cnf.Host = check(c.Host, cnf.Host).(string)
		cnf.Port = check(c.Port, cnf.Port).(uint)
		cnf.IsServer = check(c.IsServer, cnf.IsServer).(bool)
		cnf.IsLeader = check(c.IsLeader, cnf.IsLeader).(bool)
		cnf.ServersAddr = check(c.ServersAddr, cnf.ServersAddr).(string)

		if c.Service != nil {
			cnf.Services = append(cnf.Services, Service(*c.Service))
		}
	}
	return cnf, nil
}

func loadConfigsFromDir(directory string) (cnf []config, err error) {
	var configs []config

	files, err := filepath.Glob(filepath.Join(directory, "*.hcl"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {

		var config config
		err = hclsimple.DecodeFile(file, nil, &config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func check(value any, default_value any) any {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return default_value
	}
	return v.Elem().Interface()
}
