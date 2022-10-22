package config

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type Service struct {
	Name     string   `yaml:"name"`     // Name of the service
	Matcher  string   `yaml:"matcher"`  // prefix matcher to select service based on URL path
	Replicas []string `yaml:"replicas"` // URLs of the replicas of this service
}

/*
	 Config is the configuration given to lb from a config source

	 service log with 4 replicas
     log: ip1:port1 ip2:port2 ...
*/
type Config struct {
	Services []Service `yaml:"services"` // URLs of the services
	Strategy string    `yaml:"strategy"` // Name of LB strategy
}

// Server represents an instance of a running server
type Server struct {
	Url   *url.URL               // URL of the server instance
	Proxy *httputil.ReverseProxy // Proxy responsible for this server

}

func (s *Server) Forward(w http.ResponseWriter, r *http.Request) {
	s.Proxy.ServeHTTP(w, r)
}

type ServerList struct {
	Servers []*Server // List of all the servers
	Current uint32    // current server to forward the request to.

}

func (serverList *ServerList) Next() uint32 {
	next := atomic.AddUint32(&serverList.Current, uint32(1))
	lenServerList := uint32(len(serverList.Servers))
	// if next >= lenServerList {
	// 	next -= lenServerList
	// }
	return next % lenServerList
}
