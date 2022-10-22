package config

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func LoadConfig(reader io.Reader) (*Config, error) {
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	conf := Config{}
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
