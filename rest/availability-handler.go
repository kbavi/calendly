package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kbavi/calendly/pkg"
	"github.com/kbavi/calendly/pkg/schedule"
)

type AvailabilityHandler interface {
	SetAvailability(c *gin.Context)
	GetAvailability(c *gin.Context)
}

func NewAvailabilityHandler(scheduleService schedule.ScheduleService) AvailabilityHandler {
	return &availabilityHandler{scheduleService: scheduleService}
}

type availabilityHandler struct {
	scheduleService schedule.ScheduleService
}

func (h *availabilityHandler) SetAvailability(c *gin.Context) {
	var input pkg.SetAvailabilityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if errs := input.Validate(); len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	availability, err := h.scheduleService.SetAvailability(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pkg.ReturnAvailabilityResponse{
		Status: "success",
		Data: map[string]pkg.AvailabilityDTO{
			"availability": {
				ID: availability.ID,

				Calendar: pkg.CalendarDTO{
					ID: availability.Calendar.ID,
				},
				Rules: availability.Rules,
			},
		},
	})
}

func (h *availabilityHandler) GetAvailability(c *gin.Context) {
	calendarID := c.Param("calendar_id")
	if calendarID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "calendar_id is required"})
		return
	}
	availability, err := h.scheduleService.GetAvailability(c.Request.Context(), calendarID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, availability)
}
