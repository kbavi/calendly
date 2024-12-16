package schedule

import (
	"context"

	"github.com/kbavi/calendly/pkg"
)

func (s *service) SetAvailability(ctx context.Context, input *pkg.SetAvailabilityInput) (*pkg.Availability, error) {
	calendar, err := s.calendarRepo.GetCalendar(ctx, input.CalendarID)
	if err != nil {
		return nil, err
	}

	availability, err := s.availabilityRepo.SetAvailability(ctx, input)
	if err != nil {
		return nil, err
	}
	availability.Calendar = *calendar
	return availability, nil
}

func (s *service) GetAvailability(ctx context.Context, calendarID string) (*pkg.Availability, error) {
	availability, err := s.availabilityRepo.GetAvailability(ctx, calendarID)
	if err != nil {
		return nil, err
	}
	calendar, err := s.calendarRepo.GetCalendar(ctx, calendarID)
	if err != nil {
		return nil, err
	}
	availability.Calendar = *calendar
	return availability, nil
}
