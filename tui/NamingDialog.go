package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = &namingDialog{}
type namingDialog struct {
	textInput textinput.Model
	err error

	endpoint string
	parentModel tea.Model
}

type errMsg error

func NewNamingDialog(endpoint string, parent tea.Model) namingDialog {
	ti := textinput.NewModel()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return namingDialog{
		textInput: ti,
		err:       nil,
		endpoint: endpoint,
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
			//TODO:
			// add a downloader next
			panic("f")
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
		"(esc to quit)",
	) + "\n"
}