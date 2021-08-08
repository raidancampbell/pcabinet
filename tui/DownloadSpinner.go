package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/raidancampbell/pcabinet/conf"
)

var _ tea.Model = &downloadSpinner{}

type downloadSpinner struct {
	spinner spinner.Model
	err     error

	description string
	service     conf.Service
	profiling   profilingOption
	parentModel tea.Model
}

func NewDownloadSpinner(service conf.Service, profiling profilingOption, description string, parent tea.Model) tea.Model {
	sp := spinner.NewModel()

	return &downloadSpinner{
		spinner:     sp,
		err:         nil,
		description: description,
		service:     service,
		profiling:   profiling,
		parentModel: parent,
	}
}

func (d *downloadSpinner) Init() tea.Cmd {
	return nil
}

func (d *downloadSpinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return d, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		d.err = msg
		return d, nil
	}

	d.spinner, cmd = d.spinner.Update(msg)
	return d, cmd
}

func (d *downloadSpinner) View() string {
	return fmt.Sprintf("%s Capturing (^C to quit)...\n", d.spinner.View())
}
