package strategy

// TODO: rename to strategy.go

import (
	"sync/atomic"

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
			Current: uint32(0),
		}
	}
}

type RoundRobin struct {
	Current uint32 // current server to forward the request to.
}

// TODO : Check if only returning uint32 is possble
func (roundRobin *RoundRobin) Next(servers []*domain.Server) (*domain.Server, error) {

	next := atomic.AddUint32(&roundRobin.Current, uint32(1))

	lenServerList := uint32(len(servers))
	nextServer := servers[next%lenServerList]
	logger.Infof("Strategy picked server : %s", nextServer.Url.String())
	return nextServer, nil
}

// Set the LB strategy based on name. Defualts to Round Robin
func LoadStrategy(name string) BalancingStrategy {
	strategy, ok := strategies[name]
	if !ok {
		return strategies[ROUND_ROBIN]()
	}

	return strategy()
}
