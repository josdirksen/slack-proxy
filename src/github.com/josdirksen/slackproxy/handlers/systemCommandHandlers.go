package handlers
import (
	"net/http"
	"fmt"
	"bytes"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/host"
	"io"
)

type SystemHandler struct {
}

func NewSystemHandler() *SystemHandler {
	p := new(SystemHandler)
	return p
}

func (sh SystemHandler) HandleCommand(cmdToExecute *Command, w http.ResponseWriter) {
	switch cmdToExecute.SlackCommand {
	case "mem" : handleMemCommand(w)
	case "host" : handleHostCommand(w)
	}
}

func handleHostCommand(w http.ResponseWriter) {
	var buffer bytes.Buffer

	info, _ := host.HostInfo()
	buffer.WriteString(fmt.Sprintf("Boottime: %v\n", info.BootTime))
	buffer.WriteString(fmt.Sprintf("Hostname: %v\n", info.Hostname))
	buffer.WriteString(fmt.Sprintf("Uptime: %v\n", info.Uptime))

	io.WriteString(w, buffer.String())
}

func handleMemCommand(w http.ResponseWriter) {
	v, _ := mem.VirtualMemory()
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf(fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f\n", v.Total, v.Free, v.UsedPercent)))
	io.WriteString(w, buffer.String())
}