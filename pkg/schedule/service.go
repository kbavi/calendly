package schedule

import (
	"context"
	"time"

	"github.com/kbavi/calendly/pkg"
	"github.com/kbavi/calendly/pkg/user"
	"github.com/kbavi/calendly/repo"
)

type ScheduleService interface {
	AddCalendarToUser(ctx context.Context, userID string, input *pkg.CreateCalendarInput) (*pkg.Calendar, error)
	GetCalendar(ctx context.Context, calendarID string) (*pkg.Calendar, error)
	FindOverlappingAvailabilities(ctx context.Context, input *pkg.FindOverlappingAvailabilitiesInput) ([][]time.Time, error)
	GetSchedulingLink(ctx context.Context, calendarID string, durationMinutes int) string
	GetEventSlots(ctx context.Context, input *pkg.BookBySchedulingLinkInput) ([][]time.Time, error)

	SetAvailability(ctx context.Context, input *pkg.SetAvailabilityInput) (*pkg.Availability, error)
	GetAvailability(ctx context.Context, calendarID string) (*pkg.Availability, error)

	CreateEvent(ctx context.Context, input *pkg.CreateEventInput) (*pkg.Event, error)
	GetEvent(ctx context.Context, eventID string) (*pkg.Event, error)
	GetEventsByCalendarIDsInRange(ctx context.Context, calendarIDs []string, start time.Time, end time.Time) ([]pkg.Event, error)
}

type service struct {
	userService      user.Service
	calendarRepo     repo.CalendarRepository
	availabilityRepo repo.AvailabilityRepository
	eventRepo        repo.EventRepository
}

func NewService(userService user.Service, calendarRepo repo.CalendarRepository, availabilityRepo repo.AvailabilityRepository, eventRepo repo.EventRepository) ScheduleService {
	return &service{userService: userService, calendarRepo: calendarRepo, availabilityRepo: availabilityRepo, eventRepo: eventRepo}
}
