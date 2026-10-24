package configparser

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Although the config and stateTracker both deal with the servers,
// config only loads the data about the servers in memory
// While the stateTracker deals with the information
// tracking about the Connections and Servers at the runtime.

// Note to self: Based on the layout of config.yaml

// Proxy holds the proxy's runtime configuration and the configuration
// of all registered backend server groups.
//
// ProxyPort specifies the address on which the proxy listens.
// NumWorkers specifies the number of workers used to process requests.
// MaxConcReq specifies the configured maximum number of concurrent requests.
// Servers maps server identifiers to their respective configurations.
// ServerConfig contains the server configurations as a slice for convenient iteration.
type Proxy struct {
	ProxyPort    string                  `yaml:"ProxyPort"`
	NumWorkers   int                     `yaml:"NumWorkers"`
	MaxConcReq   int                     `yaml:"MaxConcReq"`
	Servers      map[string]ServerConfig `yaml:"Servers"`
	ServerConfig []ServerConfig          `yaml:"ServerConfig"`
}

// ServerConfig defines a backend server group and its routing configuration.
//
// ServerPorts contains the ports associated with the server group.
// Names contains the names associated with the server group.
// ServerRouting specifies how incoming requests are routed to the backends.
type ServerConfig struct {
	ServerPorts   []string `yaml:"ServerPorts"`
	Names         []string `yaml:"Names"`
	ServerRouting Routing  `yaml:"Routing"`
}

// Routing defines the rules for matching incoming requests and forwarding
// them to backend servers.
//
// Prefix specifies the URL path prefix associated with the routing rule.
// ProxyPass contains the list of backend destinations to which requests
// can be forwarded.
// ProxyHeaders specifies the HTTP header values configured for forwarded requests.
type Routing struct {
	Prefix       string      `yaml:"Prefix"`
	ProxyPass    []string    `yaml:"ProxyPass"`
	ProxyHeaders ProxyHeader `yaml:"ProxyHeaders"`
}

// ProxyHeader holds the configured values for HTTP headers used when
// forwarding requests to backend servers.
//
// Host specifies the value of the Host header.
// XRealIP specifies the value of the X-Real-IP header.
// XForwardedFor specifies the value of the X-Forwarded-For header.
type ProxyHeader struct {
	Host          string `yaml:"Host"`
	XRealIP       string `yaml:"XRealIP"`
	XForwardedFor string `yaml:"XForwardedFor"`
}

type configFile struct {
	Proxy struct {
		ProxyPort  string `yaml:"ProxyPort"`
		NumWorkers int    `yaml:"NumWorkers"`
		MaxConcReq int    `yaml:"MaxConcReq"`
	} `yaml:"Proxy"`
	Servers map[string]ServerConfig `yaml:"Servers"`
}

// Copies the server configurations from the Servers map into the ServerConfig slice.
// The order of iteration is not specefic since it is stored in map

func LoadConfig() (Proxy, error) {
	configPath := "config/config.yaml"
	data, err := os.ReadFile(configPath)

	if err != nil {
		return defaultProxy(), err
	}

	var rawConfig configFile
	if err := yaml.Unmarshal(data, &rawConfig); err != nil {
		return defaultProxy(), err
	}

	proxy := Proxy{
		ProxyPort:  rawConfig.Proxy.ProxyPort,
		NumWorkers: rawConfig.Proxy.NumWorkers,
		MaxConcReq: rawConfig.Proxy.MaxConcReq,
		Servers:    rawConfig.Servers,
	}

	for _, sc := range rawConfig.Servers {
		proxy.ServerConfig = append(proxy.ServerConfig, sc)
	}

	return proxy, nil
}

func defaultProxy() Proxy {
	return Proxy{
		ProxyPort:    ":8080",
		NumWorkers:   10,
		MaxConcReq:   1000,
		Servers:      make(map[string]ServerConfig),
		ServerConfig: []ServerConfig{},
	}
}
