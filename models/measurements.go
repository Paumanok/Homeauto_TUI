package models

import (
	"time"
)

type Measurement struct {
	MAC string `json:"MAC"`
	Temp float32 `json:"temp"`
	Humidity float32 `json:"humidity"`
	Pressure float32 `json:"pressure,omitempty"`
	PM25	float32  `json:"pm25,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}
