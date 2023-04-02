package tui

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/raidancampbell/pcabinet/conf"
	"github.com/raidancampbell/pcabinet/internal"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	url2 "net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

var _ tea.Model = &downloadSpinner{}

type downloadSpinner struct {
	spinner          spinner.Model
	err              error
	downloadComplete chan error
	estCompletion    time.Time

	description string
	service     conf.Service
	profiling   []profilingOption
	parentModel tea.Model
}

func NewDownloadSpinner(service conf.Service, profiling []profilingOption, description string, parent tea.Model) tea.Model {
	sp := spinner.NewModel()

	downloadComplete := make(chan error)
	go func() {
		for _, profile := range profiling {
			doDownload(service, profile, description, downloadComplete)
		}
		close(downloadComplete)
	}()

	estCompletionTime := time.Now()
	for _, profile := range profiling {
		switch profile.endpointSuffix {
		case "profile":
			if conf.C.TestCPU {
				estCompletionTime = estCompletionTime.Add(31 * time.Second)
			} else {
				estCompletionTime = estCompletionTime.Add(30 * time.Second)
			}
		case "trace":
			estCompletionTime = estCompletionTime.Add(1 * time.Second)
		default:
		}
	}

	return &downloadSpinner{
		spinner:          sp,
		err:              nil,
		downloadComplete: downloadComplete,
		estCompletion:    estCompletionTime,
		description:      description,
		service:          service,
		profiling:        profiling,
		parentModel:      parent,
	}
}

// doDownload does the actual download of the desired profile and stores it into a file.
// The rest of this repository is scaffolding for this function.
func doDownload(service conf.Service, profiling profilingOption, description string, complete chan error) {
	// create the directory for this capture.  if it didn't exist and we fail later, we'll delete it to clean up
	dirName := path.Join(conf.C.OutputBasedir, service.Name)
	err := os.Mkdir(dirName, 0755)
	dirExists := errors.Is(err, os.ErrExist)
	if err != nil && !dirExists {
		logrus.WithError(err).WithField("directory", dirName).Error("unable to create directory for file")
		complete <- err
		return
	}

	if profiling.endpointSuffix == "profile" && conf.C.TestCPU {
		usage, err := currentCPUUsage(service, profiling)
		if err != nil {
			logrus.WithError(err).Error("unable to test CPU usage for threshold")
			if !dirExists {
				os.Remove(dirName)
			}
			complete <- err
			return
		}
		if usage < 0.05 {
			err := fmt.Errorf("measured CPU usage of %s is less than threshold of 0.05", strconv.FormatFloat(usage, 'f', 3, 64))
			logrus.Errorf(err.Error())
			if !dirExists {
				os.Remove(dirName)
			}
			complete <- err
			return
		}
	}

	// create the file.  If we fail later, we'll delete it to clean up
	filename := fmt.Sprintf("%s.%s.%s.%s", service.Name, time.Now().Format("2006-01-02T15-04-05"), description, profiling.endpointSuffix)
	filename = path.Join(dirName, strings.ReplaceAll(filename, " ", "-"))
	out, err := os.Create(filename)
	if err != nil {
		logrus.WithError(err).WithField("filename", filename).Error("unable to create file")
		os.Remove(filename)
		if !dirExists {
			os.Remove(dirName)
		}
		complete <- err
		return
	}
	defer out.Close()

	// parse the input URL.  TODO: parse this at yaml unmarshaling time and propagate it through as a `url.URL`. no reason to error here.
	u, err := url2.Parse(service.Endpoint)
	if err != nil {
		logrus.WithError(err).WithField("endpoint", service.Endpoint).Error("unable to parse endpoint as URL")
		os.Remove(filename)
		if !dirExists {
			os.Remove(dirName)
		}
		complete <- err
		return
	}

	if service.Kube != nil {
		// kubectl port-forward services/service_name u.Port()
		serviceName := strings.TrimPrefix(service.Kube.Service, "services/")
		serviceName = strings.TrimPrefix(serviceName, "service/")
		cmd := exec.Command("kubectl", "port-forward", fmt.Sprintf("service/%s", service.Kube.Service), fmt.Sprintf("%s:%s", u.Port(), u.Port()))
		if service.Kube.Namespace != "" {
			cmd.Args = append(cmd.Args, fmt.Sprintf("--namespace=%s", service.Kube.Namespace))
		}
		if service.Kube.Context != "" {
			cmd.Args = append(cmd.Args, fmt.Sprintf("--context=%s", service.Kube.Namespace))
		}

		if err := cmd.Start(); err != nil {
			logrus.WithError(err).WithField("cmd", cmd.String())
			complete <- err
			return
		}
		defer cmd.Process.Signal(os.Interrupt)
	}

	// add the profile suffix to the URL, e.g. append `/trace` if the user wanted to capture a trace
	// and invoke
	u.Path = path.Join(u.Path, profiling.endpointSuffix)
	url := u.String()
	resp, err := http.Get(url)
	if err != nil {
		logrus.WithError(err).WithField("url", url).Error("Failed to get profile")
		os.Remove(filename)
		if !dirExists {
			os.Remove(dirName)
		}
		complete <- err
		return
	}
	defer resp.Body.Close()

	// if we got a non-200 it's probably due to an existing profile running, or the server is too overloaded to service
	if resp.StatusCode != 200 {
		logrus.WithError(err).WithField("url", url).WithField("status_code", resp.StatusCode).Warn("non-200 returned")
	}

	// write the HTTP response to the file.  This should be the binary data of the profile, unless the above warning was hit.
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logrus.WithError(err).WithField("url", url).Error("Failed to write response to file")
		complete <- err
		return
	}

	// log the output.  the logger is buffered because bubbletea hijacks the term.  The buffer will be dumped to stdout after bubbletea is done
	logrus.Infof("successfully wrote data to file %v", filename)
}

// currentCPUUsage takes a 1-second CPU profile and analyzes it to return the current usage as a percentage-float
func currentCPUUsage(service conf.Service, profiling profilingOption) (float64, error) {
	u, err := url2.Parse(service.Endpoint)
	if err != nil {
		logrus.WithError(err).WithField("endpoint", service.Endpoint).Error("unable to parse endpoint as URL")
		return 0, err
	}
	q := u.Query()
	q.Set("seconds", "1")
	u.RawQuery = q.Encode()
	u.Path = path.Join(u.Path, profiling.endpointSuffix)
	url := u.String()
	resp, err := http.Get(url)
	if err != nil {
		logrus.WithError(err).WithField("url", url).Error("Failed to get profile")
		return 0, err
	}
	b, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		logrus.WithError(err).WithField("url", url).Error("Failed to read profile response")
		return 0, err
	}

	usage, err := internal.CPUUsage(b)
	if err != nil {
		logrus.WithError(err).Error("Failed to parse profile response")
		return 0, err
	}
	return usage, nil
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
	case _, ok := <-d.downloadComplete:
		if !ok {
			return d, tea.Quit
		}
	default:
	}

	d.spinner, cmd = d.spinner.Update(msg)
	return d, cmd
}

func (d *downloadSpinner) View() string {
	return fmt.Sprintf("%s Capturing (^C to quit)...\nestimated time remaining: %s", d.spinner.View(), time.Until(d.estCompletion).String())
}
