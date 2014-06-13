package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"bytes"
	"encoding/json"
	"github.com/dotcloud/docker/api/client"
	"text/template"
)


var homeTempl = template.Must(template.ParseFiles("home.html"))


func CreateDockerCli() (* client.DockerCli, *bytes.Buffer){
	unixsock := "/var/run/docker.sock"
	defaultHost := fmt.Sprintf("unix://%s", unixsock)
	protoAddrParts := strings.SplitN(defaultHost, "://", 2)
	bufOut := bytes.NewBuffer(nil)
	bufErr := bytes.NewBuffer(nil)
	return client.NewDockerCli(os.Stdin, bufOut, bufErr, protoAddrParts[0], protoAddrParts[1], nil), bufOut
}


func ServeDocker(w http.ResponseWriter, r *http.Request){
	cli, _ := CreateDockerCli()
	//queryParams := r.URL.Query()
	//command := queryParams["command"][0]
	//command_name := strings.Split(command, " ")

	cli.CmdRun("ubuntu")//, command_name[0], command_name[1])

	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method nod allowed", 405)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTempl.Execute(w, r.Host)

}

type Test struct {
	Jsonrpc string
	Method string
	Params []string
	Id int
}

func ServeCmd(w http.ResponseWriter, r *http.Request){
	cli, bufOut := CreateDockerCli()
	// queryParams := r.URL.Query()
	// command := queryParams["command"][0]
	// command_name := strings.Split(command, " ")
	// cli.CmdRun("ubuntu", command_name[0], command_name[1])
	//	fmt.Println(r)
	var m Test

	dec := json.NewDecoder(r.Body)
	dec.Decode(&m)
	fmt.Println(m.Method)
	// fmt.Println(r.Body.String())
	// a, _ := httputil.DumpRequest(r.Body, true)
	// fmt.Println(string(a))
	fmt.Fprintln(w, "holaa")

}


func main() {
	http.HandleFunc("/", ServeDocker)
	http.HandleFunc("/runDocker", ServeCmd)
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	//http.HandleFunc("/ws", serveWs)
	http.ListenAndServe("localhost:4000", nil)
}
