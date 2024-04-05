package configs

type HTTPMode string

const (
	HTTPModeDevelopment HTTPMode = "development"
	HTTPModeProduction  HTTPMode = "production"
)

type HTTP struct {
	Address string   `yaml:"address"`
	Mode    HTTPMode `yaml:"mode"`
}
