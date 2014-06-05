package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"bytes"
	"github.com/dotcloud/docker/api/client"
)

func CreateDockerCli() (* client.DockerCli, *bytes.Buffer){
	unixsock := "/var/run/docker.sock"
	defaultHost := fmt.Sprintf("unix://%s", unixsock)
	protoAddrParts := strings.SplitN(defaultHost, "://", 2)
	bufOut := bytes.NewBuffer(nil)
	bufErr := bytes.NewBuffer(nil)
	return client.NewDockerCli(os.Stdin, bufOut, bufErr, protoAddrParts[0], protoAddrParts[1], nil), bufOut
}

func ServeDocker(w http.ResponseWriter, r *http.Request){
	cli, bufOut := CreateDockerCli()
	queryParams := r.URL.Query()
	command := queryParams["command"][0]
	command_name := strings.Split(command, " ")

	cli.CmdRun("ubuntu", command_name[0], command_name[1])

	fmt.Fprintln(w, bufOut)
}


func main() {
	http.HandleFunc("/string", ServeDocker)
	http.ListenAndServe("localhost:4000", nil)
}
