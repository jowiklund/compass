package checkbox

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	checked bool
	focused bool
	Label   string
}

var checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

func getText(label string, checked bool, focused bool) string {
	prefix := ""
	suffix := ""

	if checked {
		prefix += "[x] "
	} else {
		prefix += "[ ] "
	}

	if focused {
		suffix += " <"
	}

	if focused {
		return checkboxStyle.Render(prefix + label + suffix)
	}
	return fmt.Sprintf(prefix+"%s", label)
}

func (m Model) View() string {
	return getText(m.Label, m.checked, m.focused)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			m.checked = !m.checked
		}
	}
	return m, nil
}

func (m *Model) Blur() {
	m.focused = false
}

func (m Model) GetChecked() bool {
	return m.checked
}

func (m *Model) Focus() {
	m.focused = true
}

func (m Model) Focused() bool {
	return m.focused
}

func (m *Model) SetChecked(val bool) {
	m.checked = val
}

func New() Model {
	return Model{}
}
