package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/raidancampbell/pcabinet/conf"
)

var _ tea.Model = &namingDialog{}
type namingDialog struct {
	textInput textinput.Model
	err error

	service conf.Service
	profiling []profilingOption
	parentModel tea.Model
}

type errMsg error

func NewNamingDialog(service conf.Service, profiling []profilingOption, parent tea.Model) tea.Model {
	ti := textinput.NewModel()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return namingDialog{
		textInput: ti,
		err:       nil,
		service: service,
		profiling: profiling,
		parentModel: parent,
	}
}


func (n namingDialog) Init() tea.Cmd {
	return textinput.Blink
}

func (n namingDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return n, tea.Quit
		case tea.KeyEsc:
			return n.parentModel, nil
		case tea.KeyEnter:
			return NewDownloadSpinner(n.service, n.profiling, n.textInput.Value(), n), spinner.Tick
		}

	// We handle errors just like any other message
	case errMsg:
		n.err = msg
		return n, nil
	}

	n.textInput, cmd = n.textInput.Update(msg)
	return n, cmd
}

func (n namingDialog) View() string {
	return fmt.Sprintf(
		"Enter a short description to insert into the filename:\n\n%s\n\n%s",
		n.textInput.View(),
		"(^C to quit)",
	) + "\n"
}