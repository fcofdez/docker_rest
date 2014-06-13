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
	// unixsock := "/var/run/docker.sock"
	// defaultHost := fmt.Sprintf("unix://%s", unixsock)
	defaultHost := "tcp://192.168.59.103:2375"
	protoAddrParts := strings.SplitN(defaultHost, "://", 2)
	bufOut := bytes.NewBuffer(nil)
	bufErr := bytes.NewBuffer(nil)
	return client.NewDockerCli(os.Stdin, bufOut, bufErr, protoAddrParts[0], protoAddrParts[1], nil), bufOut
}


func ServeDocker(w http.ResponseWriter, r *http.Request){

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

type RequestParams struct {
	Jsonrpc string
	Method string
	Params []string
	Id int
}

func ServeCmd(w http.ResponseWriter, r *http.Request){
	cli, bufOut := CreateDockerCli()
	var m RequestParams

	dec := json.NewDecoder(r.Body)
	dec.Decode(&m)
	err := cli.CmdRun("ubuntu", m.Method, strings.Join(m.Params, " "))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintln(w, bufOut)

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
