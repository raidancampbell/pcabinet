package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/raidancampbell/pcabinet/conf"
)

// profileSelector gives a menu for the user to select the profile to capture
// it does not implement the tea.Model interface, since it's a sub-namingDialog
type profileSelector struct {
	options []profilingOption
	idx     int
	chosen  []int
	service conf.Service

	parentModel tea.Model
}

var profiles = []profilingOption{
	{"CPU profile", "profile"},
	{"heap profile", "heap"},
	{"goroutine profile", "goroutine"},
	{"block profile", "block"},
	{"mutex profile", "mutex"},
	{"trace", "trace"},
}

var _ tea.Model = &profileSelector{}

func NewProfileSelector(service conf.Service, parentModel tea.Model) tea.Model {
	return &profileSelector{
		options:     profiles,
		chosen:      []int{},
		service:     service,
		parentModel: parentModel,
	}
}

type profilingOption struct {
	name           string
	endpointSuffix string
}

func (s *profileSelector) Init() tea.Cmd {
	return nil
}

func (s *profileSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return s, tea.Quit
		case "q", "left":
			return s.parentModel, nil
		case "up", "k":
			if s.idx%len(s.options) == 0 {
				s.idx = len(s.options) - 1
			} else {
				s.idx--
			}
		case "down", "j":
			s.idx++
		case " ":
			removed := false
			for i, chosen := range s.chosen {
				if chosen == s.idx%len(s.options) {
					removed = true
					s.chosen = append(s.chosen[:i], s.chosen[i+1:]...)
				}
			}
			if !removed {
				s.chosen = append(s.chosen, s.idx%len(s.options))
			}
		case "enter", "right":
			opts := []profilingOption{}
			for _, i := range s.chosen {
				opts = append(opts, s.options[i])
			}
			if len(opts) == 0 {
				opts = append(opts, s.options[s.idx%len(s.options)])
			}
			return NewNamingDialog(s.service, opts, s), nil
		}
	}
	return s, nil
}

func (s *profileSelector) View() string {
	str := "Which profile would you like to capture?\n\n"
	for i := 0; i < len(s.options); i++ {
		cursor := " "
		if s.idx%len(s.options) == i {
			cursor = ">"
		}
		isSelected := " "
		for _, chosen := range s.chosen {
			if chosen == i {
				isSelected = "X"
			}
		}
		str += fmt.Sprintf("%s [%s] %s\n", cursor, isSelected, s.options[i].name)
	}
	str += "\nPress q or left arrow to go back, space to select, right arrow or enter to continue\n"
	return str
}
