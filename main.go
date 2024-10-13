package main

import (
	"compass/client"
	"compass/global"
	"compass/scope"
	"compass/session"
	"compass/views/login"
	"compass/views/tokencreate"
	"compass/views/zoneselector"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var SSIHost string

func init() {
	flag.StringVar(&SSIHost, "h", "", "Where is SSI?")
}

func main() {
	flag.Parse()
	if SSIHost == "" {
		log.Fatal("No host provided")
	}
	global.SSIHost = SSIHost
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("We ran into an error: %v", err)
		os.Exit(1)
	}
}

type Model struct {
	loginView    tea.Model
	resourceView tea.Model
	tokenForm    tea.Model

	zonesList []scope.ZoneData
	client    *client.Client
	session   session.SessionInterface

	mode int

	output []string
}

const (
	ModeLogin = iota
	ModeSelectResources
	ModeCreateToken
)

func initialModel() tea.Model {
	c := client.NewClient(SSIHost,
		client.WithHeader("Content-Type", "application/json"),
	)
	return Model{
		client:    c,
		loginView: login.New(),
		zonesList: []scope.ZoneData{},
		mode:      ModeLogin,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.output = m.output[:0]

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q", "ctrl+c":
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case scope.TokenInformation:
		m.output = append(m.output, fmt.Sprintf("%+v", msg))

	case zoneselector.PermissionCollection:
		m.tokenForm = tokencreate.New(m.client, msg)
		m.mode = ModeCreateToken

	case zoneselector.Model:
		m.resourceView = msg
		m.mode = ModeSelectResources

	case *client.Client:
		m.client = msg

	case *session.Session:
		m.session = msg
		m.client.Config(
			client.WithInterceptor(func(c *client.Client) {
				c.SetHeader("Authorization", m.session.GetToken())
			}),
			client.WithStatusHandler(401, func() {
				fmt.Print("Session was invalid. Next call will prompt a sign in")
				defer m.session.Save()
				m.session.SetToken("")
			}),
		)
		return m, createZonesView(m)

	case error:
		m.output = append(m.output, "Error: \n"+msg.Error())
	}

	if m.session == nil || m.session.GetToken() == "" {
		m.mode = ModeLogin
	}

	if m.mode == ModeLogin {
		newLoginview, cmd := m.loginView.Update(msg)
		loginViewModel, ok := newLoginview.(login.Model)
		if !ok {
			panic("Could not assert login view")
		}
		m.loginView = loginViewModel
		return m, cmd
	}

	if m.mode == ModeSelectResources {
		newZonesView, cmd := m.resourceView.Update(msg)
		zonesViewModel, ok := newZonesView.(zoneselector.Model)
		if !ok {
			panic("Could not assert zones view")
		}
		m.resourceView = zonesViewModel
		return m, cmd
	}

	if m.mode == ModeCreateToken {
		newTokenForm, cmd := m.tokenForm.Update(msg)
		tokenFormModel, ok := newTokenForm.(tokencreate.Model)
		if !ok {
			panic("Could not assert token view")
		}
		m.tokenForm = tokenFormModel
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.mode == ModeLogin {
		return lipgloss.JoinVertical(
			lipgloss.Top,
			m.loginView.View(),
			strings.Join(m.output, ", "),
		)
	}
	if m.mode == ModeSelectResources {
		return lipgloss.JoinVertical(
			lipgloss.Top,
			m.resourceView.View(),
			strings.Join(m.output, ",\n"),
		)
	}
	if m.mode == ModeCreateToken {
		return lipgloss.JoinVertical(
			lipgloss.Top,
			m.tokenForm.View(),
			strings.Join(m.output, ",\n"),
		)
	}
	return "Loading..."
}

func getZones(c client.ClientInterface) ([]scope.ZoneData, error) {
	zones := []scope.ZoneData{}
	err := c.Get("/zones", &zones)
	if err != nil {
		return zones, err
	}

	return zones, nil
}
func createZonesView(m Model) tea.Cmd {
	zones, err := getZones(m.client)
	if err != nil {
		return func() tea.Msg {
			return err
		}
	}
	zonesView := zoneselector.New(zones)
	return func() tea.Msg {
		return zonesView
	}

}
