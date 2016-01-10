package main

import (
	"log"
	"net/http"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/josdirksen/slackproxy/config"
	"github.com/josdirksen/slackproxy/handlers"
	"os"
	"flag"
)

// setup the listener
func MainListener(w http.ResponseWriter, req *http.Request) {
	// get the contents from the request body, convert it to a command
	req.ParseForm()
	cmdToExecute := handlers.ParseInput(req.Form)
	client := setupDockerClient(cmdToExecute.Environment)

	// execute the command
	handlers.HandleCommand(cmdToExecute, client, w)
}

func setupDockerClient(env string) *docker.Client {

	// first get the environment from the config
	cfg, err := config.GetEnvironmentConfig(env)

	if (err != nil) {
		println(err)
		log.Fatal("Can't parse config, exiting")
		os.Exit(1)
	}

	if (cfg.Tls) {
		endpoint := cfg.Host
		path := cfg.Path
		ca := fmt.Sprintf("%s/%s", path, cfg.Ca)
		cert := fmt.Sprintf("%s/%s", path, cfg.Cert)
		key := fmt.Sprintf("%s/%s", path, cfg.Key)

		client,_ := docker.NewTLSClient(endpoint, cert, key, ca)
		return client
	} else {
		client,_ := docker.NewClient(cfg.Host)
		return client
	}

}

func StartListening() {
	http.HandleFunc("/glitchrequest", MainListener)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	var configLocation = flag.String("config", "config.json", "specify the config file")
	flag.Parse()
	// first parse the config
	config.ParseConfig(*configLocation)
	// setup the handler that listens to 9000
	StartListening()
}