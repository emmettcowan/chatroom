package client

import (
	"bufio"
	"log"
	"net"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
	"github.com/emmettcowan/chatroom/internal/client/username"
)

type model struct {
	textInput textinput.Model
	quitting  bool
	w, h      int
	conn      net.Conn
	ch        chan incommingMsg
	recvied   string
	username  string
}

// TODO: Use state to swap between the login and test screen at some point
// currently running one tui before the next

func initialModel(conn net.Conn, ch chan incommingMsg, username string) model {
	ti := textinput.New()
	ti.Placeholder = "Say hello"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(20)

	return model{textInput: ti, conn: conn, ch: ch, username: username}
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
				"Disconnected from server\n",
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
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			_, err := m.conn.Write([]byte(m.username + ":  " + m.textInput.Value() + "\n"))
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
		return m, reviceFromServer(m.ch)
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginLeft(1)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)

	subtleStyle = lipgloss.NewStyle().
			Foreground(compat.AdaptiveColor{
			Light: lipgloss.Color("#9B9B9B"),
			Dark:  lipgloss.Color("#5C5C5C"),
		})
)

func (m model) View() tea.View {
	if m.w == 0 {
		return tea.NewView("")
	}

	header := titleStyle.Render("Chatroom")
	usernameDisplay := subtleStyle.Render("Logged in as: ") + lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Render(m.username)

	messages := m.recvied
	if messages == "" {
		messages = subtleStyle.Render("No messages yet...")
	}
	msgBox := borderStyle.Width(m.w - 4).Height(m.h - 9).Render(messages)

	inputBox := borderStyle.Width(m.w - 4).BorderForeground(lipgloss.Color("170")).Render(m.textInput.View())

	help := subtleStyle.Render("esc: quit enter: send")

	ui := lipgloss.JoinVertical(
		lipgloss.Left,
		header+" "+usernameDisplay,
		msgBox,
		inputBox,
		help,
	)

	out := lipgloss.NewStyle().Padding(1, 2).Render(ui)

	v := tea.NewView(out)
	v.AltScreen = true
	return v
}

func Run() {
	usernameScreen := tea.NewProgram(username.UsernameScreen())
	m, err := usernameScreen.Run()
	if err != nil {
		log.Fatal(err)
	}
	usernameModel := m.(username.UsernameModel)
	conn, err := net.Dial("tcp", ":8090")
	if err != nil {
		log.Fatal("Error connecting: ", err)
	}
	defer conn.Close()

	ch := make(chan incommingMsg)

	go listenToServer(conn, ch)
	p := tea.NewProgram(initialModel(conn, ch, usernameModel.Username))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
