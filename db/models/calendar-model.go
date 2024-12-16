package models

import "gorm.io/gorm"

type CalendarModel struct {
	gorm.Model
	ID             string
	UserID         string
	AvailabilityID *string
	Availability   *AvailabilityModel `gorm:"foreignKey:AvailabilityID"`
}

func (CalendarModel) TableName() string {
	return "calendars"
}
