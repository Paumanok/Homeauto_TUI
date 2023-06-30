package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"datapaddock.lan/ht_client/client"
	"datapaddock.lan/ht_client/models"
)


type homepage struct {
	id     string
	height int
	width  int

	active string
	items []string

	client *client.Client
}

func formatItem(device models.Device, measurement models.Measurement) string {
	retstr := fmt.Sprintf( "%s\nMAC: %s\nTemp: %f\nHum: %f", device.Nickname, device.MAC, 
		measurement.Temp, measurement.Humidity)
	return retstr
}

//trying to do this here https://github.com/charmbracelet/bubbletea/blob/master/examples/realtime/main.go
func (h *homepage) populateItems(sub chan []string) tea.Cmd {
	devs := h.client.GetDevices()
	meas := h.client.GetLast()

	//START FROM HERE
	//we want to formwat a string 
	//that can be passed as an "item"
	//using the above data. 
	//perhaps mapping devs to MAC and
	//using the last measurements' MAC
	//to sort things would be a good idea. 
	
	devmap := make(map[string]models.Device)

	for _, dev := range devs {
		devmap[dev.MAC] = dev
	}
	
	var new_items []string
	
	for _, measurement := range meas {
		new_items = append(new_items, formatItem(devmap[measurement.MAC], measurement))
	}
	//h.items = nil
	//h.items = new_items
	return
}


func (h homepage) Init() tea.Cmd {
	return nil
}

func (h homepage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	h.populateItems()
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.height = msg.Height
		h.width = msg.Width
	case tea.MouseMsg:

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
	//fmt.Println(h.items)
	for _, item := range h.items{
		out = append(out, zone.Mark(h.id+item, homepageStyle.Render(item)))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, out...)
}
