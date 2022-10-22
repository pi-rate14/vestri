package config

import (
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	conf, err := LoadConfig(strings.NewReader(`
services:
  - name: "test service"
    replicas:
    - "localhost:8081"
    - "localhost:8082"
strategy: "RoundRobin"`))

	if err != nil {
		t.Errorf("Error should be nil %v", err)
	}

	if conf.Strategy != "RoundRobin" {
		t.Errorf("Strategy undefined. Expected RoundRobin, Got : %s", conf.Strategy)
	}

	if len(conf.Services) != 1 {
		t.Errorf("Expected service count to be 1, Got : %d", len(conf.Services))
	}

	if conf.Services[0].Name != "test service" {
		t.Errorf("Expected service name to be 'test service', Got : %s", conf.Services[0].Name)
	}

	if len(conf.Services[0].Replicas) != 2 {
		t.Errorf("Expected service replica count to be 3, Got : %d", len(conf.Services[0].Replicas))
	}

	if conf.Services[0].Replicas[0] != "localhost:8081" {
		t.Errorf("Expected first service replica to be 'localhost:8081', Got : %s", conf.Services[0].Replicas[0])
	}

	if conf.Services[0].Replicas[1] != "localhost:8082" {
		t.Errorf("Expected second service replica to be 'localhost:8082', Got : %s", conf.Services[0].Replicas[1])
	}
}
