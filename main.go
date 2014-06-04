package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"bytes"
	"github.com/dotcloud/docker/api/client"
)

type Hello struct{}

type Struct struct {
	Greeting string
	Punct string
	Who string
}

func (s Struct) ServeHTTP(w http.ResponseWriter, r *http.Request){
	unixsock := "/var/run/docker.sock"
	defaultHost := fmt.Sprintf("unix://%s", unixsock)
	protoAddrParts := strings.SplitN(defaultHost, "://", 2)
	bufOut := bytes.NewBuffer(nil)
	bufErr := bytes.NewBuffer(nil)
	cli := client.NewDockerCli(os.Stdin, bufOut, bufErr, protoAddrParts[0], protoAddrParts[1], nil)
	queryParams := r.URL.Query()
	command := queryParams["command"][0]
	command_name := strings.Split(command, " ")

	cli.CmdRun("ubuntu", command_name[0], command_name[1])

	fmt.Fprintln(w, bufOut)
}

func (h Hello) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request) {
	fmt.Println(r)
	fmt.Fprint(w, "Hello!")
}

func main() {
	http.Handle("/string", Struct{"this", "is", "test"})
	http.ListenAndServe("localhost:4001", nil)
}
