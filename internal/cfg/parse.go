package cfg

import (
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/web-programming-fall-2022/digivision-backend/internal/cfg_utils"
)

const devConfigFileName = "config.dev.yml"

func ParseDevConfig() Config {
	_, currentFilePath, _, _ := runtime.Caller(0)
	return ParseConfig(path.Join(path.Dir(currentFilePath), devConfigFileName))
}

// ParseConfig parses config file specified by path into config struct and validates it.
func ParseConfig(path string) Config {
	config := Config{}
	err := cfg_utils.UnmarshalConfig("DVS", path, &config)
	if err != nil {
		logrus.Panic(err.Error())
	}
	return config
}
