package configparser

type Proxy struct {
	ProxyPort string
	NumWorkers int
	serverconfig []ServerConfig
}

type ServerConfig struct {
	ServerPorts []string
	Names []string
	ServerRouting Routing
}

type Routing struct {
	Prefix string
	ProxyPass []string
	ProxyHeaders ProxyHeader
}

type ProxyHeader struct {
	Host string
	XRealIP string
	XForwardedFor string
}

func LoadConfig() Proxy {
	...
}