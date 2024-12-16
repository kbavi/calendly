package repo

import (
	"context"
	"time"

	"github.com/kbavi/calendly/pkg"
)

type UserRepository interface {
	Create(ctx context.Context, input *pkg.CreateUserInput) (*pkg.User, error)
	Get(ctx context.Context, id string) (*pkg.User, error)
	Delete(ctx context.Context, id string) error
}

type CalendarRepository interface {
	AddCalendarToUser(ctx context.Context, userID string, input *pkg.CreateCalendarInput) (*pkg.Calendar, error)
	GetCalendar(ctx context.Context, calendarID string) (*pkg.Calendar, error)
	GetCalendarsByIDs(ctx context.Context, calendarIDs []string) ([]pkg.Calendar, error)
}

type EventRepository interface {
	CreateEvent(ctx context.Context, input *pkg.CreateEventInput) (*pkg.Event, error)
	GetEvent(ctx context.Context, eventID string) (*pkg.Event, error)
	GetEventsByCalendarIDsInRange(ctx context.Context, calendarIDs []string, start time.Time, end time.Time) ([]pkg.Event, error)
}

type AvailabilityRepository interface {
	SetAvailability(ctx context.Context, input *pkg.SetAvailabilityInput) (*pkg.Availability, error)
	GetAvailability(ctx context.Context, calendarID string) (*pkg.Availability, error)
	GetAvailabilityRulesByCalendarIDs(ctx context.Context, calendarIDs []string) (map[string][]pkg.AvailabilityRule, error)
}
