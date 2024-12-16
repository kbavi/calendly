package schedule

import "time"

// TimeInterval represents a time period with start and end times
type TimeInterval struct {
	Start time.Time
	End   time.Time
}

// MinuteInterval represents a time period in minutes from start of day
type MinuteInterval struct {
	Start int
	End   int
}

// DailySchedule represents available time slots for a specific weekday
type DailySchedule struct {
	Intervals []MinuteInterval
}

// CalendarSchedule maps weekdays to their available time slots
type CalendarSchedule map[time.Weekday][]MinuteInterval
