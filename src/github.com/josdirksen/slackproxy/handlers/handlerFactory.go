package handlers
import (
	"net/http"
	"github.com/josdirksen/slackproxy/config"
	"strings"
	"net/url"
)

// Define a simple interface. Which provides access to the command that needs to be
// executed and the response writer, so that the handler can send information back.
type CommandHandler interface {
	HandleCommand(cmdToExecute *Command, w http.ResponseWriter)
}

// We expect the post input to look something like this:
// token=5cLHiZjpWaRDb0fP6ka02XCR&team_id=T0001&team_domain=example&channel_id=C2147483705&channel_name=test&user_id=U2147483697&user_name=Steve&command=/docker&text=local+ps
type Command struct {
	Token        string
	Team_id      string
	Team_domain  string
	Channel_id   string
	Channel_name string
	User_id      string
	User_name    string
	Command      string
	Text         string
	Environment  string
	SlackCommand string
}

// Parse a key value map with the correct header names to a command
func NewCommand(kvmap map[string]string) *Command {
	c := Command{kvmap["token"], kvmap["team_id"], kvmap["team_domain"],
		kvmap["channel_id"], kvmap["channel_name"], kvmap["user_id"],
		kvmap["user_name"], kvmap["command"], kvmap["text"],
		kvmap["Environment"], kvmap["SlackCommand"]}
	return &c
}

// parse the input to a command
func ParseInput(input url.Values) *Command {
	m := make(map[string]string)

	cmdParts := strings.Split(input.Get("text"), " ")
	m["Environment"] = cmdParts[0]
	m["SlackCommand"] = cmdParts[1]

	for key, entry := range input {
		m[key] = entry[0]
	}

	return NewCommand(m)
}


type DummyHandler struct {}
func NewDummyHandler() *DummyHandler {
	p := new(DummyHandler)
	return p
}

func (dh DummyHandler) HandleCommand(cmdToExecute *Command, w http.ResponseWriter) {}

func GetHandler(handlerName string, config *config.Configuration) CommandHandler {
	switch handlerName {
	case "docker": return NewDockerHandler(config)
	default: return NewDummyHandler()
	}
}





