package zoneselector

import (
	"compass/bubbles/checkbox"
	"compass/bubbles/permissioneditor"
	"compass/scope"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jowiklund/goutil/queue"
)

type ZoneSelect struct {
	zone  scope.ZoneData
	input checkbox.Model
}

type PermissionCollection map[string]scope.Permission

type Model struct {
	zones            []ZoneSelect
	zoneQueue        queue.Queue[scope.ZoneData]
	permissions      PermissionCollection
	permissionEditor permissioneditor.Model
	selected         int
	editMode         bool
	finished         bool
}

func New(zones []scope.ZoneData) Model {
	m := Model{
		zoneQueue:        queue.NewQueue[scope.ZoneData](len(zones)),
		permissionEditor: permissioneditor.New("", ""),
		permissions:      PermissionCollection{},
		finished:         false,
	}
	for _, zone := range zones {
		cb := checkbox.New()
		cb.Label = zone.Name
		m.zones = append(m.zones, ZoneSelect{
			zone:  zone,
			input: cb,
		})
	}
	return m
}

func permissionZoneName(z scope.ZoneData) string {
	return "Zone" + z.Id.String()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case permissioneditor.PermissionMessage:
		m.permissions[msg.Name] = scope.CreatePermission(msg.Flag)
		m.editMode = false
		if m.finished {
			return m, func() tea.Msg {
				return m.permissions
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			for index, input := range m.zones {
				if input.input.GetChecked() {
					m.zones[index].input.SetChecked(false)
					m.zoneQueue.Enqueue(input.zone)
				}
			}
		case "tab", "j", "down":
			m.selected = min(m.selected+1, len(m.zones)-1)
		case "shift+tab", "k", "up":
			m.selected = max(m.selected-1, 0)
		}
	}

	if m.editMode {
		m.permissionEditor, cmd = m.permissionEditor.Update(msg)
		return m, cmd
	}

	if !m.zoneQueue.Empty() && !m.editMode {
		item, err := m.zoneQueue.Dequeue()
		if err != nil {
			return m, func() tea.Msg {
				return err
			}
		}
		m.permissionEditor = permissioneditor.New(
			permissionZoneName(item),
			item.Name,
		)

		m.editMode = true
		if m.zoneQueue.Empty() {
			m.finished = true
		}
	}

	if !m.editMode {
		for index := range m.zones {
			m.zones[index].input.Blur()
		}
		m.zones[m.selected].input.Focus()
		m.zones[m.selected].input, cmd = m.zones[m.selected].input.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	if m.editMode {
		return m.permissionEditor.View()
	}
	if m.finished {
		doc := strings.Builder{}
		for _, item := range m.permissions {
			doc.WriteString(fmt.Sprintf("%+v\n\n", item))
		}
		return doc.String()
	}
	doc := strings.Builder{}
	for _, zone := range m.zones {
		doc.WriteString(zone.input.View())
		doc.WriteString("\n")
	}
	return lipgloss.JoinVertical(
		lipgloss.Top,
		doc.String(),
		fmt.Sprintf("You have %d permissions added", len(m.permissions)),
		fmt.Sprintf("%+v", m.permissions),
	)
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
