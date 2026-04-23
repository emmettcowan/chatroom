package client

import (
	"bufio"
	"fmt"
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
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Say hello"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(20)

	return model{textInput: ti}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
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
			fmt.Println(m.textInput.Value())

		}
	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	if m.w == 0 {
		return tea.NewView("")
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(1)

	top := box.Width(m.w).Render("top\ncontent")
	bottom := box.Width(m.w).Render(m.textInput.View())

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

	go listenToServer(conn)
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	// reader := bufio.NewReader(os.Stdin)
	// for {
	// 	fmt.Print(">: ")
	// 	text, err := reader.ReadString('\n')
	// 	if err != nil {
	// 		if err != os.ErrClosed {
	// 			log.Printf("Error reading stdin: %v", err)
	// 		}
	// 		break
	// 	}
	//
	// 	_, err = conn.Write([]byte(text))
	// 	if err != nil {
	// 		fmt.Printf("Error writing to connection: %v\n", err)
	// 		break
	// 	}
	// }
}

func listenToServer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\nDisconnected from server: %v\n", err)
			return
		}
		fmt.Printf("\rServer: %s>: ", response)
	}
}
