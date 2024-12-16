package pkg

import (
	"errors"
	"strings"
	"time"
)

type CreateUserInput struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required"`
}

func (i *CreateUserInput) Validate() []error {
	var errs []error

	if i.Email == "" {
		errs = append(errs, errors.New("email is required"))
	}
	if i.Name == "" {
		errs = append(errs, errors.New("name is required"))
	}

	return errs
}

type UserDTO struct {
	ID    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}

func (u *UserDTO) FromUser(user *User) {
	u.ID = user.ID
	u.Email = string(user.Email)
	u.Name = user.Name
}

type ReturnUserResponse struct {
	Status string             `json:"status"`
	Data   map[string]UserDTO `json:"data"`
}

type FindOverlappingAvailabilitiesInput struct {
	CalendarIDs []string `json:"calendar_ids" binding:"required"`
	From        string   `json:"from" binding:"required"`
	To          string   `json:"to" binding:"required"`
}

func (i *FindOverlappingAvailabilitiesInput) Validate() []error {
	var errs []error

	if len(i.CalendarIDs) == 0 {
		errs = append(errs, errors.New("calendar_ids is required"))
	}

	if i.From == "" {
		errs = append(errs, errors.New("from is required"))
	}

	if i.To == "" {
		errs = append(errs, errors.New("to is required"))
	}

	fromTime, err := time.Parse(time.RFC3339, i.From)
	if err != nil {
		errs = append(errs, errors.New("invalid from time format, should be RFC3339"))
	}

	toTime, err := time.Parse(time.RFC3339, i.To)
	if err != nil {
		errs = append(errs, errors.New("invalid to time format, should be RFC3339"))
	}

	if !fromTime.Before(toTime) {
		errs = append(errs, errors.New("from time must be before to time"))
	}

	return errs
}

type BookBySchedulingLinkInput struct {
	DurationMinutes int       `json:"duration_minutes" binding:"required"`
	From            time.Time `json:"from" binding:"required"`
	To              time.Time `json:"to" binding:"required"`
	CalendarID      string    `json:"calendar_id" binding:"required"`
}

type CreateCalendarInput struct {
	UserID string `json:"user_id" binding:"required"`
}

func (i *CreateCalendarInput) Validate() []error {
	var errs []error

	if i.UserID == "" {
		errs = append(errs, errors.New("user_id is required"))
	}

	return errs
}

type CalendarDTO struct {
	ID   string  `json:"id"`
	User UserDTO `json:"user,omitempty"`
}

func (c *CalendarDTO) FromCalendar(calendar *Calendar) {
	c.ID = calendar.ID
	if calendar.User != nil {
		c.User.FromUser(calendar.User)
	}
}

type ReturnCalendarResponse struct {
	Status string                 `json:"status"`
	Data   map[string]CalendarDTO `json:"data"`
}

type AvailabilityDTO struct {
	ID       string             `json:"id"`
	User     UserDTO            `json:"user"`
	Calendar CalendarDTO        `json:"calendar"`
	Rules    []AvailabilityRule `json:"rules"`
}

type SetAvailabilityInput struct {
	CalendarID string             `json:"calendar_id" binding:"required"`
	Rules      []AvailabilityRule `json:"rules" binding:"required"`
}

func (i *SetAvailabilityInput) Validate() []error {
	var errs []error

	if i.CalendarID == "" {
		errs = append(errs, errors.New("calendar_id is required"))
	}

	if len(i.Rules) == 0 {
		errs = append(errs, errors.New("rules are required"))
	}

	for _, rule := range i.Rules {
		if rule.Type == "" {
			errs = append(errs, errors.New("rule type is required"))
		}

		if rule.Type == AvailabilityRuleTypeDay && rule.DayRule == nil {
			errs = append(errs, errors.New("day rule is required"))
		}

		if rule.Type == AvailabilityRuleTypeDate && rule.DateRule == nil {
			errs = append(errs, errors.New("date rule is required"))
		}
	}

	return errs
}

type ReturnAvailabilityResponse struct {
	Status string                     `json:"status"`
	Data   map[string]AvailabilityDTO `json:"data"`
}

type CreateEventInput struct {
	UserID      string  `json:"user_id"`
	CalendarID  string  `json:"calendar_id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Invitees    string  `json:"invitees"`
	Start       string  `json:"start"`
	End         string  `json:"end"`
}

type EventDTO struct {
	ID          string      `json:"id"`
	Calendar    CalendarDTO `json:"calendar"`
	Title       string      `json:"title"`
	Description *string     `json:"description"`
	Invitees    string      `json:"invitees"`
	Start       time.Time   `json:"start"`
	End         time.Time   `json:"end"`
}

func (e *EventDTO) FromEvent(event *Event) {
	e.ID = event.ID
	if event.Calendar != nil {
		e.Calendar.FromCalendar(event.Calendar)
	}
	e.Title = event.Title
	e.Description = event.Description
	e.Invitees = event.Invitees
	e.Start = event.Start
	e.End = event.End
}

func (i *CreateEventInput) Validate() []string {
	var errs []string

	if i.CalendarID == "" {
		errs = append(errs, "calendar_id is required")
	}

	if i.UserID == "" {
		errs = append(errs, "user_id is required")
	}

	if i.Title == "" {
		errs = append(errs, "title is required")
	}

	fromTime, err := time.Parse(time.RFC3339, i.Start)
	if err != nil {
		errs = append(errs, "invalid from time format, should be RFC3339")
	}

	toTime, err := time.Parse(time.RFC3339, i.End)
	if err != nil {
		errs = append(errs, "invalid to time format, should be RFC3339")
	}

	if !fromTime.Before(toTime) {
		errs = append(errs, "from time must be before to time")
	}

	if i.Invitees == "" {
		errs = append(errs, "invitees are required")
	}

	invitees := strings.Split(i.Invitees, ",")
	for _, invitee := range invitees {
		if invitee == "" {
			errs = append(errs, "invitees must be a comma separated list of emails")
		}
	}

	return errs
}

type ReturnEventResponse struct {
	Status string              `json:"status"`
	Data   map[string]EventDTO `json:"data"`
}

type ReturnEventsResponse struct {
	Status string                `json:"status"`
	Data   map[string][]EventDTO `json:"data"`
}
