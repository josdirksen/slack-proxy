package handlers

import (
	"io"
	"net/http"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"bytes"
	"time"
	"github.com/josdirksen/slackproxy/config"
	"log"
	"os"
)

type DockerHandler struct {
	config	*config.Configuration
}

func NewDockerHandler(configuration *config.Configuration) *DockerHandler {
	p := new(DockerHandler)
	p.config = configuration
	return p
}

func (dh DockerHandler) HandleCommand(cmdToExecute *Command, w http.ResponseWriter) {
	client := setupDockerClient(cmdToExecute.Environment)

	switch cmdToExecute.SlackCommand {
	case "ps" : handlePsCommand(client, w)
	case "imgs" : handleListImagesCommand(client, w)
	}
}

func handleListImagesCommand(client *docker.Client, w http.ResponseWriter) {
	images, _ := client.ListImages(docker.ListImagesOptions{All: false})
	for _, img := range images {
		fmt.Println("ID: ", img.ID)
		fmt.Println("RepoTags: ", img.RepoTags)
		fmt.Println("Created: ", img.Created)
		fmt.Println("Size: ", img.Size)
		fmt.Println("VirtualSize: ", img.VirtualSize)
		fmt.Println("ParentId: ", img.ParentID)
	}
}

func handlePsCommand(client *docker.Client, w http.ResponseWriter) {
	containers, _ := client.ListContainers(docker.ListContainersOptions{All: false})
	var buffer bytes.Buffer
	for _, container := range containers {
		buffer.WriteString(fmt.Sprintf("ID: %s\n", container.ID))
		buffer.WriteString(fmt.Sprintf("Command: %s\n", container.Command))
		buffer.WriteString(fmt.Sprintf("Created: %s\n", time.Unix(container.Created, 0)))
		buffer.WriteString(fmt.Sprintf("Image: %s\n", container.Image))
		buffer.WriteString(fmt.Sprintf("Status: %s\n", container.Status))
		buffer.WriteString(fmt.Sprintf("Names: %s\n", container.Names))
		if (len(container.Ports) > 0) {
			buffer.WriteString("Ports: \n")
			for _, port := range container.Ports {

				buffer.WriteString(fmt.Sprintf("\t type: %s IP: %s private: %d public: %d\n", port.Type, port.IP, port.PrivatePort, port.PublicPort))
			}
		}
		buffer.WriteString(fmt.Sprint("\n"))
	}

	io.WriteString(w, buffer.String())
}

func setupDockerClient(env string) *docker.Client {

	// first get the environment from the config
	cfg, err := config.GetDockerEnvironmentConfig(env)

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