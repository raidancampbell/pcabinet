package main

import (
	"bytes"
	"fmt"
	"github.com/raidancampbell/pcabinet/conf"
	"github.com/raidancampbell/pcabinet/tui"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	conf.Initialize()
	buf := bytes.Buffer{}
	logrus.SetOutput(&buf)

	// no need to open HTTP listen ports by default. this is only useful for debugging pcabinet itself.
	if os.Getenv("PCABINET_DEBUG") == "yes" {
		go defaultWebServer()
	}

	model := tui.NewServiceSelector(conf.C.Services)

	p := tea.NewProgram(model)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	fmt.Print(buf.String())
}

// defaultWebServer exists so that I can test it on itself. I don't wanna keep another long-running debug service around.
func defaultWebServer() {
	logrus.Fatal(http.ListenAndServe("127.0.0.1:8080", http.DefaultServeMux))
}
