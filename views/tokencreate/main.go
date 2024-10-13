package tokencreate

import (
	"compass/client"
	"compass/scope"
	"compass/views/zoneselector"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputName string

const (
	INPUT_PASS     InputName = "PASSWORD"
	INPUT_NAME     InputName = "NAME"
	INPUT_VALIDITY InputName = "VALIDITY"
)

type Input struct {
	model textinput.Model
	name  InputName
}

type Model struct {
	inputs        []Input
	selectedInput int
	client        client.ClientInterface
	authReq       byte
	permissions   map[string]scope.Permission
}

func New(c client.ClientInterface, perms zoneselector.PermissionCollection) tea.Model {
	nameInput := textinput.New()
	nameInput.Prompt = "Name: "
	name := Input{
		model: nameInput,
		name:  INPUT_NAME,
	}

	validityInput := textinput.New()
	validityInput.Prompt = "Validity (seconds): "
	validityInput.Validate = intValidator
	validity := Input{
		model: validityInput,
		name:  INPUT_VALIDITY,
	}

	passInput := textinput.New()
	passInput.Prompt = "Password: "
	passInput.EchoMode = textinput.EchoPassword
	pass := Input{
		model: passInput,
		name:  INPUT_PASS,
	}

	inputs := []Input{name, validity, pass}

	return Model{
		client:        c,
		selectedInput: 0,
		inputs:        inputs,
		permissions:   perms,
	}
}

func intValidator(s string) error {
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}

func createToken(m Model) tea.Cmd {
	tokenData := scope.NewTokenRequest{}
	s := scope.SPScope{}
	s.Permissions = m.permissions
	s.Clients = []scope.Client{"SynkzoneSSI"}

	for _, input := range m.inputs {
		switch input.name {
		case INPUT_PASS:
			val := input.model.Value()
			charr := strings.Split(string(val), "")
			tokenData.Password = charr
		case INPUT_NAME:
			val := input.model.Value()
			tokenData.TokenName = &val
			s.Name = val
		case INPUT_VALIDITY:
			val := input.model.Value()
			num, _ := strconv.ParseInt(val, 10, 64)
			tokenData.Validity = &num
			s.Name = val
		}
	}

	tokenData.Scope = s

	tokenRes := scope.TokenInformation{}
	err := m.client.Post("/tokens", tokenData, &tokenRes)
	if err != nil {
		return func() tea.Msg {
			return err
		}
	}

	return func() tea.Msg {
		return tokenRes
	}
}

func (m Model) View() string {
	form := strings.Builder{}
	for _, i := range m.inputs {
		form.WriteString(i.model.View() + "\n")
	}
	return lipgloss.JoinVertical(
		lipgloss.Top,
		form.String(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, createToken(m)
		case "tab":
			m.selectedInput = min(m.selectedInput+1, len(m.inputs)-1)
		case "shift+tab":
			m.selectedInput = max(m.selectedInput-1, 0)
		}
	}

	for index := range m.inputs {
		if index != m.selectedInput {
			m.inputs[index].model.Blur()
		}
	}

	m.inputs[m.selectedInput].model.Focus()

	var cmd tea.Cmd
	m.inputs[m.selectedInput].model, cmd = m.inputs[m.selectedInput].model.Update(msg)
	return m, cmd
}

func (m Model) Init() tea.Cmd {
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
