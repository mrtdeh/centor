package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type config struct {
	Name               *string `hcl:"name"`
	Host               *string `hcl:"host" json:"host,omitempty"`
	Port               *uint   `hcl:"port"`
	IsServer           *bool   `hcl:"is_server"`
	IsLeader           *bool   `hcl:"is_leader"`
	ServersAddr        *string `hcl:"servers_address"`
	PrimaryServersAddr *string `hcl:"primaries_address"`
	Service            *struct {
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
	Name               string    // id of the agent
	Host               string    // hostname of the agent
	AltHost            string    // hostname of the agent (alternative) (optional)
	Port               uint      // port of the agent
	IsServer           bool      // is this node a server or not
	IsLeader           bool      // is this node leader or not
	ServersAddr        string    // address of the servers in the cluster
	PrimaryServersAddr string    // address of the primary servers in the cluster
	Services           []Service // services in the cluster
}

var (
	Verbose bool // verbose mode
	WithAPI bool // with api endpoint
)

func LoadConfiguration() *Config {

	path := "/etc/centor.d/"

	// load configuration from file
	configs, err := loadConfigsFromDir(path)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	// compile hcl configuration to struct
	cnf, err := compile(configs)
	if err != nil {
		log.Fatalf("Failed to compile configuration: %s", err)
	}

	// get environment variables
	if e := os.Getenv("PORT"); e != "" {
		p, _ := strconv.Atoi(e)
		cnf.Port = uint(p)
	}
	if e := os.Getenv("NAME"); e != "" {
		cnf.Name = e
	}
	if e := os.Getenv("HOST"); e != "" {
		cnf.Host = e
	}
	if e := os.Getenv("JOIN"); e != "" {
		cnf.ServersAddr = e
	}
	if e := os.Getenv("PRIMARIES"); e != "" {
		cnf.PrimaryServersAddr = e
	}
	if e := os.Getenv("SERVER"); e != "" {
		cnf.IsServer = true
	}
	if e := os.Getenv("LEADER"); e != "" {
		cnf.IsLeader = true
	}
	if e := os.Getenv("ALTERNATIVE_HOST"); e != "" {
		cnf.AltHost = e
	}

	// load config from cli arguments
	flag.BoolVar(&Verbose, "v", false, "")
	flag.BoolVar(&WithAPI, "api", false, "")
	flag.StringVar(&cnf.Name, "n", cnf.Name, "")
	flag.StringVar(&cnf.Host, "h", cnf.Host, "")
	flag.StringVar(&cnf.AltHost, "ah", cnf.AltHost, "")
	flag.UintVar(&cnf.Port, "p", cnf.Port, "")
	flag.StringVar(&cnf.ServersAddr, "join", cnf.ServersAddr, "")
	flag.StringVar(&cnf.PrimaryServersAddr, "primaries-addr", cnf.PrimaryServersAddr, "")
	flag.BoolVar(&cnf.IsServer, "server", cnf.IsServer, "")
	flag.BoolVar(&cnf.IsLeader, "leader", cnf.IsLeader, "")
	flag.Parse()

	// print configuration an verbose mode
	if Verbose {
		cb, _ := json.MarshalIndent(cnf, "", " ")
		fmt.Printf("%s\n", cb)
	}

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
