package tui

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/raidancampbell/pcabinet/conf"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	url2 "net/url"
	"os"
	"path"
	"strings"
	"time"
)

var _ tea.Model = &downloadSpinner{}

type downloadSpinner struct {
	spinner spinner.Model
	err     error
	downloadComplete chan error

	description string
	service     conf.Service
	profiling   profilingOption
	parentModel tea.Model
}

func NewDownloadSpinner(service conf.Service, profiling profilingOption, description string, parent tea.Model) tea.Model {
	sp := spinner.NewModel()

	downloadComplete := make(chan error)
	go doDownload(service, profiling, description, downloadComplete)

	return &downloadSpinner{
		spinner:     sp,
		err:         nil,
		downloadComplete: downloadComplete,
		description: description,
		service:     service,
		profiling:   profiling,
		parentModel: parent,
	}
}

func doDownload(service conf.Service, profiling profilingOption, description string, complete chan error) {
	filename := fmt.Sprintf("%s.%s.%s.%s", service.Name, time.Now().Format("2006-01-02T15-04-05"), description, profiling.endpointSuffix)
	filename = strings.ReplaceAll(filename, " ", "-")
	err := os.Mkdir(service.Name, 0755)
	dirExists := errors.Is(err, os.ErrExist)
	if err != nil && !dirExists {
		logrus.WithError(err).WithField("directory", service.Name).Error("unable to create directory for file")
		complete <- err
		return
	}
	out, err := os.Create(path.Join(service.Name, filename))
	if err != nil {
		logrus.WithError(err).WithField("filename", filename).Error("unable to create file")
		if !dirExists {
			os.Remove(service.Name)
		}
		os.Remove(path.Join(service.Name, filename))
		complete <- err
		return
	}
	defer out.Close()
	u, err := url2.Parse(service.Endpoint)
	if err != nil {
		logrus.WithError(err).WithField("endpoint", service.Endpoint).Error("unable to parse endpoint as URL")
		if !dirExists {
			os.Remove(service.Name)
		}
		os.Remove(path.Join(service.Name, filename))
		complete <- err
		return
	}
	u.Path = path.Join(u.Path, profiling.endpointSuffix)
	url := u.String()
	resp, err := http.Get(url)
	if err != nil {
		logrus.WithError(err).WithField("url", url).Error("Failed to get profile")
		if !dirExists {
			os.Remove(service.Name)
		}
		os.Remove(path.Join(service.Name, filename))
		complete <- err
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logrus.WithError(err).WithField("url", url).WithField("status_code", resp.StatusCode).Warn("non-200 returned")
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logrus.WithError(err).WithField("url", url).Error("Failed to write response to file")
		complete <- err
		return
	}
	logrus.Infof("successfully wrote data to file '%v'", path.Join(service.Name, filename))
	close(complete)
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

	// this is a massive abuse of the system: I'm treating the spinner spin messages as a non-blocking poll whether the download is complete
	// TODO: figure out how to inject commands into the framework outside of the update method
	select {
	case <- d.downloadComplete:
		return d, tea.Quit
	default:
	}

	d.spinner, cmd = d.spinner.Update(msg)
	return d, cmd
}

func (d *downloadSpinner) View() string {
	return fmt.Sprintf("%s Capturing (^C to quit)...\n", d.spinner.View())
}
