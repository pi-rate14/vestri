package main

import (
	"flag"
	"fmt"
	"net/http"

	logger "github.com/sirupsen/logrus"
)

var port = flag.Int("port", 8081, "Port to start demo service on")

type DemoServer struct {
}

func (demoServer *DemoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.Infof("Recieved Forwarded Request from: %v", r.URL.Host)
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Received Response from port %d\n", *port)))
}

func main() {
	flag.Parse()
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), &DemoServer{})
	if err != nil {
		logger.Fatal(err)
	}
}
