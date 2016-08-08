package main

import (
	"io"
	"fmt"
	"os/exec"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

var cmdList map[string]string

func CmdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cmdString := cmdList[vars["cmd"]]
	if cmdString == "" {
		log.Print("Command not found for ", vars["cmd"])
		return
	}
	cmd := exec.Command("/bin/sh", "-c", cmdString)

	cmd.Stdout = io.Writer(w)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(w, "Hello, %s", "HERE GOES ERR MSG")
		log.Print("HERE GOES ERR MSG")
	}
}

func main() {
	cmdList = make(map[string]string)
	cmdList["screenoff"] = "/usr/bin/xset dpms force off"
	cmdList["volumeup"] = "/usr/bin/pulseaudio-ctl up"
	cmdList["volumedown"] = "/usr/bin/pulseaudio-ctl down"
	cmdList["volumeinfo"] = "/usr/bin/pulseaudio-ctl full-status"
	r := mux.NewRouter()
	r.HandleFunc("/{cmd}", CmdHandler)
	log.Fatal(http.ListenAndServe(":8080", r))
}
