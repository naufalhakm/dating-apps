package models

import (
	"time"
)

type User struct {
	ID          uint
	Name        string
	Age         int
	Gender      string
	Latitude    float64
	Longitude   float64
	Interests   []string
	Preferences *Preferences
	LastActive  time.Time
}

type Preferences struct {
	PreferredAgeRange [2]int `json:"age"`
	PreferredGender   string `json:"gender"`
	MaxDistanceKm     int    `json:"distance"`
}
