package main

import (
	"io/ioutil"
	"os"

	"sigs.k8s.io/yaml"

	"github.com/opensourceways/xihe-training-center/domain"
	"github.com/opensourceways/xihe-training-center/huaweicloud/training"
)

type Config struct {
	Training domain.TrainingConfig `json:"training"`
	Cloud    training.HuaweiCloud  `json:"cloud"`
}

func (cfg *Config) setDefault() {
	cfg.Training.Setdefault()
}

func (cfg *Config) validate() error {
	return nil
}

func loadConfig(path string) (*Config, error) {
	v := new(Config)

	if err := loadFromYaml(path, v); err != nil {
		return nil, err
	}

	v.setDefault()

	if err := v.validate(); err != nil {
		return nil, err
	}

	return v, nil
}

func loadFromYaml(path string, cfg interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	content := []byte(os.ExpandEnv(string(b)))

	return yaml.Unmarshal(content, cfg)
}
