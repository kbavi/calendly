package pg_repo

import (
	"context"
	"encoding/json"

	"github.com/kbavi/calendly/db/models"
	"github.com/kbavi/calendly/pkg"
	"github.com/kbavi/calendly/repo"
	"gorm.io/gorm"
)

type availabilityRepository struct {
	db *gorm.DB
}

func NewAvailabilityRepository(db *gorm.DB) repo.AvailabilityRepository {
	db.AutoMigrate(&models.AvailabilityModel{})
	return &availabilityRepository{db: db}
}

func (r *availabilityRepository) SetAvailability(ctx context.Context, input *pkg.SetAvailabilityInput) (*pkg.Availability, error) {
	id := repo.GenerateID()

	rules, err := json.Marshal(input.Rules)
	if err != nil {
		return nil, err
	}

	availability := &models.AvailabilityModel{
		ID:         id,
		CalendarID: input.CalendarID,
		Rules:      string(rules),
	}

	existingAvailability, err := r.GetAvailabilityByCalendarID(ctx, input.CalendarID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err := r.db.Create(availability).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&models.CalendarModel{}).Where("id = ?", input.CalendarID).Update("availability_id", id).Error; err != nil {
		return nil, err
	}

	if existingAvailability != nil {
		err = r.DeleteAvailability(ctx, existingAvailability.ID)
		if err != nil {
			return nil, err
		}
	}

	rulesMarshalled, err := parseAvailabilityRules(availability.Rules)
	if err != nil {
		return nil, err
	}
	return &pkg.Availability{
		ID:    availability.ID,
		Rules: rulesMarshalled,
	}, nil
}

func parseAvailabilityRules(rules string) ([]pkg.AvailabilityRule, error) {
	var rulesMarshalled []pkg.AvailabilityRule
	if err := json.Unmarshal([]byte(rules), &rulesMarshalled); err != nil {
		return nil, err
	}
	return rulesMarshalled, nil
}

func (r *availabilityRepository) DeleteAvailability(ctx context.Context, availabilityID string) error {
	return r.db.Model(&models.AvailabilityModel{}).Unscoped().Delete(&models.AvailabilityModel{}, "id = ?", availabilityID).Error
}

func (r *availabilityRepository) GetAvailabilityByCalendarID(ctx context.Context, calendarID string) (*pkg.Availability, error) {
	availability := &models.AvailabilityModel{}
	if err := r.db.First(availability, "calendar_id = ?", calendarID).Error; err != nil {
		return nil, err
	}
	return availability.ToAvailability(), nil
}

func (r *availabilityRepository) GetAvailability(ctx context.Context, calendarID string) (*pkg.Availability, error) {
	availability := &models.AvailabilityModel{}
	if err := r.db.First(availability, "calendar_id = ?", calendarID).Error; err != nil {
		return nil, err
	}
	return availability.ToAvailability(), nil
}

func (r *availabilityRepository) GetAvailabilityRulesByCalendarIDs(ctx context.Context, calendarIDs []string) (map[string][]pkg.AvailabilityRule, error) {
	availabilities := []models.AvailabilityModel{}
	if err := r.db.Find(&availabilities, "calendar_id IN (?)", calendarIDs).Error; err != nil {
		return nil, err
	}
	availabilityRules := map[string][]pkg.AvailabilityRule{}

	for _, availability := range availabilities {
		rulesMarshalled, err := parseAvailabilityRules(availability.Rules)
		if err != nil {
			return nil, err
		}
		availabilityRules[availability.CalendarID] = rulesMarshalled
	}
	return availabilityRules, nil
}
