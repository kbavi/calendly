package integration

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kbavi/calendly/app"
	"github.com/kbavi/calendly/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	os.Setenv("TZ", "UTC")
}

var testApp *app.App

func TestMain(m *testing.M) {
	// Load test environment variables
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading test env file: %v", err)
	}

	// Verify required env variables
	if os.Getenv("PG_DSN") == "" {
		log.Fatal("PG_DSN environment variable is required")
	}

	// Initialize the application
	testApp = app.NewApp()

	// Run the tests
	code := m.Run()

	os.Exit(code)
}

func TestUserIntegration(t *testing.T) {
	ctx := context.Background()

	t.Run("Create_Get_Delete_User_Flow", func(t *testing.T) {
		// Test creating a user
		createInput := &pkg.CreateUserInput{
			Name:  "Test User",
			Email: "test@example.com",
		}

		createdUser, err := testApp.Services.UserService.Create(ctx, createInput)
		require.NoError(t, err)
		require.NotNil(t, createdUser)
		assert.Equal(t, createInput.Name, createdUser.Name)
		assert.Equal(t, createInput.Email, string(createdUser.Email))
		assert.NotEmpty(t, createdUser.ID)

		// Test getting the created user
		fetchedUser, err := testApp.Services.UserService.Get(ctx, createdUser.ID)
		require.NoError(t, err)
		require.NotNil(t, fetchedUser)
		assert.Equal(t, createdUser.ID, fetchedUser.ID)
		assert.Equal(t, createdUser.Name, fetchedUser.Name)
		assert.Equal(t, createdUser.Email, fetchedUser.Email)

		// Test deleting the user
		err = testApp.Services.UserService.Delete(ctx, createdUser.ID)
		require.NoError(t, err)

		// Verify user is deleted
		deletedUser, err := testApp.Services.UserService.Get(ctx, createdUser.ID)
		assert.Error(t, err)
		assert.Nil(t, deletedUser)
	})

	t.Run("Get_NonExistent_User", func(t *testing.T) {
		user, err := testApp.Services.UserService.Get(ctx, "non-existent-id")
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
