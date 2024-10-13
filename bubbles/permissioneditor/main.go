package permissioneditor

import (
	"compass/bubbles/checkbox"
	"compass/scope"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	DisplayName string
	Name        string

	All    checkbox.Model
	Access checkbox.Model
	Create checkbox.Model
	Delete checkbox.Model
	Get    checkbox.Model
	List   checkbox.Model
	Modify checkbox.Model

	cursor int
}

var checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

func (m Model) View() string {
	doc := strings.Builder{}
	doc.WriteString(m.All.View() + "\n")
	doc.WriteString(m.Access.View() + "\n")
	doc.WriteString(m.Create.View() + "\n")
	doc.WriteString(m.Delete.View() + "\n")
	doc.WriteString(m.Get.View() + "\n")
	doc.WriteString(m.List.View() + "\n")
	doc.WriteString(m.Modify.View() + "\n")
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.DisplayName+"\n",
		"----------",
		doc.String(),
	)
}

const NUM_OF_OPTIONS = 7

func BlurAll(m Model) Model {
	m.All.Blur()
	m.Access.Blur()
	m.Create.Blur()
	m.Delete.Blur()
	m.Get.Blur()
	m.List.Blur()
	m.Modify.Blur()
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, sendPermission(m)
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "j", "down":
			if m.cursor < NUM_OF_OPTIONS {
				m.cursor++
			}
		}
	}

	if m.All.Focused() {
		var cmd tea.Cmd
		m.All, cmd = m.All.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.Access.Focused() {
		var cmd tea.Cmd
		m.Access, cmd = m.Access.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.Create.Focused() {
		var cmd tea.Cmd
		m.Create, cmd = m.Create.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.Delete.Focused() {
		var cmd tea.Cmd
		m.Delete, cmd = m.Delete.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.Get.Focused() {
		var cmd tea.Cmd
		m.Get, cmd = m.Get.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.List.Focused() {
		var cmd tea.Cmd
		m.List, cmd = m.List.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.Modify.Focused() {
		var cmd tea.Cmd
		m.Modify, cmd = m.Modify.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.cursor > 0 {
		m = BlurAll(m)
	}

	switch m.cursor {
	case 0:
		m = BlurAll(m)
	case 1:
		m.All.Focus()
	case 2:
		m.Access.Focus()
	case 3:
		m.Create.Focus()
	case 4:
		m.Delete.Focus()
	case 5:
		m.Get.Focus()
	case 6:
		m.List.Focus()
	case 7:
		m.Modify.Focus()
	}

	return m, tea.Batch(cmds...)
}

func New(name string, displayName string) Model {
	all := checkbox.New()
	all.Label = "All"

	access := checkbox.New()
	access.Label = "Access"

	create := checkbox.New()
	create.Label = "Create"

	deletec := checkbox.New()
	deletec.Label = "Delete"

	get := checkbox.New()
	get.Label = "Get"

	list := checkbox.New()
	list.Label = "List"

	modify := checkbox.New()
	modify.Label = "Modify"

	return Model{
		Name:        name,
		DisplayName: displayName,
		cursor:      0,
		All:         all,
		Access:      access,
		Create:      create,
		Delete:      deletec,
		Get:         get,
		List:        list,
		Modify:      modify,
	}
}

type PermissionMessage struct {
	Flag byte
	Name string
}

func sendPermission(m Model) tea.Cmd {
	flag := byte(0)
	if m.All.GetChecked() {
		flag |= scope.S_ALL
	}
	if m.Access.GetChecked() {
		flag |= scope.S_ACCESS
	}
	if m.Create.GetChecked() {
		flag |= scope.S_CREATE
	}
	if m.Delete.GetChecked() {
		flag |= scope.S_DELETE
	}
	if m.Get.GetChecked() {
		flag |= scope.S_GET
	}
	if m.List.GetChecked() {
		flag |= scope.S_LIST
	}
	if m.Modify.GetChecked() {
		flag |= scope.S_MODIFY
	}
	return func() tea.Msg {
		return PermissionMessage{
			Name: m.Name,
			Flag: flag,
		}
	}
}
