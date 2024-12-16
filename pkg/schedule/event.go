package schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/kbavi/calendly/pkg"
)

func (s *service) CreateEvent(ctx context.Context, input *pkg.CreateEventInput) (*pkg.Event, error) {
	_, err := s.userService.Get(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}
	_, err = s.calendarRepo.GetCalendar(ctx, input.CalendarID)
	if err != nil {
		return nil, fmt.Errorf("invalid calendar id")
	}
	event, err := s.eventRepo.CreateEvent(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create event")
	}
	return event, nil
}

func (s *service) GetEvent(ctx context.Context, eventID string) (*pkg.Event, error) {
	event, err := s.eventRepo.GetEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *service) GetEventsByCalendarIDsInRange(ctx context.Context, calendarIDs []string, start time.Time, end time.Time) ([]pkg.Event, error) {
	events, err := s.eventRepo.GetEventsByCalendarIDsInRange(ctx, calendarIDs, start, end)
	if err != nil {
		return nil, err
	}
	return events, nil
}
