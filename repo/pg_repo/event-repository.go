package pg_repo

import (
	"context"
	"time"

	"github.com/kbavi/calendly/db/models"
	"github.com/kbavi/calendly/pkg"
	"github.com/kbavi/calendly/repo"
	"gorm.io/gorm"
)

func NewEventRepository(db *gorm.DB) repo.EventRepository {
	db.AutoMigrate(&models.EventModel{})
	return &eventRepository{db: db}
}

type eventRepository struct {
	db *gorm.DB
}

func (r *eventRepository) CreateEvent(ctx context.Context, input *pkg.CreateEventInput) (*pkg.Event, error) {
	id := repo.GenerateID()

	start, err := time.Parse(time.RFC3339, input.Start)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse(time.RFC3339, input.End)
	if err != nil {
		return nil, err
	}
	event := &models.EventModel{
		ID:          id,
		CalendarID:  input.CalendarID,
		Title:       input.Title,
		Description: input.Description,
		Start:       start,
		Ending:      end,
		Invitees:    input.Invitees,
	}
	if err := r.db.Create(event).Error; err != nil {
		return nil, err
	}
	return event.ToEvent(), nil
}

func (r *eventRepository) GetEvent(ctx context.Context, eventID string) (*pkg.Event, error) {
	event := &models.EventModel{}
	if err := r.db.First(event, "id = ?", eventID).Error; err != nil {
		return nil, err
	}
	return event.ToEvent(), nil
}

func (r *eventRepository) GetEventsByCalendarIDsInRange(ctx context.Context, calendarIDs []string, start time.Time, end time.Time) ([]pkg.Event, error) {
	events := []models.EventModel{}
	if err := r.db.Where("calendar_id IN (?) AND start >= ? AND ending <= ?", calendarIDs, start, end).Find(&events).Error; err != nil {
		return nil, err
	}
	pkgEvents := []pkg.Event{}
	for _, event := range events {
		pkgEvents = append(pkgEvents, *event.ToEvent())
	}
	return pkgEvents, nil
}
