package main

import (
	"fmt"
	"time"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"datapaddock.lan/cli_client/client"
	"datapaddock.lan/cli_client/models"
)


type homepage struct {
	id     string
	height int
	width  int

	active string
	items []string

	client *client.Client
}

// A message used to indicate that activity has occurred. 
type responseMsg struct{}

func listenForActivity(sub chan struct{}) tea.Cmd {
	client := &client.Client {
		Url: "192.168.0.151",
		Port: "8080",
	}
	return func() tea.Msg {
		for {
			syncTime := client.GetSyncInt()
			if syncTime == 0 {
				time.Sleep(time.Second * time.Duration(3))
			} else {
				time.Sleep(time.Second * time.Duration(syncTime))
				sub <- struct{}{}
			}
		}
	}
}

// A command that waits for the activity on a channel.
func waitForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}


func formatItem(device models.Device, measurement models.Measurement) string {
	retstr := fmt.Sprintf( "%s\nMAC: %s\nTemp: %f\nHum: %f\ncreatedat:%s", device.Nickname, device.MAC, 
		measurement.Temp, measurement.Humidity, measurement.CreatedAt)
	return retstr
}



func (h *homepage) populateItems() {
	devs := h.client.GetDevices()
	meas := h.client.GetLast()

	devmap := make(map[string]models.Device)

	for _, dev := range devs {
		devmap[dev.MAC] = dev
	}
	
	var new_items []string
	
	for _, measurement := range meas {
		new_items = append(new_items, formatItem(devmap[measurement.MAC], measurement))
	}
	h.items = nil
	h.items = new_items
	return
}


func (h homepage) Init() tea.Cmd {
	return nil
}

func (h *homepage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	h.populateItems()
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.height = msg.Height
		h.width = msg.Width
	case responseMsg:
		return h, nil
	case tea.MouseMsg:
		if msg.Type != tea.MouseLeft {
			return h, nil
		}

		for _, item := range h.items {
			// Check each item to see if it's in bounds.
			if zone.Get(h.id + item).InBounds(msg) {
				h.active = item
				break
			}
		}
	}
	return h, nil
}

func (h homepage) View() string {
	homepageStyle := lipgloss.NewStyle().
	Align(lipgloss.Center).
	Foreground(lipgloss.Color("#c0caf5")).
	//Background(lipgloss.Color("#1a1b26")).
	Background(subtle).
	Margin(1).
	Padding(1,2).
	Width((h.width/ len(h.items)) -2).
	Height(h.height - 2).
	MaxHeight(h.height)

	out := []string{}
	for _, item := range h.items{
		out = append(out, zone.Mark(h.id+item, homepageStyle.Render(item)))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, out...)
}
