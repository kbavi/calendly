package rest

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kbavi/calendly/pkg"
	"github.com/kbavi/calendly/pkg/schedule"
	"gorm.io/gorm"
)

type EventHandler struct {
	scheduleService schedule.ScheduleService
}

func NewEventHandler(scheduleService schedule.ScheduleService) *EventHandler {
	return &EventHandler{scheduleService: scheduleService}
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	var input pkg.CreateEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if errs := input.Validate(); len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errs})
		return
	}
	event, err := h.scheduleService.CreateEvent(c.Request.Context(), &input)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user or calendar id"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	eventDTO := pkg.EventDTO{}
	eventDTO.FromEvent(event)
	c.JSON(http.StatusOK, pkg.ReturnEventResponse{
		Status: "success",
		Data: map[string]pkg.EventDTO{
			"event": eventDTO,
		},
	})
}

func (h *EventHandler) GetEvent(c *gin.Context) {
	event, err := h.scheduleService.GetEvent(c.Request.Context(), c.Param("event_id"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	eventDTO := pkg.EventDTO{}
	eventDTO.FromEvent(event)
	c.JSON(http.StatusOK, pkg.ReturnEventResponse{
		Status: "success",
		Data: map[string]pkg.EventDTO{
			"event": eventDTO,
		},
	})
}

func (h *EventHandler) GetEventsByCalendarIDsInRange(c *gin.Context) {
	calendarIDsStr := c.Query("calendar_ids")
	if calendarIDsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "calendar_ids is required"})
		return
	}
	calendarIDs := strings.Split(calendarIDsStr, ",")
	if len(calendarIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "calendar_ids is required"})
		return
	}
	start, err := time.Parse(time.RFC3339, c.Query("start"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
		return
	}
	end, err := time.Parse(time.RFC3339, c.Query("end"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}
	events, err := h.scheduleService.GetEventsByCalendarIDsInRange(c.Request.Context(), calendarIDs, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	eventsDTO := []pkg.EventDTO{}
	for _, event := range events {
		eventDTO := pkg.EventDTO{}
		eventDTO.FromEvent(&event)
		eventsDTO = append(eventsDTO, eventDTO)
	}
	c.JSON(http.StatusOK, pkg.ReturnEventsResponse{
		Status: "success",
		Data: map[string][]pkg.EventDTO{
			"events": eventsDTO,
		},
	})
}
