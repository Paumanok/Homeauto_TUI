package main

import (
	"fmt"
	"time"
	//"math/rand"
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

// A message used to indicate that activity has occurred. In the real world (for
// example, chat) this would contain actual data.
type responseMsg struct{}

// Simulate a process that sends events at an irregular interval in real time.
// In this case, we'll send events on the channel at a random interval between
// 100 to 1000 milliseconds. As a command, Bubble Tea will run this
// asynchronously.
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


//func waitForUpdate(h *homepage) tea.Cmd {
//	return func() tea.Msg {
//		for {
//			syncTime := h.client.GetSyncInt()
//			fmt.Println("starting sleep")
//			time.Sleep(time.Second * time.Duration(syncTime))
//			h.populateItems()
//			fmt.Println("populated")
//			h.sub <- "new"
//		}
//	}
//}
//
//func respondToUpdate(h *homepage) tea.Cmd {
//	return func() tea.Msg {
//		return <-h.sub
//	}
//}

//trying to do this here https://github.com/charmbracelet/bubbletea/blob/master/examples/realtime/main.go
func (h *homepage) populateItems() {
	devs := h.client.GetDevices()
	meas := h.client.GetLast()

	//START FROM HERE
	//we want to formwat a string 
	//that can be passed as an "item"
	//using the above data. 
	//perhaps mapping devs to MAC and
	//using the last measurements' MAC
	//to sort things would be a good idea. 
	//fmt.Println("populatin")	
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
	//return tea.Batch (
	//	listenForActivity(h.Sub),
	//	waitForActivity(h.sub),
	//)
	return nil
}

func (h *homepage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//fmt.Println(h.items)
	//fmt.Println("##### populating ###")
	h.populateItems()
	//fmt.Println(h.items)
	//fmt.Println("update")
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.height = msg.Height
		h.width = msg.Width
	case responseMsg:
		//fmt.Println("realtime")                  // record external activity
		return h, nil //waitForActivity(h.Sub) // wait for next event
	case tea.MouseMsg:
		if msg.Type != tea.MouseLeft {
			//fmt.Println("update2")
			return h, nil
		}

		for _, item := range h.items {
			// Check each item to see if it's in bounds.
			if zone.Get(h.id + item).InBounds(msg) {
				h.active = item
				break
			}
		}
//	case string:
//		if msg == "new" {
//			return h, respondToUpdate(&h)	
//		}
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
