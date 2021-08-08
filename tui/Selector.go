package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sirupsen/logrus"
	"sort"
)

// ensure that the Model interface is implemented.
var _ tea.Model = &ServiceSelector{}

// ServiceSelector gives a menu for the user to select the service to profile
//TODO add a text input and fuzzy matching
type ServiceSelector struct {
	Options     map[string]string
	nameIndexes map[int]string
	idx         int
	chosen      int
}

func (s *ServiceSelector) Init() tea.Cmd {
	if len(s.Options) == 0 {
		logrus.Error("no options to select")
		tea.Quit()
	}
	// initialize data
	s.nameIndexes = map[int]string{}
	s.chosen = -1

	// maps are unsorted. this builds the index for sorted names of the options
	strs := make([]string, len(s.Options))
	i := 0
	for name := range s.Options {
		strs[i] = name
		i++
	}
	sort.Strings(strs)
	for i, str := range strs {
		s.nameIndexes[i] = str
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
			if s.idx > 0 {
				s.idx--
			}
		case "down", "j":
			if s.idx < len(s.Options)-1 {
				s.idx++
			}
		case "enter", " ":
			s.chosen = s.idx
			panic(s.Options[s.nameIndexes[s.idx]])
		}
	}
	return s, nil
}

func (s *ServiceSelector) View() string {
	str := "Which server would you like to profile?\n\n"
	for i := 0; i < len(s.nameIndexes); i++ {
		cursor := " "
		if s.idx == i {
			cursor = ">"
		}
		str += fmt.Sprintf("%s [%s] %s\n", cursor, s.nameIndexes[i], s.Options[s.nameIndexes[i]])
	}
	str += "\nPress q to quit.\n"
	return str
}
