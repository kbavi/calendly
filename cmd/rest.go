package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kbavi/calendly/app"
	"github.com/kbavi/calendly/rest"
)

func main() {
	os.Setenv("TZ", "UTC")
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	if env == "local" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Warning: Error loading .env file: %v\n", err)
		}
	}
	app := app.NewApp()
	router := gin.Default()
	initRoutes(router, app.Services)
	router.Run(":8080")
}

func initRoutes(router *gin.Engine, services app.Services) {

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	userHandler := rest.NewUserHandler(services.UserService)
	router.POST("/api/v1/users", userHandler.Create)
	router.GET("/api/v1/users/:id", userHandler.Get)
	router.DELETE("/api/v1/users/:id", userHandler.Delete)

	calendarHandler := rest.NewCalendarHandler(services.ScheduleService)
	router.POST("/api/v1/calendars", calendarHandler.Create)
	router.GET("/api/v1/calendars/:calendar_id", calendarHandler.Get)
	router.POST("/api/v1/calendars/availabilities/overlap", calendarHandler.OverlappingAvailabilities)
	router.GET("/book/:calendar_id/slots/:duration_minutes", calendarHandler.GetEventSlots)

	availabilityHandler := rest.NewAvailabilityHandler(services.ScheduleService)
	router.POST("/api/v1/availabilities/:user_id", availabilityHandler.SetAvailability)
	router.GET("/api/v1/availabilities/:user_id", availabilityHandler.GetAvailability)

	eventHandler := rest.NewEventHandler(services.ScheduleService)
	router.POST("/api/v1/events", eventHandler.CreateEvent)
	router.GET("/api/v1/events", eventHandler.GetEventsByCalendarIDsInRange)
	router.GET("/api/v1/events/:event_id", eventHandler.GetEvent)
}
