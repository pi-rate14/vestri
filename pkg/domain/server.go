package domain

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
)

// TODO: make metadata own object
type Replica struct {
	Url      string            `yaml:"url"`
	Metadata map[string]string `yaml:"metadata"`
}

type Service struct {
	Name     string    `yaml:"name"`     // Name of the service (test-service-1)
	Matcher  string    `yaml:"matcher"`  // prefix matcher to select service based on URL path (/api/v1)
	Strategy string    `yaml:"strategy"` // LB strategy for this service
	Replicas []Replica `yaml:"replicas"` // URLs of the replicas of this service (8081, 8082)
}

// Server represents an instance of a running server
type Server struct {
	Url      *url.URL               // URL of the server instance
	Proxy    *httputil.ReverseProxy // Proxy responsible for this server
	Metadata map[string]string      // connection count of server and other metadata
	rwMutex  sync.Mutex             // Mutex to update health of server
	alive    bool                   // health status of server
}

func (server *Server) Forward(w http.ResponseWriter, r *http.Request) {
	server.Proxy.ServeHTTP(w, r)
}

// returns the string value associated with the given key in server metadata or returns def
func (server *Server) GetMetaOrDefault(key, def string) string {
	value, ok := server.Metadata[key]
	if !ok {
		return def
	}

	return value
}

// returns the int value associated with the given key in server metadata or returns def
func (server *Server) GetMetaOrDefaultInt(key string, def int) int {
	value := server.GetMetaOrDefault(key, fmt.Sprintf("%d", def))

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return def // default return is 1 because primary use case of function is WRR and defaults to 1 for RR
	}

	return intValue
}

// Changes the currnt alive field value and returns the old value
func (server *Server) SetLiveness(value bool) bool {
	server.rwMutex.Lock()
	defer server.rwMutex.Unlock()

	oldVal := server.alive
	server.alive = value

	return oldVal
}

// returns the health status of the server
func (server *Server) IsAlive() bool {
	server.rwMutex.Lock()
	defer server.rwMutex.Unlock()

	return server.alive
}
