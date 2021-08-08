package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

// ProfileSelector gives a menu for the user to select the profile to capture
// it does not implement the tea.Model interface, since it's a sub-model
type ProfileSelector struct {
	Options  []string
	idx      int
	chosen   int
	endpoint string

	parentModel tea.Model
}

var _ tea.Model = &ProfileSelector{}

func NewProfileSelector(endpoint string, parentModel tea.Model) tea.Model {
	return &ProfileSelector{
		Options:  []string{"CPU profile", "heap profile", "trace"},
		chosen:   -1,
		endpoint: endpoint,
		parentModel: parentModel,
	}
}

func (s *ProfileSelector) Init() tea.Cmd {
	return nil
}

func (s *ProfileSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return s, tea.Quit
		case "q", "left":
			return s.parentModel, nil
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
			panic(s.Options[s.idx])
		}
	}
	return s, nil
}

func (s *ProfileSelector) View() string {
	str := "Which profile would you like to capture?\n\n"
	for i := 0; i < len(s.Options); i++ {
		cursor := " "
		if s.idx == i {
			cursor = ">"
		}
		str += fmt.Sprintf("%s %s\n", cursor, s.Options[i])
	}
	str += "\nPress q or left arrow to go back.\n"
	return str
}
