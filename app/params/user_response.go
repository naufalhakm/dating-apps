package params

import (
	"time"
)

type UserResponse struct {
	ID          uint                 `json:"user_id"`
	Name        string               `json:"name"`
	Age         int                  `json:"age"`
	Gender      string               `json:"gender"`
	Latitude    float64              `json:"latitude"`
	Longitude   float64              `json:"longitude"`
	Interests   []string             `json:"interests"`
	Preferences *PreferencesResponse `json:"preferences"`
	LastActive  time.Time            `json:"last_active"`
}

type PreferencesResponse struct {
	PreferredAgeRange [2]int `json:"age"`
	PreferredGender   string `json:"gender"`
	MaxDistanceKm     int    `json:"distance"`
}

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"per_page"`
	PageCount  int `json:"page_count"`
	TotalCount int `json:"total_count"`
}
