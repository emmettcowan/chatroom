package client

import (
	"bufio"
	"log"
	"net"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	textInput textinput.Model
	err       error
	quitting  bool
	w, h      int
	conn      net.Conn
	ch        chan incommingMsg
	recvied   string
}

func initialModel(conn net.Conn, ch chan incommingMsg) model {
	ti := textinput.New()
	ti.Placeholder = "Say hello"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(20)

	return model{textInput: ti, conn: conn, ch: ch}
}

type incommingMsg struct {
	message string
}

func listenToServer(conn net.Conn, ch chan incommingMsg) {
	reader := bufio.NewReader(conn)
	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			ch <- incommingMsg{
				"\nDisconnected from server\n",
			}
			return
		}
		ch <- incommingMsg{
			response,
		}
	}
}

func reviceFromServer(sub chan incommingMsg) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		reviceFromServer(m.ch),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			_, err := m.conn.Write([]byte(m.textInput.Value() + "\n"))
			if err != nil {
				log.Printf("Error writing to server: %v", err)
			}
			m.textInput.SetValue("")

		}
	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height
	case incommingMsg:
		m.recvied += msg.message
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	if m.w == 0 {
		return tea.NewView("")
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1)

	top := box.Width(m.w).Height(m.h - 5).Render(m.recvied)
	bottom := box.Width(m.w).Height(1).Render(m.textInput.View())

	out := lipgloss.JoinVertical(
		lipgloss.Top,
		top,
		bottom,
	)

	return tea.NewView(out)
}

func Run() {
	conn, err := net.Dial("tcp", ":8090")
	if err != nil {
		log.Fatal("Error connecting: ", err)
	}
	defer conn.Close()

	ch := make(chan incommingMsg)

	go listenToServer(conn, ch)
	p := tea.NewProgram(initialModel(conn, ch))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
