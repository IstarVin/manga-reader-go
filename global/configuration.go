package global

type Configuration struct {
	DataDir string `json:"dataDir"`

	Host string `json:"host"`
	Port int    `json:"port"`
}

var DefaultConfig = &Configuration{
	DataDir: ".",
	Host:    "0.0.0.0",
	Port:    8080,
}

var Config *Configuration
