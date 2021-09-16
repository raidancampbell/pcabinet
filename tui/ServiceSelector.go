package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/raidancampbell/pcabinet/conf"
	"github.com/sirupsen/logrus"
)

// ensure that the Model interface is implemented.
var _ tea.Model = &ServiceSelector{}

// ServiceSelector gives a menu for the user to select the service to profile
//TODO add a text input and fuzzy matching
type ServiceSelector struct {
	Options     []conf.Service
	idx         int
}

func NewServiceSelector(services []conf.Service) tea.Model {

	return &ServiceSelector{
		Options:     services,
	}
}

func (s *ServiceSelector) Init() tea.Cmd {
	if len(s.Options) == 0 {
		logrus.Error("no options to select")
		tea.Quit()
	}
	return nil
}

func (s *ServiceSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return s, tea.Quit
		case "up", "k":
			if s.idx % len(s.Options) == 0 {
				s.idx = len(s.Options) - 1
			} else {
				s.idx--
			}
		case "down", "j":
			s.idx++
		case "enter", " ", "right":

			return NewProfileSelector(s.Options[s.idx % len(s.Options)], s), nil
		}
	}
	return s, nil
}

func (s *ServiceSelector) View() string {
	str := "Which server would you like to profile?\n\n"
	for i := 0; i < len(s.Options); i++ {
		cursor := " "
		if s.idx % len(s.Options) == i {
			cursor = ">"
		}
		str += fmt.Sprintf("%s [%s] %s\n", cursor, s.Options[i].Name, s.Options[i].Endpoint)
	}
	str += "\nPress q to quit.\n"
	return str
}
