package cfg

import (
	"github.com/arimanius/digivision-backend/internal/bootstrap"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
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

	Img2Vec struct {
		Addr string
	}

	Milvus struct {
		Addr           string
		VectorDim      int
		MetricType     entity.MetricType
		NProbe         int
		CollectionName string
	}

	ObjectDetector struct {
		Addr string
	}

	Redis struct {
		Addr string
	}
}

func (c *Config) Validate() error {
	return validation.Errors{
		"env": validation.Validate(c.Env, validation.Required),
		"log.level": validation.Validate(c.Log.Level, validation.Required, validation.In(
			"panic", "fatal", "error", "warn", "info", "debug", "trace",
		)),
		"server.host":           validation.Validate(c.Server.Host, validation.Required),
		"server.port":           validation.Validate(c.Server.Port, validation.Required),
		"milvus.vectorDim":      validation.Validate(c.Milvus.VectorDim, validation.Required),
		"milvus.metricType":     validation.Validate(c.Milvus.MetricType, validation.Required),
		"milvus.nProbe":         validation.Validate(c.Milvus.NProbe, validation.Required),
		"milvus.collectionName": validation.Validate(c.Milvus.CollectionName, validation.Required),
	}.Filter()
}
