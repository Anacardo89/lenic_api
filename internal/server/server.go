package server

type Config struct {
	Host     string `yaml:"host"`
	GrpcPort string `yaml:"grpcPort"`
}

var (
	Server *Config
)
