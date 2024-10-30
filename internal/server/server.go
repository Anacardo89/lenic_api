package server

type Config struct {
	Host      string `yaml:"host"`
	ProxyPORT string `yaml:"proxyPort"`
	HttpPORT  string `yaml:"httpPort"`
	HttpsPORT string `yaml:"httpsPort"`
}

var (
	Server *Config
)
