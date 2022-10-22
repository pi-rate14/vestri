package config

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type Service struct {
	// Name of the service
	Name string `yaml:"name"`
	// URLs of the replicas of this service
	Replicas []string `yaml:"replicas"`
}

/*
	 Config is the configuration given to lb from a config source

	 service log with 4 replicas
     log: ip1:port1 ip2:port2 ...
*/
type Config struct {
	// URLs of the services
	Services []Service `yaml:"services"`
	// Name of LB strategy
	Strategy string `yaml:"strategy"`
}

// Server represents an instance of a running server
type Server struct {
	// URL of the server instance
	Url *url.URL
	// Proxy responsible for this server
	Proxy *httputil.ReverseProxy
}

func (s *Server) Forward(w http.ResponseWriter, r *http.Request) {
	s.Proxy.ServeHTTP(w, r)
}

type ServerList struct {
	// List of all the servers
	Servers []*Server
	// current server to forward the request to.
	// Next server should be (cur + 1) % len(Servers)
	Current uint32
}

func (serverList *ServerList) Next() uint32 {
	next := atomic.AddUint32(&serverList.Current, uint32(1))
	lenServerList := uint32(len(serverList.Servers))
	// if next >= lenServerList {
	// 	next -= lenServerList
	// }
	return next % lenServerList
}
