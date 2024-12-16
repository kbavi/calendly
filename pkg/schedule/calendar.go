package schedule

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/kbavi/calendly/pkg"
)

func (s *service) AddCalendarToUser(ctx context.Context, userID string, input *pkg.CreateCalendarInput) (*pkg.Calendar, error) {
	user, err := s.userService.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	calendar, err := s.calendarRepo.AddCalendarToUser(ctx, userID, input)
	if err != nil {
		return nil, err
	}
	calendar.User = user

	return calendar, nil
}

func (s *service) GetCalendar(ctx context.Context, calendarID string) (*pkg.Calendar, error) {
	return s.calendarRepo.GetCalendar(ctx, calendarID)
}

func (s *service) GetSchedulingLink(ctx context.Context, calendarID string, durationMinutes int) string {
	return fmt.Sprintf("/book/%s/slots/%d-mins", calendarID, durationMinutes)
}

func (s *service) GetEventSlots(ctx context.Context, input *pkg.BookBySchedulingLinkInput) ([][]time.Time, error) {
	availableIntervals, err := s.FindOverlappingAvailabilities(ctx, &pkg.FindOverlappingAvailabilitiesInput{
		CalendarIDs: []string{input.CalendarID},
		From:        input.From.Format(time.RFC3339),
		To:          input.To.Format(time.RFC3339),
	})
	if err != nil {
		return nil, err
	}

	intervalChunks := [][]time.Time{}

	for _, interval := range availableIntervals {
		if interval[0].Minute()%input.DurationMinutes != 0 {
			// set lower bound to next multiple of 30 minutes
			// e.g. if durationMinutes is 15 and interval[0] is 10:12, then lowerBound should be 10:30
			lowerBound := interval[0].Add(time.Duration(30-interval[0].Minute()%30) * time.Minute)
			// set lower bound seconds to 0
			lowerBound = time.Date(lowerBound.Year(), lowerBound.Month(), lowerBound.Day(), lowerBound.Hour(), lowerBound.Minute(), 0, 0, lowerBound.Location())

			if interval[1].Sub(lowerBound) < time.Duration(input.DurationMinutes)*time.Minute {
				continue
			}
			interval[0] = lowerBound
		}
		chunks := splitIntervalIntoMinuteChunks(interval, input.DurationMinutes)
		intervalChunks = append(intervalChunks, chunks...)
	}
	return intervalChunks, nil
}

func splitIntervalIntoMinuteChunks(interval []time.Time, durationMinutes int) [][]time.Time {
	chunks := [][]time.Time{}
	for i := 0; i < int(interval[1].Sub(interval[0]).Minutes()); i += durationMinutes {
		chunks = append(chunks, []time.Time{interval[0].Add(time.Duration(i) * time.Minute), interval[0].Add(time.Duration(i+durationMinutes) * time.Minute)})
	}
	return chunks
}

func (s *service) FindOverlappingAvailabilities(ctx context.Context, input *pkg.FindOverlappingAvailabilitiesInput) ([][]time.Time, error) {
	calendars, err := s.calendarRepo.GetCalendarsByIDs(ctx, input.CalendarIDs)
	if err != nil {
		return nil, err
	}

	if len(calendars) != len(input.CalendarIDs) {
		return nil, fmt.Errorf("some calendars not found")
	}

	from, err := time.Parse(time.RFC3339, input.From)
	if err != nil {
		return nil, err
	}

	to, err := time.Parse(time.RFC3339, input.To)
	if err != nil {
		return nil, err
	}

	// Get availability rules for all calendars
	calendarAvailabilityRules, err := s.availabilityRepo.GetAvailabilityRulesByCalendarIDs(ctx, input.CalendarIDs)
	if err != nil {
		return nil, err
	}

	// Verify all calendars have availability rules
	if err := s.validateCalendarRules(input.CalendarIDs, calendarAvailabilityRules); err != nil {
		return nil, err
	}

	dailySchedules, err := s.buildDaySchedules(calendarAvailabilityRules)
	if err != nil {
		return nil, err
	}
	// Remove days where not all calendars have availability
	expectedCalendarCount := len(input.CalendarIDs)
	for weekday, schedules := range dailySchedules {
		if len(schedules) != expectedCalendarCount {
			delete(dailySchedules, weekday)
		}
	}

	dailyOverlaps := s.findDailyOverlaps(dailySchedules)

	events, err := s.eventRepo.GetEventsByCalendarIDsInRange(ctx, input.CalendarIDs, from, to)
	if err != nil {
		return nil, err
	}

	freeIntervals, err := s.findFreeIntervals(dailyOverlaps, events, from, to)
	if err != nil {
		return nil, err
	}

	return freeIntervals, nil
}

func (s *service) validateCalendarRules(calendarIDs []string, rules map[string][]pkg.AvailabilityRule) error {
	for _, calendarID := range calendarIDs {
		if _, ok := rules[calendarID]; !ok {
			return fmt.Errorf("availability rules for calendar %s not found", calendarID)
		}
	}
	return nil
}

func (s *service) buildDaySchedules(calendarRules map[string][]pkg.AvailabilityRule) (map[time.Weekday][][][]int, error) {
	daySchedules := make(map[time.Weekday][][][]int)

	for _, rules := range calendarRules {
		for _, rule := range rules {
			if rule.Type != pkg.AvailabilityRuleTypeDay {
				continue
			}
			calendarIntervals := [][]int{}
			for _, interval := range rule.DayRule.Intervals {
				start, end, err := parseTimeInterval(interval)
				if err != nil {
					return nil, err
				}
				calendarIntervals = append(calendarIntervals, []int{start, end})
			}
			if _, ok := daySchedules[rule.DayRule.Day]; !ok {
				daySchedules[rule.DayRule.Day] = [][][]int{}
			}
			daySchedules[rule.DayRule.Day] = append(daySchedules[rule.DayRule.Day], calendarIntervals)
		}
	}

	return daySchedules, nil
}

func parseTimeInterval(interval pkg.AvailabilityInterval) (start int, end int, err error) {
	startTime, err := time.Parse("15:04", interval.From)
	if err != nil {
		return 0, 0, err
	}

	endTime, err := time.Parse("15:04", interval.To)
	if err != nil {
		return 0, 0, err
	}

	return timeToMinutes(startTime), timeToMinutes(endTime), nil
}

func timeToMinutes(t time.Time) int {
	return t.Hour()*60 + t.Minute()
}

func (s *service) findDailyOverlaps(daySchedules map[time.Weekday][][][]int) map[time.Weekday][][]int {
	result := make(map[time.Weekday][][]int)

	for day, calendars := range daySchedules {
		if len(calendars) == 0 {
			continue
		}

		// Start with first calendar's intervals
		dayResult := calendars[0]

		// Intersect with each subsequent calendar
		for i := 1; i < len(calendars); i++ {
			var newResult [][]int
			for _, r1 := range dayResult {
				for _, r2 := range calendars[i] {
					// Find overlap between intervals
					start := max(r1[0], r2[0])
					end := min(r1[1], r2[1])

					// If there is an overlap, add it to results
					if start < end {
						newResult = append(newResult, []int{start, end})
					}
				}
			}
			dayResult = newResult
		}

		// Merge overlapping intervals in final result
		if len(dayResult) == 0 {
			continue
		}

		// Sort intervals by start time
		sort.Slice(dayResult, func(i, j int) bool {
			return dayResult[i][0] < dayResult[j][0]
		})

		merged := [][]int{dayResult[0]}
		for i := 1; i < len(dayResult); i++ {
			last := merged[len(merged)-1]
			curr := dayResult[i]

			if curr[0] <= last[1] {
				// Merge overlapping intervals
				last[1] = max(last[1], curr[1])
			} else {
				// Add non-overlapping interval
				merged = append(merged, curr)
			}
		}

		result[day] = merged
	}

	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *service) findFreeIntervals(
	overlaps map[time.Weekday][][]int,
	events []pkg.Event,
	from time.Time,
	to time.Time,
) ([][]time.Time, error) {
	result := [][]time.Time{}

	// Process each day's overlaps
	for weekday, intervals := range overlaps {
		dayIntervals := [][2]time.Time{}

		// Convert minute intervals to time.Time intervals for each day in the range
		current := from
		for current.Before(to) {
			if current.Weekday() == weekday {
				for _, interval := range intervals {
					// Create time.Time for interval start and end
					start := time.Date(
						current.Year(), current.Month(), current.Day(),
						interval[0]/60, interval[0]%60, 0, 0,
						current.Location(),
					)
					end := time.Date(
						current.Year(), current.Month(), current.Day(),
						interval[1]/60, interval[1]%60, 0, 0,
						current.Location(),
					)

					// Only include intervals that overlap with the requested range
					if end.After(from) && start.Before(to) {
						// Adjust interval boundaries to respect the requested range
						if start.Before(from) {
							start = from
						}
						if end.After(to) {
							end = to
						}
						dayIntervals = append(dayIntervals, [2]time.Time{start, end})
					}
				}
			}
			current = current.AddDate(0, 0, 1)
		}

		// Split intervals based on events
		for _, event := range events {
			var newIntervals [][2]time.Time
			for _, interval := range dayIntervals {
				if event.Start.Before(interval[1]) && event.End.After(interval[0]) {
					// Event overlaps with interval
					if interval[0].Before(event.Start) {
						newIntervals = append(newIntervals, [2]time.Time{interval[0], event.Start})
					}
					if event.End.Before(interval[1]) {
						newIntervals = append(newIntervals, [2]time.Time{event.End, interval[1]})
					}
				} else {
					newIntervals = append(newIntervals, interval)
				}
			}
			dayIntervals = newIntervals
		}

		if len(dayIntervals) > 0 {
			// Convert [2]time.Time to []time.Time for the result
			intervals := make([][]time.Time, len(dayIntervals))
			for i, interval := range dayIntervals {
				intervals[i] = []time.Time{interval[0], interval[1]}
			}
			result = append(result, intervals...)
		}
	}

	return result, nil
}
