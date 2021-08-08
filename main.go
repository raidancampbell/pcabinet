package main

import (
	"github.com/raidancampbell/pcabinet/conf"
	"github.com/raidancampbell/pcabinet/tui"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg := conf.Initialize()

	//TODO: ServiceSelector should operate on services
	serviceNames := map[string]string{}
	for name, service := range cfg.Services {
		serviceNames[name] = service.Endpoint
	}
	p := tea.NewProgram(&tui.ServiceSelector{Options: serviceNames})
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	//TODO: figure out how to chain together bubbles in bubbletea
	//TODO: after a value is selected, display a download spinner
	//TODO: HTTP GET that endpoint, write to some file
}
