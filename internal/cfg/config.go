package cfg

import (
	"github.com/arimanius/digivision-backend/internal/bootstrap"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	Env string
	Log struct {
		Level string
	}

	bootstrap.GrpcServerRunnerConfig `mapstructure:",squash" yaml:",inline"`

	HttpServer struct {
		Port int
	}
}

func (c *Config) Validate() error {
	return validation.Errors{
		"env": validation.Validate(c.Env, validation.Required),
		"log.level": validation.Validate(c.Log.Level, validation.Required, validation.In(
			"panic", "fatal", "error", "warn", "info", "debug", "trace",
		)),
		"server.host": validation.Validate(c.Server.Host, validation.Required),
		"server.port": validation.Validate(c.Server.Port, validation.Required),
	}.Filter()
}
