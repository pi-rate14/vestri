package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	config "github.com/pi-rate14/simple-lb/pkg/config"
	logger "github.com/sirupsen/logrus"
)

var (
	port = flag.Int("port", 8080, "THe port where the load balancer starts")
)

type Vestri struct {
	Config     *config.Config
	ServerList *config.ServerList
}

func NewVestri(configuration *config.Config) *Vestri {

	servers := make([]*config.Server, 0)
	for _, service := range configuration.Services {
		for _, replica := range service.Replicas {
			replicaURL, err := url.Parse(replica)
			if err != nil {
				logger.Fatal(err)
			}

			proxy := httputil.NewSingleHostReverseProxy(replicaURL)
			// newServer := config.NewServer(replicaURL, proxy)

			servers = append(servers, &config.Server{
				Url:   replicaURL,
				Proxy: proxy,
			})
		}
	}

	return &Vestri{
		Config: configuration,
		ServerList: &config.ServerList{
			Servers: servers,
			Current: uint32(0),
		},
	}
}

func (vestri *Vestri) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// for each service, read request path host:port/remaining/url
	// and load balance it to service named "service" with url host[i]:port[i]/remaining/url
	logger.Infof("Received New Request: url='%s'", r.Host)
	nextServer := vestri.ServerList.Next()
	// forward request to proxy
	logger.Infof("Forwarding to the server='%s'", vestri.ServerList.Servers[nextServer].Url.String())
	vestri.ServerList.Servers[nextServer].Proxy.ServeHTTP(w, r)
	// vestri.ServerList.Servers[nextServer].Forward(w, r)
}

func main() {
	flag.Parse()

	config := &config.Config{
		Services: []config.Service{
			{
				Name: "Test",
				Replicas: []string{
					"http://localhost:8081",
					"http://localhost:8082",
				},
			},
		},
	}

	vestri := NewVestri(config)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: vestri,
	}

	err := server.ListenAndServe()
	if err != nil {
		logger.Fatal()
	}
}
