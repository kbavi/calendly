package pkg

import (
	"time"
)

type Email string

type User struct {
	ID    string
	Email Email
	Name  string
}

type Calendar struct {
	ID           string
	User         *User
	Availability *Availability
}

type AvailabilityInterval struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type AvailabilityDayRule struct {
	Day       time.Weekday           `json:"day"`
	Intervals []AvailabilityInterval `json:"intervals"`
}

type AvailabilityDateRule struct {
	Date      time.Time              `json:"date"`
	Intervals []AvailabilityInterval `json:"intervals"`
}

type AvailabilityRuleType string

const (
	AvailabilityRuleTypeDay  AvailabilityRuleType = "day"
	AvailabilityRuleTypeDate AvailabilityRuleType = "date"
)

type AvailabilityRule struct {
	Type     AvailabilityRuleType  `json:"type"`
	DayRule  *AvailabilityDayRule  `json:"day,omitempty"`
	DateRule *AvailabilityDateRule `json:"date,omitempty"`
}

type Availability struct {
	ID       string
	Calendar Calendar
	Rules    []AvailabilityRule
}

type Event struct {
	ID          string
	Calendar    *Calendar
	Title       string
	Description *string
	Invitees    string
	Start       time.Time
	End         time.Time
	Status      EventStatus
}

type EventStatus string

const (
	EventStatusBooked    EventStatus = "booked"
	EventStatusCancelled EventStatus = "cancelled"
)
