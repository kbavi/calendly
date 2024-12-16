package integration

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/kbavi/calendly/pkg"
	"github.com/stretchr/testify/require"
)

func init() {
	os.Setenv("TZ", "UTC")
}

type TestUser struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Availability struct {
		Rules []pkg.AvailabilityRule `json:"rules"`
	} `json:"availability"`
}

type TestEvent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Invitees    string `json:"invitees"`
}

type TestUserEvents struct {
	UserEmail string      `json:"userEmail"`
	Events    []TestEvent `json:"events"`
}

type TestData struct {
	Users []TestUser `json:"users"`
}

type EventsData struct {
	Events []TestUserEvents `json:"events"`
}

func loadUserMockData(t *testing.T) TestData {
	jsonFile, err := os.Open(filepath.Join("mockdata", "users.json"))
	require.NoError(t, err)
	defer jsonFile.Close()

	var testData TestData
	err = json.NewDecoder(jsonFile).Decode(&testData)
	require.NoError(t, err)
	return testData
}

func loadEventsMockData(t *testing.T) EventsData {
	jsonFile, err := os.Open(filepath.Join("mockdata", "events.json"))
	require.NoError(t, err)
	defer jsonFile.Close()

	var eventsData EventsData
	err = json.NewDecoder(jsonFile).Decode(&eventsData)
	require.NoError(t, err)
	return eventsData
}

func TestMultipleUserAvailability(t *testing.T) {
	ctx := context.Background()
	userMockData := loadUserMockData(t)
	eventsMockData := loadEventsMockData(t)

	// Create users and calendars
	userMap := make(map[string]*pkg.User)         // map[email]user
	calendarMap := make(map[string]*pkg.Calendar) // map[email]calendar

	for _, userData := range userMockData.Users {
		user, err := testApp.Services.UserService.Create(ctx, &pkg.CreateUserInput{
			Name:  userData.Name,
			Email: userData.Email,
		})
		require.NoError(t, err)
		require.NotNil(t, user)
		userMap[userData.Email] = user

		// Create calendar for each user
		calendar, err := testApp.Services.ScheduleService.AddCalendarToUser(ctx, user.ID, &pkg.CreateCalendarInput{
			UserID: user.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, calendar)
		calendarMap[userData.Email] = calendar
	}

	// Clean up users after tests
	defer func() {
		for _, user := range userMap {
			testApp.Services.UserService.Delete(ctx, user.ID)
		}
	}()

	// Set availabilities for each user
	for _, userData := range userMockData.Users {
		calendar := calendarMap[userData.Email]

		availabilityInput := &pkg.SetAvailabilityInput{
			CalendarID: calendar.ID,
			Rules:      userData.Availability.Rules,
		}

		_, err := testApp.Services.ScheduleService.SetAvailability(ctx, availabilityInput)
		require.NoError(t, err)
	}

	// Create events for each user
	for _, userEvents := range eventsMockData.Events {
		user := userMap[userEvents.UserEmail]
		calendar := calendarMap[userEvents.UserEmail]
		require.NotNil(t, user, "User not found for email: %s", userEvents.UserEmail)
		require.NotNil(t, calendar, "Calendar not found for email: %s", userEvents.UserEmail)

		for _, eventData := range userEvents.Events {
			description := eventData.Description
			createEventInput := &pkg.CreateEventInput{
				UserID:      user.ID,
				CalendarID:  calendar.ID,
				Title:       eventData.Title,
				Description: &description,
				Start:       eventData.Start,
				End:         eventData.End,
				Invitees:    eventData.Invitees,
			}

			event, err := testApp.Services.ScheduleService.CreateEvent(ctx, createEventInput)
			require.NoError(t, err)
			require.NotNil(t, event)
			require.Equal(t, eventData.Title, event.Title)
			require.Equal(t, eventData.Description, *event.Description)
		}
	}

	// Verify availabilities
	expectedRules := map[string]int{
		"weekday@example.com": 5, // Weekday user (Mon-Fri)
		"midweek@example.com": 4, // Midweek user (Wed-Sat)
		"weekend@example.com": 2, // Weekend user (Sat-Sun)
	}

	for email := range userMap {
		availability, err := testApp.Services.ScheduleService.GetAvailability(ctx, calendarMap[email].ID)
		require.NoError(t, err)
		require.NotNil(t, availability)
		require.Len(t, availability.Rules, expectedRules[email],
			"User %s should have %d availability rules", email, expectedRules[email])
	}
}

func TestFindOverlappingAvailabilities(t *testing.T) {
	ctx := context.Background()
	userMockData := loadUserMockData(t)
	eventsMockData := loadEventsMockData(t)

	// Create users and their calendars first
	userMap := make(map[string]*pkg.User)         // map[email]user
	calendarMap := make(map[string]*pkg.Calendar) // map[email]calendar
	userToCalendarMap := make(map[string]string)  // map[userID]calendarID

	for _, userData := range userMockData.Users {
		user, err := testApp.Services.UserService.Create(ctx, &pkg.CreateUserInput{
			Name:  userData.Name,
			Email: userData.Email,
		})
		require.NoError(t, err)
		require.NotNil(t, user)
		userMap[userData.Email] = user

		calendar, err := testApp.Services.ScheduleService.AddCalendarToUser(ctx, user.ID, &pkg.CreateCalendarInput{
			UserID: user.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, calendar)
		calendarMap[userData.Email] = calendar
		userToCalendarMap[user.ID] = calendar.ID
	}

	// Clean up users after tests
	defer func() {
		for _, user := range userMap {
			testApp.Services.UserService.Delete(ctx, user.ID)
		}
	}()

	// Set availabilities for each user
	for _, userData := range userMockData.Users {
		calendar := calendarMap[userData.Email]

		availabilityInput := &pkg.SetAvailabilityInput{
			CalendarID: calendar.ID,
			Rules:      userData.Availability.Rules,
		}

		_, err := testApp.Services.ScheduleService.SetAvailability(ctx, availabilityInput)
		require.NoError(t, err)
	}

	for _, userEvents := range eventsMockData.Events {
		user := userMap[userEvents.UserEmail]
		calendar := calendarMap[userEvents.UserEmail]
		require.NotNil(t, user, "User not found for email: %s", userEvents.UserEmail)
		require.NotNil(t, calendar, "Calendar not found for email: %s", userEvents.UserEmail)

		for _, eventData := range userEvents.Events {
			description := eventData.Description
			createEventInput := &pkg.CreateEventInput{
				UserID:      user.ID,
				CalendarID:  calendar.ID,
				Title:       eventData.Title,
				Description: &description,
				Start:       eventData.Start,
				End:         eventData.End,
				Invitees:    eventData.Invitees,
			}

			event, err := testApp.Services.ScheduleService.CreateEvent(ctx, createEventInput)
			require.NoError(t, err)
			require.NotNil(t, event)
			require.Equal(t, eventData.Title, event.Title)
			require.Equal(t, eventData.Description, *event.Description)
		}
	}

	// Test cases for finding overlapping availabilities
	testCases := []struct {
		name          string
		userEmails    []string
		startTime     string
		endTime       string
		expectedSlots int
	}{
		{
			name:          "All_Users_No_Overlap",
			userEmails:    []string{"weekday@example.com", "midweek@example.com", "weekend@example.com"},
			startTime:     "2024-12-16T09:00:00Z", // Monday
			endTime:       "2024-12-16T17:00:00Z",
			expectedSlots: 0,
		},
		{
			name:          "Weekday_Midweek_Friday_Overlap",
			userEmails:    []string{"weekday@example.com", "midweek@example.com"},
			startTime:     "2024-12-20T09:00:00Z", // Friday
			endTime:       "2024-12-20T17:00:00Z",
			expectedSlots: 2,
		},
		{
			name:          "Midweek_Weekend_Saturday_Overlap",
			userEmails:    []string{"midweek@example.com", "weekend@example.com"},
			startTime:     "2024-12-21T09:00:00Z", // Saturday
			endTime:       "2024-12-21T17:00:00Z",
			expectedSlots: 1,
		},
		{
			name:          "No_Common_Availability",
			userEmails:    []string{"weekday@example.com", "weekend@example.com"},
			startTime:     "2024-12-21T09:00:00Z", // Saturday
			endTime:       "2024-12-21T17:00:00Z",
			expectedSlots: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var calendarIDs []string
			for _, email := range tc.userEmails {
				user := userMap[email]
				require.NotNil(t, user, "User not found for email: %s", email)
				calendarID := userToCalendarMap[user.ID]
				require.NotEmpty(t, calendarID, "Calendar not found for user: %s", user.ID)
				calendarIDs = append(calendarIDs, calendarID)
			}

			overlappingSlots, err := testApp.Services.ScheduleService.FindOverlappingAvailabilities(ctx, &pkg.FindOverlappingAvailabilitiesInput{
				CalendarIDs: calendarIDs,
				From:        tc.startTime,
				To:          tc.endTime,
			})

			require.NoError(t, err)
			require.NotNil(t, overlappingSlots)

			require.Equal(t, tc.expectedSlots, len(overlappingSlots),
				"Expected %d overlapping slots but got %d",
				tc.expectedSlots, len(overlappingSlots))
		})
	}
}
