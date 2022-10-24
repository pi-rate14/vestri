package health

import (
	"errors"
	"net"
	"time"

	"github.com/pi-rate14/simple-lb/pkg/domain"
	logger "github.com/sirupsen/logrus"
)

type HealthChecker struct {
	servers []*domain.Server
	period  int
}

// Create a new instance of a health checker
func NewHealthChecker(servers []*domain.Server) (*HealthChecker, error) {
	if len(servers) == 0 {
		return nil, errors.New("server list is empty. cannot check health")
	}

	return &HealthChecker{
		servers: servers,
		period:  1,
	}, nil
}

// loops infinietely to check health of every server
func (checker *HealthChecker) Start() {
	logger.Info("Starting health checker...")
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, server := range checker.servers {
				go checkHealth(server)
			}
		}
	}
}

// toggles liveness of server
func checkHealth(server *domain.Server) {
	// server is healthy if a tcp connection can be established
	// within a reasonable time frame
	conn, err := net.DialTimeout("tcp", server.Url.Host, time.Second*5)

	if err != nil {
		logger.Errorf("Could not connect to the server at '%s'", server.Url.Host)
		oldStatus := server.SetLiveness(false)
		if oldStatus {
			logger.Warnf("Transitioning server '%s' from Live to Unavailable", server.Url.Host)
		}
		return
	}

	defer conn.Close()

	oldStatus := server.SetLiveness(true)
	if !oldStatus {
		logger.Infof("Transitioning server '%s' from Unavailable to Live", server.Url.Host)
	}
}
