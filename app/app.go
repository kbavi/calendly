package app

import (
	"log"
	"os"

	"github.com/kbavi/calendly/pkg/schedule"
	"github.com/kbavi/calendly/pkg/user"
	"github.com/kbavi/calendly/repo/pg_repo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Services struct {
	UserService     user.Service
	ScheduleService schedule.ScheduleService
}

type App struct {
	Services Services
}

func NewApp() *App {
	client, err := gorm.Open(postgres.Open(os.Getenv("PG_DSN")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	userRepo := pg_repo.NewUserRepository(client)
	calendarRepo := pg_repo.NewCalendarRepository(client)
	availabilityRepo := pg_repo.NewAvailabilityRepository(client)
	eventRepo := pg_repo.NewEventRepository(client)

	userService := user.NewService(userRepo)
	scheduleService := schedule.NewService(userService, calendarRepo, availabilityRepo, eventRepo)

	return &App{
		Services: Services{
			UserService:     userService,
			ScheduleService: scheduleService,
		},
	}
}
