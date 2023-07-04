package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"datapaddock.lan/cli_client/client"
)

// This is a modified version of this example, supporting full screen, dynamic
// resizing, and clickable models (tabs, lists, dialogs, etc).
// 	https://github.com/charmbracelet/lipgloss/blob/master/example

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(subtle).
		String()
)

type model struct {
	height int
	width  int
	tabs tea.Model
	home tea.Model
	sub chan struct{}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(listenForActivity(m.sub), waitForActivity(m.sub))
}

func (m model) isInitialized() bool {
	return m.height != 0 && m.width != 0
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.isInitialized() {
		if _, ok := msg.(tea.WindowSizeMsg); !ok {
			return m, nil
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Example of toggling mouse event tracking on/off.
		if msg.String() == "ctrl+e" {
			zone.SetEnabled(!zone.Enabled())
			return m, nil
		}

		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		msg.Height -= 2
		msg.Width -= 4
		return m.propagate(msg), nil
	case responseMsg:
		return m.propagate(msg), waitForActivity(m.sub) //waitFor
	}

	return m.propagate(msg), nil
}

func (m *model) propagate(msg tea.Msg) tea.Model {
	m.tabs, _ = m.tabs.Update(msg)
	m.home, _ = m.home.Update(msg)

	//if msg, ok := msg.(tea.WindowSizeMsg); ok {
	//	msg.Height -= m.tabs.(tabs).height + m.list1.(list).height
	//	return m 
	//}
	return m
}

func (m model) View() string {
	if !m.isInitialized() {
		return ""
	}
	s := lipgloss.NewStyle().MaxHeight(m.height).MaxWidth(m.width).Padding(1,2,1,2)

	switch m.tabs.(tabs).active {
		case "Home":
			return zone.Scan(s.Render(lipgloss.JoinVertical(lipgloss.Top,
				m.tabs.View(), "", m.home.View(), )))
	default:
		return zone.Scan(s.Render(lipgloss.JoinVertical(lipgloss.Top, m.tabs.View(), "", )))
	}

	return zone.Scan(s.Render(lipgloss.JoinVertical(lipgloss.Top, m.tabs.View(), "", )))
}

func main() {

	client := &client.Client {
		Url: "192.168.0.151",
		Port: "8080",
	}

	fmt.Println(client.GetDevices())
	fmt.Println(client.GetLast())
	fmt.Println(client.GetSyncInt())

	zone.NewGlobal()

	m := &model{
		tabs: &tabs{
			id:     zone.NewPrefix(), // Give each type an ID, so no zones will conflict.
			height: 3,
			active: "Home",
			items:  []string{"Home", "Particulate Matter", "Remind Me", "Config"},
		},
		home: &homepage{
			id: zone.NewPrefix(),
			client: client,
			items: []string{
				"test12",
			},
		},
		sub: make(chan struct{}),
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _,err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
