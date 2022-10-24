package strategy

// TODO: rename to strategy.go

import (
	"fmt"
	"sync"

	"github.com/pi-rate14/simple-lb/pkg/domain"
	logger "github.com/sirupsen/logrus"
)

// Implemented LB strategies
const (
	ROUND_ROBIN          = "RoundRobin"
	WEIGHTED_ROUND_ROBIN = "WeightedRoundRobin"
	UNKNOWN              = "Unknown"
)

// BalancingStrategy is the Load Balancing interface that every LB strategy must implement
type BalancingStrategy interface {
	Next([]*domain.Server) (*domain.Server, error)
}

// generate BalancingStrategy for each call
var strategies map[string]func() BalancingStrategy

func init() {
	strategies = make(map[string]func() BalancingStrategy, 0)
	strategies[ROUND_ROBIN] = func() BalancingStrategy {
		return &RoundRobin{
			mutex:   sync.Mutex{},
			current: uint32(0),
		}
	}

	strategies[WEIGHTED_ROUND_ROBIN] = func() BalancingStrategy {
		return &WeightedRoundRobin{
			mutex: sync.Mutex{},
		}
	}
}

type RoundRobin struct {
	current uint32     // current server to forward the request to.
	mutex   sync.Mutex // lock to grant exlusive access to struct variables
}

// TODO : Check if only returning uint32 is possble
func (roundRobin *RoundRobin) Next(servers []*domain.Server) (*domain.Server, error) {
	roundRobin.mutex.Lock()
	defer roundRobin.mutex.Unlock()
	seen := 0
	var nextServer *domain.Server
	for seen < len(servers) {
		nextServer = servers[roundRobin.current]
		roundRobin.current = (roundRobin.current + 1) % uint32(len(servers))
		if nextServer.IsAlive() {
			break
		}
		seen += 1
	}
	if nextServer == nil || seen == len(servers) {
		logger.Error("All servers are down")
		return nil, fmt.Errorf("all %d servers are unavailable", seen)
	}
	// lenServerList := uint32(len(servers))
	// nextServer := servers[next%lenServerList]
	logger.Infof("Strategy picked server : %s", nextServer.Url.Host)
	return nextServer, nil
}

type WeightedRoundRobin struct {
	mutex   sync.Mutex // lock to grant exlusive access to struct variables
	count   []int      // keep track of the number of request server i processed
	current int        // index of the last server that executed the request
}

func (wrr *WeightedRoundRobin) Next(servers []*domain.Server) (*domain.Server, error) {
	wrr.mutex.Lock()
	defer wrr.mutex.Unlock()

	if wrr.count == nil {
		wrr.count = make([]int, len(servers))
		wrr.current = 0
	}

	seen := 0
	var nextServer *domain.Server

	for seen < len(servers) {

		nextServer = servers[wrr.current]
		// find capacity of the current server ()
		capacity := nextServer.GetMetaOrDefaultInt("weight", 1)
		if !nextServer.IsAlive() {
			seen += 1
			wrr.count[wrr.current] = 0
			wrr.current = (wrr.current + 1) % len(servers)
			continue
		}

		if wrr.count[wrr.current] < capacity {
			wrr.count[wrr.current] += 1
			logger.Infof("Strategy picked server '%s'", nextServer.Url.Host)
			return nextServer, nil
		}

		wrr.count[wrr.current] = 0
		wrr.current = (wrr.current + 1) % len(servers)
	}

	if nextServer == nil || seen == len(servers) {
		logger.Error("All servers are down")
		return nil, fmt.Errorf("all %d servers are unavailable", seen)
	}

	return nextServer, nil
}

// Set the LB strategy based on name. Defualts to Round Robin
func LoadStrategy(name string) BalancingStrategy {
	strategy, ok := strategies[name]
	if !ok {
		logger.Warnf("Strategy with name %s not found. Falling back to Round Robin.", name)
		return strategies[ROUND_ROBIN]()
	}
	logger.Infof("Picked strategy: %s", name)
	return strategy()
}
