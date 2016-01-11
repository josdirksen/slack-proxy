package main

import (
	"flag"
	"github.com/josdirksen/slackproxy/config"
	"github.com/josdirksen/slackproxy/handlers"
	"log"
	"net/http"
)

// setup the listener, with a config passed in.
func GetConfigListener(config *config.Configuration) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// get the contents from the request body, convert it to a command
		req.ParseForm()
		cmdToExecute := handlers.ParseInput(req.Form)

		// execute the command
		handler := handlers.GetHandler("docker", config)
		handler.HandleCommand(cmdToExecute, w)
	}
}

func StartListening(config *config.Configuration) {
	http.HandleFunc("/handleSlackCommand", GetConfigListener(config))
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
	StartListening(config.GetConfig())
}
