package main

import (
	"io"
	"fmt"
	"os/exec"
	"log"
	"net/http"
	"html/template"
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

type Page struct {
	Buttons []string
}

func buttonsHandler(file_name string, w http.ResponseWriter, r *http.Request) {
	var buttons []string
	for key := range cmdList {
		buttons = append(buttons, key)
	}
	p := &Page{Buttons: buttons}
	t, err := template.ParseFiles(file_name)
	if err != nil {
		log.Printf("buttonsHandler, ParseFiles: %v", err)
		return
	}
	t.Execute(w, p)
}

func main() {
	cmdList = make(map[string]command)
	cmdList["screenoff"] = command{"/usr/bin/xset", []string{"dpms", "force", "off"}, ""}
	cmdList["volumeup"] =   command{"/usr/bin/pulseaudio-ctl", []string{"up"}, ""}
	cmdList["volumedown"] = command{"/usr/bin/pulseaudio-ctl", []string{"down"}, ""}
	cmdList["volumeinfo"] = command{"/usr/bin/pulseaudio-ctl", []string{"full-status"}, ""}
	cmdList["F"] = command{"/usr/bin/xvkbd", []string{"-text"}, "F"}
	cmdList["F11"] = command{"/usr/bin/xvkbd", []string{"-text"}, "\\{F11}"}
	r := mux.NewRouter()
	for i := range cmdList {
		r.HandleFunc("/gomme-api/{cmd}" + cmdList[i].Pattern, CmdHandler)
	}
	r.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		buttonsHandler("html/gomme.html", w, r)
	})
	r.HandleFunc("/gomme.js", func (w http.ResponseWriter, r *http.Request) {
		buttonsHandler("html/gomme.js", w, r)
	})
	fs := http.FileServer(http.Dir("html"))
	r.PathPrefix("/js/").Handler(fs)
	r.PathPrefix("/css/").Handler(fs)
	r.PathPrefix("/fonts/").Handler(fs)
	log.Fatal(http.ListenAndServe(":8080", r)) // Change to localhost on prod
}
