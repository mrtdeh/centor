package envoy

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// NewEnvoyConfig creates a new EnvoyConfig instance
// with the given configuration settings and TLS configuration
//
// The configuration settings :
//
//	envoy.EnvoyConfig{
//		ListenerPort: 8081,
//		EndpointPort: 3000,
//		TLSConfig: envoy.TLSConfig{
//			Secure:         true,
//			CA:             "./pkg/envoy/testData/certs/ca.crt",
//			Cert:           "./pkg/envoy/testData/certs/server.crt",
//			Key:            "./pkg/envoy/testData/certs/server.key",
//			SessionTimeout: "6000s",
//		},
//	}
func NewEnvoy(cnf EnvoyConfig) error {
	strConfig := DefaultConfig

	if cnf.LogPath == "" {
		cnf.LogPath = "/var/log/centor.log"
	}
	if cnf.ListenerAddress == "" {
		cnf.ListenerAddress = "0.0.0.0"
	}
	if cnf.ListenerPort == 0 {
		cnf.ListenerPort = 80
	}
	if cnf.TLSConfig.Secure {
		strConfig = fmt.Sprintf(strConfig, downstreamTLS)
	} else {
		strConfig = fmt.Sprintf(strConfig, "")
	}

	if cnf.EndpointAddress == "" {
		cnf.EndpointAddress = "127.0.0.1"
	}
	if cnf.EndpointPort == 0 {
		return fmt.Errorf("endpoint port not specified")
	}

	strConfig = strings.ReplaceAll(strConfig, "{listener_address}", cnf.ListenerAddress)
	strConfig = strings.ReplaceAll(strConfig, "{listener_port}", fmt.Sprintf("%d", cnf.ListenerPort))
	strConfig = strings.ReplaceAll(strConfig, "{log_path}", cnf.LogPath)
	strConfig = strings.ReplaceAll(strConfig, "{session_timeout}", cnf.TLSConfig.SessionTimeout)
	if cnf.TLSConfig.DisableSessionTicket {
		strConfig = strings.ReplaceAll(strConfig, "{disable_session_ticket}", "true")
	} else {
		strConfig = strings.ReplaceAll(strConfig, "{disable_session_ticket}", "false")
	}
	strConfig = strings.ReplaceAll(strConfig, "{ssl_cert}", cnf.TLSConfig.Cert)
	strConfig = strings.ReplaceAll(strConfig, "{ssl_key}", cnf.TLSConfig.Key)
	strConfig = strings.ReplaceAll(strConfig, "{ssl_ca}", cnf.TLSConfig.CA)
	strConfig = strings.ReplaceAll(strConfig, "{endpoint_address}", cnf.EndpointAddress)
	strConfig = strings.ReplaceAll(strConfig, "{endpoint_port}", fmt.Sprintf("%d", cnf.EndpointPort))

	cmd := exec.Command("envoy", "--config-yaml", strConfig)
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err.Error())
	}

	return nil
}
