package main

import (
	"io"
	"fmt"
	"os/exec"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type command struct {
	BinPath string
	Args []string
	Pattern string
}

var cmdList map[string]command

func CmdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cmd, ok := cmdList[vars["cmd"]]
	if !ok { //Todo add test for that
		log.Print("Command not found for ", vars["cmd"])
		return
	}
	cmdExec := exec.Command(cmd.BinPath, cmd.Args...)

	cmdExec.Stdout = io.Writer(w)
	if err := cmdExec.Run(); err != nil {
		fmt.Fprintf(w, "error: %s", err)
		log.Print("cmd: ", cmd.BinPath, " error: ", err)
	}
}

func main() {
	cmdList = make(map[string]command)
	cmdList["screenoff"] = command{"/usr/bin/xset", []string{"dpms", "force", "off"}, ""}
	cmdList["volumeup"] =   command{"/usr/bin/pulseaudio-ctl", []string{"up"}, ""}
	cmdList["volumedown"] = command{"/usr/bin/pulseaudio-ctl", []string{"down"}, ""}
	cmdList["volumeinfo"] = command{"/usr/bin/pulseaudio-ctl", []string{"full-status"}, ""}
	r := mux.NewRouter()
	for i := range cmdList {
		r.HandleFunc("/gomme-api/{cmd}" + cmdList[i].Pattern, CmdHandler)
	}
	log.Fatal(http.ListenAndServe(":8080", r)) // Change to localhost on prod
}
