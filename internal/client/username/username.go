package username

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type UsernameModel struct {
	textInput textinput.Model
	Username  string
	err       error
	quitting  bool
}

func UsernameScreen() UsernameModel {
	ti := textinput.New()
	ti.Placeholder = "Enter a username"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(20)

	return UsernameModel{textInput: ti, Username: ""}
}

func (m UsernameModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m UsernameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			m.Username = m.textInput.Value()
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m UsernameModel) View() tea.View {
	var c *tea.Cursor
	if !m.textInput.VirtualCursor() {
		c = m.textInput.Cursor()
		c.Y += lipgloss.Height(m.headerView())
	}

	str := lipgloss.JoinVertical(lipgloss.Top, m.headerView(), m.textInput.View(), m.footerView())
	if m.quitting {
		str += "\n"
	}

	v := tea.NewView(str)
	v.Cursor = c
	return v
}

func (m UsernameModel) headerView() string { return "Whats your username?\n" }
func (m UsernameModel) footerView() string { return "\n(esc to quit)" }

