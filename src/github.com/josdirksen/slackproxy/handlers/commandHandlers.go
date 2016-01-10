package handlers

import (
	"io"
	"net/http"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"bytes"
	"time"
	"strings"
	"net/url"
)

type CommandHandler interface {



}

// We don't do much input checking. if we can't parse it, we just fail the command.
// We expect the post input to look something like this:
// token=5cLHiZjpWaRDb0fP6ka02XCR&team_id=T0001&team_domain=example&channel_id=C2147483705&channel_name=test&user_id=U2147483697&user_name=Steve&command=/docker&text=local+ps
type Command struct {
	Token string
	Team_id string
	Team_domain string
	Channel_id string
	Channel_name string
	User_id string
	User_name string
	Command string
	Text string
	Environment string
	DockerCommand string
}

// Parse a key value map with the correct header names to a command
func NewCommand(kvmap map[string]string) *Command {
	c := Command{kvmap["token"], kvmap["team_id"], kvmap["team_domain"],
		kvmap["channel_id"], kvmap["channel_name"], kvmap["user_id"],
		kvmap["user_name"], kvmap["command"], kvmap["text"],
		kvmap["Environment"], kvmap["DockerCommand"]}
	return &c
}

// parse the input to a command
func ParseInput(input url.Values) *Command {
	m := make(map[string]string)

	cmdParts := strings.Split(input.Get("text"), " ")
	m["Environment"] = cmdParts[0]
	m["DockerCommand"] = cmdParts[1]

	for key, entry := range input {
		m[key] = entry[0]
	}

	return NewCommand(m)
}

func HandleCommand(cmdToExecute *Command, client *docker.Client, w http.ResponseWriter) {
	switch cmdToExecute.DockerCommand {
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