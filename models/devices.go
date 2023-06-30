package models


type Device struct {
	Nickname string	`json:"nickname"`
	MAC	string `json:"MAC"`
	HumidityComp int `json:"humidityComp"`
	TemperatureComp int `json:"temperatureComp"`
}

type Devices []Device

