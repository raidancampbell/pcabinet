package main

import (
	"github.com/raidancampbell/pcabinet/conf"
	"github.com/raidancampbell/pcabinet/tui"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	_ "net/http/pprof"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg := conf.Initialize()

	go defaultWebServer()

	model := tui.NewServiceSelector(cfg.Services)

	p := tea.NewProgram(model)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	//TODO: figure out how to chain together bubbles in bubbletea
	//TODO: after a value is selected, display a download spinner
	//TODO: HTTP GET that endpoint, write to some file
}

// defaultWebServer exists so that I can test it on itself. I don't wanna keep another long-running debug service around.
func defaultWebServer() {
	logrus.Fatal(http.ListenAndServe("127.0.0.1:8080", http.DefaultServeMux))
}
