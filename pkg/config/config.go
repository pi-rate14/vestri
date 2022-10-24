package config

import (
	"github.com/pi-rate14/simple-lb/pkg/domain"
	"github.com/pi-rate14/simple-lb/pkg/health"
	"github.com/pi-rate14/simple-lb/pkg/strategy"
)

/*
	 Config is the configuration given to lb from a config source

	 service log with 4 replicas
     log: ip1:port1 ip2:port2 ...
*/
type Config struct {
	Services []domain.Service `yaml:"services"` // URLs of the services
	Strategy string           `yaml:"strategy"` // Name of LB strategy
}

type ServerList struct {
	Name     string                     // Name of ther service that has this serverList
	Servers  []*domain.Server           // List of all the servers
	Strategy strategy.BalancingStrategy // how this serverlist will be load balanced. Defaults to round robin
	Checker  *health.HealthChecker
}
