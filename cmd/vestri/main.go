package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	config "github.com/pi-rate14/simple-lb/pkg/config"
	"github.com/pi-rate14/simple-lb/pkg/domain"
	"github.com/pi-rate14/simple-lb/pkg/strategy"
	logger "github.com/sirupsen/logrus"
)

var (
	port       = flag.Int("port", 8080, "THe port where the load balancer starts")
	configFile = flag.String("config", "", "Path to the configuration file ")
)

type Vestri struct {
	Config            *config.Config                // conguration of the app loaded from yaml file
	ServiceServersMap map[string]*config.ServerList // ServiceServersMap contains map between service matcher and its replicas
}

func NewVestri(configuration *config.Config) *Vestri {
	serverMap := make(map[string]*config.ServerList, 0)

	for _, service := range configuration.Services {

		servers := make([]*domain.Server, 0)

		for _, replica := range service.Replicas {

			replicaURL, err := url.Parse(replica.Url)
			if err != nil {
				logger.Fatal(err)
			}

			proxy := httputil.NewSingleHostReverseProxy(replicaURL)
			servers = append(servers, &domain.Server{
				Url:   replicaURL,
				Proxy: proxy,
			})
		}

		serverMap[service.Matcher] = &config.ServerList{
			Servers:  servers,
			Name:     service.Name,
			Strategy: strategy.LoadStrategy(service.Strategy),
		}
	}

	return &Vestri{
		Config:            configuration,
		ServiceServersMap: serverMap,
	}
}

func (vestri *Vestri) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// for each service, read request path host:port/remaining/url
	// and load balance it to service named "service" with url host[i]:port[i]/remaining/url
	logger.Infof("Received New Request: url='%s'", r.Host)

	serverList, err := vestri.findService(r.URL.Path)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	nextServer, err := serverList.Strategy.Next(serverList.Servers)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// forward request to proxy
	logger.Infof("Forwarding to the server='%s'", nextServer.Url.String())

	nextServer.Forward(w, r)
}

// findService looks for the first serverList that matches the request path (matcher)
// returns an error if no matcher found
func (vestri *Vestri) findService(path string) (*config.ServerList, error) {
	logger.Infof("Trying to find matcher for request: %s", path)

	serverList, ok := vestri.ServiceServersMap[path]
	if !ok {
		return nil, fmt.Errorf("could not find a matcher for the request '%s'", path)
	}

	logger.Infof("Found service '%s' matching the request", serverList.Name)
	return serverList, nil

}

func main() {
	flag.Parse()

	file, err := os.Open(*configFile)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()
	config, err := config.LoadConfig(file)
	if err != nil {
		logger.Fatal(err)
	}

	vestri := NewVestri(config)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: vestri,
	}

	err = server.ListenAndServe()
	if err != nil {
		logger.Fatal()
	}
}
