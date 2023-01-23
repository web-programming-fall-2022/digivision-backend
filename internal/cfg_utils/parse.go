package cfg_utils

import (
	"bytes"
	"fmt"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"gopkg.in/yaml.v2"
	"strings"

	"github.com/spf13/viper"
)

// UnmarshalConfig parses config file specified by path into config struct and validates it.
// config should always be a pointer (reference).
func UnmarshalConfig(prefix, path string, config Config) error {
	viperInstance, err := newViper(prefix, path, config)
	if err != nil {
		return err
	}
	if err := viperInstance.UnmarshalExact(config); err != nil {
		return fmt.Errorf("could not unmarshal configuration file: %s", err)
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %s", err)
	}
	return nil
}

func newViper(prefix, configPath string, config Config) (*viper.Viper, error) {
	v := viper.New()

	v.SetDefault("prometheus.buckets", []float64{0.05, 0.1, 0.2, 0.5, 1, 2, 5})
	v.SetDefault("prometheus.prefix", "metrics")
	v.SetDefault("milvus.vectorDim", 768)
	v.SetDefault("milvus.metricType", entity.L2)
	v.SetDefault("milvus.nProbe", 16)
	v.SetDefault("milvus.collectionName", "products_revis_digikala_clip_ViT_L_14_336px")

	v.SetConfigType("yaml")
	v.SetEnvPrefix(prefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()
	if err := addEmptyConfig(v, config); err != nil {
		return nil, err
	}

	if configPath == "" {
		return v, nil
	}
	v.SetConfigFile(configPath)
	if err := v.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("could not read config from env or config file: %v", err)
	}
	return v, nil
}

func addEmptyConfig(v *viper.Viper, config Config) error {
	// We set default empty values in viper using the config struct so that it knows which env
	// variables to bind. Otherwise it will ignore any env vars not provided in config.
	// https://github.com/spf13/viper/issues/188
	b, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("could not marshal config struct to yaml: %v", err)
	}
	if err := v.ReadConfig(bytes.NewReader(b)); err != nil {
		return fmt.Errorf("could not configure viper with empty config: %v", err)
	}
	return nil
}
