package configparser

import (
	"os"
	
	"gopkg.in/yaml.v3"
)

// Note to self: Based on the layout of config.yaml
type Proxy struct {
	ProxyPort    string                  `yaml:"ProxyPort"`
	NumWorkers   int                     `yaml:"NumWorkers"`
	MaxConcReq   int                     `yaml:"MaxConcReq"`
	Servers      map[string]ServerConfig `yaml:"Servers"`
	ServerConfig []ServerConfig          `yaml:"ServerConfig"`
}

type ServerConfig struct {
	ServerPorts   []string `yaml:"ServerPorts"`
	Names         []string `yaml:"Names"`
	ServerRouting Routing  `yaml:"Routing"`
}

type Routing struct {
	Prefix       string      `yaml:"Prefix"`
	ProxyPass    []string    `yaml:"ProxyPass"`
	ProxyHeaders ProxyHeader `yaml:"ProxyHeaders"`
}

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
