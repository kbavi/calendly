package models

import (
	"encoding/json"

	"github.com/kbavi/calendly/pkg"
	"gorm.io/gorm"
)

type AvailabilityModel struct {
	gorm.Model
	ID         string
	CalendarID string
	Rules      string
}

func (AvailabilityModel) TableName() string {
	return "availabilities"
}

func (a AvailabilityModel) ToAvailability() *pkg.Availability {
	var rules []pkg.AvailabilityRule
	err := json.Unmarshal([]byte(a.Rules), &rules)
	if err != nil {
		return nil
	}
	return &pkg.Availability{
		ID: a.ID,
		Calendar: pkg.Calendar{
			ID: a.CalendarID,
		},
		Rules: rules,
	}
}
