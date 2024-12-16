package models

import (
	"time"

	"github.com/kbavi/calendly/pkg"
	"gorm.io/gorm"
)

type EventModel struct {
	gorm.Model
	ID          string
	CalendarID  string
	Calendar    CalendarModel `gorm:"foreignKey:CalendarID"`
	Title       string
	Description *string
	Start       time.Time
	Ending      time.Time
	Invitees    string `gorm:"default:[]"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (EventModel) TableName() string {
	return "events"
}

func (e EventModel) ToEvent() *pkg.Event {
	return &pkg.Event{
		ID: e.ID,
		Calendar: &pkg.Calendar{
			ID: e.CalendarID,
		},
		Title:       e.Title,
		Description: e.Description,
		Start:       e.Start,
		End:         e.Ending,
		Invitees:    e.Invitees,
	}
}
