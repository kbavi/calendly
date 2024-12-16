package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/kbavi/calendly/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	os.Setenv("TZ", "UTC")
}

func TestCalendarIntegration(t *testing.T) {
	ctx := context.Background()

	// Create a test user first
	user, err := testApp.Services.UserService.Create(ctx, &pkg.CreateUserInput{
		Name:  "Calendar Test User",
		Email: "calendar_test@example.com",
	})
	require.NoError(t, err)
	require.NotNil(t, user)

	// Clean up user after all tests
	defer testApp.Services.UserService.Delete(ctx, user.ID)

	calendar, err := testApp.Services.ScheduleService.AddCalendarToUser(ctx, user.ID, &pkg.CreateCalendarInput{
		UserID: user.ID,
	})
	require.NoError(t, err)
	require.NotNil(t, calendar)

	t.Run("Availability_Service_Tests", func(t *testing.T) {
		t.Run("Set_And_Get_Availability", func(t *testing.T) {
			createInput := &pkg.SetAvailabilityInput{
				CalendarID: calendar.ID,
				Rules: []pkg.AvailabilityRule{
					{
						Type: pkg.AvailabilityRuleTypeDay,
						DayRule: &pkg.AvailabilityDayRule{
							Day: 0,
							Intervals: []pkg.AvailabilityInterval{
								{
									From: "09:00",
									To:   "17:00",
								},
							},
						},
					},
				},
			}

			createdAvailability, err := testApp.Services.ScheduleService.SetAvailability(ctx, createInput)
			require.NoError(t, err)
			require.NotNil(t, createdAvailability)

			fetchedAvailability, err := testApp.Services.ScheduleService.GetAvailability(ctx, calendar.ID)
			require.NoError(t, err)
			require.NotEmpty(t, fetchedAvailability)
		})
	})

	t.Run("Event_Service_Tests", func(t *testing.T) {
		// First set up availability for Monday
		availabilityInput := &pkg.SetAvailabilityInput{
			CalendarID: calendar.ID,
			Rules: []pkg.AvailabilityRule{
				{
					Type: pkg.AvailabilityRuleTypeDay,
					DayRule: &pkg.AvailabilityDayRule{
						Day: 1, // Monday
						Intervals: []pkg.AvailabilityInterval{
							{
								From: "09:00",
								To:   "17:00",
							},
						},
					},
				},
			},
		}

		_, err := testApp.Services.ScheduleService.SetAvailability(ctx, availabilityInput)
		require.NoError(t, err)

		t.Run("Create_And_Get_Event", func(t *testing.T) {
			// Find next Monday
			now := time.Now()
			daysUntilMonday := (8 - int(now.Weekday())) % 7
			if daysUntilMonday == 0 {
				daysUntilMonday = 7
			}
			description := "Integration Test Meeting"

			createEventInput := &pkg.CreateEventInput{
				UserID:      user.ID,
				CalendarID:  calendar.ID,
				Title:       "Test Meeting",
				Description: &description,
				Start:       "2024-12-14T14:30:00Z",
				End:         "2024-12-14T14:30:00Z",
				Invitees:    "attendee@example.com",
			}

			createdEvent, err := testApp.Services.ScheduleService.CreateEvent(ctx, createEventInput)
			require.NoError(t, err)
			require.NotNil(t, createdEvent)
			assert.Equal(t, createEventInput.Title, createdEvent.Title)
			assert.Equal(t, createEventInput.Description, createdEvent.Description)

			// Test getting the event
			fetchedEvent, err := testApp.Services.ScheduleService.GetEvent(ctx, createdEvent.ID)
			require.NoError(t, err)
			require.NotNil(t, fetchedEvent)
			assert.Equal(t, createdEvent.ID, fetchedEvent.ID)

			// // Test getting calendar events
			// calendarEvents, err := testApp.Services.ScheduleService.GetCalendarEvents(ctx, calendar.ID)
			// require.NoError(t, err)
			// require.NotEmpty(t, calendarEvents)
			// assert.Contains(t, calendarEvents, createdEvent)

			// // Test canceling the event
			// err = testApp.Services.ScheduleService.CancelEvent(ctx, user.ID, createdEvent.ID)
			// require.NoError(t, err)

			// // Verify event is canceled
			// canceledEvent, err := testApp.Services.ScheduleService.GetEvent(ctx, createdEvent.ID)
			// require.NoError(t, err)
			// assert.Equal(t, pkg.EventStatusCanceled, canceledEvent.Status)
		})
	})
}
