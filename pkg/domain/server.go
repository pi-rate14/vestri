package domain

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Service struct {
	Name     string   `yaml:"name"`     // Name of the service (test-service-1)
	Matcher  string   `yaml:"matcher"`  // prefix matcher to select service based on URL path (/api/v1)
	Strategy string   `yaml:"strategy"` // LB strategy for this service
	Replicas []string `yaml:"replicas"` // URLs of the replicas of this service (8081, 8082)
}

// Server represents an instance of a running server
type Server struct {
	Url   *url.URL               // URL of the server instance
	Proxy *httputil.ReverseProxy // Proxy responsible for this server

}

func (s *Server) Forward(w http.ResponseWriter, r *http.Request) {
	s.Proxy.ServeHTTP(w, r)
}
