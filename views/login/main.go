package login

import (
	"compass/client"
	"compass/global"
	"compass/scope"
	"compass/session"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	REQ_NONE uint8 = 1 << iota
	REQ_PASS
	REQ_OTP
)

type InputName string

const (
	INPUT_USER InputName = "USER"
	INPUT_PASS InputName = "PASSWORD"
	INPUT_OTP  InputName = "OTP"
)

type Input struct {
	model textinput.Model
	name  InputName
}

type Model struct {
	inputs        []Input
	selectedInput int
	client        client.ClientInterface
	session       session.SessionInterface
	authReq       byte
}

func New() tea.Model {
	s := session.New(
		session.WithStore(session.FileStorage(global.CONFIG_FILE)),
	)

	c := client.NewClient(global.SSIHost,
		client.WithHeader("Content-Type", "application/json"),
	)

	return Model{
		session:       s,
		client:        c,
		authReq:       REQ_NONE,
		selectedInput: 0,
	}
}

type LogonRequest struct {
	AcceptSessions   *bool          `json:"acceptSessions,omitempty"`
	AutoLogon        *bool          `json:"autoLogon,omitempty"`
	LogonSessionId   *string        `json:"logonSessionId,omitempty"`
	OrganizationName *string        `json:"organizationName,omitempty"`
	Password         *[]string      `json:"password,omitempty"`
	PasswordAsString *string        `json:"passwordAsString,omitempty"`
	SavePassword     *bool          `json:"savePassword,omitempty"`
	Scope            *scope.SPScope `json:"scope,omitempty"`
	SetCookie        *bool          `json:"setCookie,omitempty"`
	Timeout          *int64         `json:"timeout,omitempty"`
	TimeoutInSeconds *int64         `json:"timeoutInSeconds,omitempty"`
	Token            *string        `json:"token,omitempty"`
	Username         *string        `json:"username,omitempty"`
	VerificationCode *string        `json:"verificationCode,omitempty"`
}

type LogonResult struct {
	ExpiresInMillis          *int64  `json:"expiresInMillis,omitempty"`
	IdpURI                   *string `json:"idpURI,omitempty"`
	LogonOK                  bool    `json:"logonOK"`
	PasswordChangeIsRequired *bool   `json:"passwordChangeIsRequired,omitempty"`
	SessionId                *string `json:"sessionId,omitempty"`
}

func getSession(m Model) tea.Cmd {
	defer m.session.Save()

	credentials := LogonRequest{}

	validity := int64(60 * 15)
	credentials.TimeoutInSeconds = &validity

	for _, input := range m.inputs {
		switch input.name {
		case INPUT_USER:
			val := ""
			val = input.model.Value()
			credentials.Username = &val
		case INPUT_PASS:
			val := ""
			val = input.model.Value()
			charr := strings.Split(string(val), "")
			credentials.Password = &charr
		case INPUT_OTP:
			val := ""
			val = input.model.Value()
			credentials.VerificationCode = &val
		}
	}

	loginRes := LogonResult{}
	if err := m.client.Post("/access/logon", credentials, &loginRes); err != nil {
		return func() tea.Msg {
			return err
		}
	}

	m.session.SetToken(*loginRes.SessionId)

	return func() tea.Msg {
		return m.session
	}
}

func getInputs(m Model) tea.Cmd {
	u := textinput.New()
	u.Prompt = "User name: "
	inputs := []Input{
		{
			name:  INPUT_USER,
			model: u,
		},
	}
	if scope.ConfHasOpt(m.authReq, REQ_OTP) {
		i := textinput.New()
		i.Prompt = "One time password: "
		inputs = append(inputs, Input{
			name:  INPUT_OTP,
			model: i,
		})
	}

	if scope.ConfHasOpt(m.authReq, REQ_PASS) {
		i := textinput.New()
		i.Prompt = "Password: "
		i.EchoMode = textinput.EchoPassword
		inputs = append(inputs, Input{
			name:  INPUT_PASS,
			model: i,
		})
	}

	return func() tea.Msg {
		return inputs
	}
}

type AuthReq byte
type AuthorizationRequirements struct {
	CanSavePassword  *bool   `json:"canSavePassword,omitempty"`
	OrganizationName *string `json:"organizationName,omitempty"`
	OtpRequired      *bool   `json:"otpRequired,omitempty"`
	OtpType          *string `json:"otpType,omitempty"`
	PasswordRequired *bool   `json:"passwordRequired,omitempty"`
	WorldName        *string `json:"worldName,omitempty"`
}

func getAuthReq(m Model) tea.Cmd {
	authReq := byte(0)
	authReqData := AuthorizationRequirements{}
	if err := m.client.Get("/access/authorizationrequirements", &authReqData); err != nil {
		return func() tea.Msg {
			return err
		}
	}

	if *authReqData.OtpRequired {
		authReq |= REQ_OTP
	}

	if *authReqData.PasswordRequired {
		authReq |= REQ_PASS
	}

	return func() tea.Msg {
		return AuthReq(authReq)
	}
}

func (m Model) View() string {
	credentials := strings.Builder{}
	for _, i := range m.inputs {
		credentials.WriteString(i.model.View())
		credentials.WriteString("\n")
	}
	return lipgloss.JoinVertical(
		lipgloss.Top,
		credentials.String(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.session.Load()
	if m.session.GetToken() != "" {
		return m, func() tea.Msg {
			return m.session
		}
	}

	switch msg := msg.(type) {
	case AuthReq:
		m.authReq = byte(msg)
		return m, getInputs(m)
	case []Input:
		m.inputs = msg
		m.inputs[0].model.Focus()
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, getSession(m)
		case "tab":
			m.selectedInput = min(m.selectedInput+1, len(m.inputs)-1)
		case "shift+tab":
			m.selectedInput = max(m.selectedInput-1, 0)
		}
	}

	if m.authReq == REQ_NONE {
		return m, getAuthReq(m)
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
