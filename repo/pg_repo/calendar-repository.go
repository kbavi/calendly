package pg_repo

import (
	"context"

	"github.com/kbavi/calendly/db/models"
	"github.com/kbavi/calendly/pkg"
	"github.com/kbavi/calendly/repo"
	"gorm.io/gorm"
)

func NewCalendarRepository(db *gorm.DB) repo.CalendarRepository {
	db.AutoMigrate(&models.CalendarModel{})
	return &calendarRepository{db: db}
}

type calendarRepository struct {
	db *gorm.DB
}

func (r *calendarRepository) AddCalendarToUser(ctx context.Context, userID string, input *pkg.CreateCalendarInput) (*pkg.Calendar, error) {
	id := repo.GenerateID()

	calendar := &models.CalendarModel{
		ID:     id,
		UserID: userID,
	}
	if err := r.db.Create(calendar).Error; err != nil {
		return nil, err
	}
	return &pkg.Calendar{
		ID: calendar.ID,
	}, nil
}

func (r *calendarRepository) GetCalendar(ctx context.Context, calendarID string) (*pkg.Calendar, error) {
	calendar := &models.CalendarModel{}
	if err := r.db.First(calendar, "id = ?", calendarID).Error; err != nil {
		return nil, err
	}
	return &pkg.Calendar{
		ID: calendar.ID,
	}, nil
}

func (r *calendarRepository) GetCalendarsByIDs(ctx context.Context, calendarIDs []string) ([]pkg.Calendar, error) {
	calendars := []models.CalendarModel{}
	if err := r.db.Where("id IN (?)", calendarIDs).Find(&calendars).Error; err != nil {
		return nil, err
	}
	calendarsDTO := []pkg.Calendar{}
	for _, calendar := range calendars {
		calendarsDTO = append(calendarsDTO, pkg.Calendar{
			ID: calendar.ID,
		})
	}
	return calendarsDTO, nil
}
