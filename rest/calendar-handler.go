package rest

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kbavi/calendly/pkg"
	"github.com/kbavi/calendly/pkg/schedule"
	"gorm.io/gorm"
)

type CalendarHandler struct {
	scheduleService schedule.ScheduleService
}

func NewCalendarHandler(scheduleService schedule.ScheduleService) *CalendarHandler {
	return &CalendarHandler{scheduleService: scheduleService}
}

func (h *CalendarHandler) Create(c *gin.Context) {
	var input pkg.CreateCalendarInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if errs := input.Validate(); len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errs})
		return
	}
	calendar, err := h.scheduleService.AddCalendarToUser(c.Request.Context(), input.UserID, &input)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	calendarDTO := pkg.CalendarDTO{}
	calendarDTO.FromCalendar(calendar)
	c.JSON(http.StatusOK, pkg.ReturnCalendarResponse{
		Status: "success",
		Data: map[string]pkg.CalendarDTO{
			"calendar": calendarDTO,
		},
	})
}

func (h *CalendarHandler) Get(c *gin.Context) {
	calendarID := c.Param("calendar_id")
	if calendarID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "calendar_id is required"})
		return
	}
	fromStr := c.Query("from")
	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from time format, should be RFC3339"})
		return
	}
	toStr := c.Query("to")
	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to time format, should be RFC3339"})
		return
	}
	// get events from between from and to
	events, err := h.scheduleService.GetEventsByCalendarIDsInRange(c.Request.Context(), []string{calendarID}, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// get availability rules for the calendar
	availability, err := h.scheduleService.GetAvailability(c.Request.Context(), calendarID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// get free intervals for the calendar between from and to
	var input pkg.FindOverlappingAvailabilitiesInput
	input.CalendarIDs = []string{calendarID}
	input.From = from.Format(time.RFC3339)
	input.To = to.Format(time.RFC3339)
	freeIntervals, err := h.scheduleService.FindOverlappingAvailabilities(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fifteenMinsSchedulingLink := h.scheduleService.GetSchedulingLink(c.Request.Context(), calendarID, 15)
	thirtyMinsSchedulingLink := h.scheduleService.GetSchedulingLink(c.Request.Context(), calendarID, 30)
	oneHourSchedulingLink := h.scheduleService.GetSchedulingLink(c.Request.Context(), calendarID, 60)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"calendar_id":    calendarID,
			"events":         events,
			"availability":   availability,
			"free_intervals": freeIntervals,
			"scheduling_links": map[string]string{
				"15-mins": fifteenMinsSchedulingLink,
				"30-mins": thirtyMinsSchedulingLink,
				"60-mins": oneHourSchedulingLink,
			},
		},
	})
}

func (h *CalendarHandler) GetEventSlots(c *gin.Context) {
	calendarID := c.Param("calendar_id")
	if calendarID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "calendar_id is required"})
		return
	}
	durationMinutesStr := c.Param("duration_minutes")
	durationMinutesStr = strings.TrimSuffix(durationMinutesStr, "-mins")
	durationMinutes, err := strconv.Atoi(durationMinutesStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid duration_minutes"})
		return
	}
	fromStr := c.Query("from")
	var from time.Time
	if fromStr == "" {
		from = time.Now()
	} else {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from time format, should be RFC3339"})
			return
		}
	}
	toStr := c.Query("to")
	var to time.Time
	if toStr == "" {
		to = from.AddDate(0, 0, 1)
	} else {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to time format, should be RFC3339"})
			return
		}
	}
	var input pkg.BookBySchedulingLinkInput
	input.CalendarID = calendarID
	input.DurationMinutes = durationMinutes
	input.From = from
	input.To = to
	slots, err := h.scheduleService.GetEventSlots(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	schedulingLink := h.scheduleService.GetSchedulingLink(c.Request.Context(), calendarID, durationMinutes)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{
		"scheduling_link": schedulingLink,
		"slots":           slots,
	}})
}

func (h *CalendarHandler) OverlappingAvailabilities(c *gin.Context) {
	var input pkg.FindOverlappingAvailabilitiesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if errs := input.Validate(); len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errs})
		return
	}
	overlaps, err := h.scheduleService.FindOverlappingAvailabilities(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   overlaps,
	})
}
