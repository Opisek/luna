package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/cache"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var _ binding.BindUnmarshaler = (*types.Color)(nil) // necessary for https://gin-gonic.com/en/docs/binding/bind-custom-unmarshaler/

type exposedEvent struct {
	Id         types.ID         `json:"id"`
	Calendar   types.ID         `json:"calendar"`
	Name       string           `json:"name"`
	Desc       string           `json:"desc,omitempty"`
	Color      *types.Color     `json:"color"`
	Date       *types.EventDate `json:"date"`
	Overridden bool             `json:"overridden"`
	CanEdit    bool             `json:"can_edit"` // TODO: might exclude from here and add to "detailed" view instead
	CanDelete  bool             `json:"can_delete"`
	Settings   any              `json:"settings,omitempty"` // TODO: delete, this is temporary for debugging recurrences
}

func GetEvents(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	calendarId, tr := util.GetId(c, "calendar")
	if tr != nil {
		u.Error(tr)
		return
	}

	// Get the requested calendar
	calendar, tr := cache.GetCached(u.Config.Cache, userId, calendarId, u.Context, func() (types.Calendar, *errors.ErrorTrace) {
		return u.Tx.Queries().GetCalendar(userId, calendarId, u.Context, u.Config)
	})
	if tr != nil {
		u.Error(tr)
		return
	}

	// Get the associated events
	startStr := c.Query("start")
	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		u.Error(errors.New().
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Missing or malformed start time"))
		return
	}
	endStr := c.Query("end")
	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		u.Error(errors.New().
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Missing or malformed end time"))
		return
	}
	if startTime.After(endTime) {
		u.Error(errors.New().
			Append(errors.LvlPlain, "Start time must not be after end time"))
		return
	}
	if endTime.Sub(startTime) > time.Hour*24*365 {
		endTime = startTime.Add(time.Hour * 24 * 365)
	}

	eventsFromCal, tr := calendar.GetEvents(startTime, endTime, u.Tx.Queries())
	if tr != nil {
		u.Error(tr)
		return
	}

	// Save in the database (the "parent" events need to exist in the database when inserting recurrence instances!)
	if _, tr = u.Tx.Queries().OverrideEvents(eventsFromCal); tr != nil {
		u.Error(tr)
		return
	}

	// Expand recurring events
	expandedEvents := make([]types.Event, len(eventsFromCal))
	count := 0
	for _, event := range eventsFromCal {
		expanded, tr := types.ExpandRecurrence(event, &startTime, &endTime)
		if tr != nil {
			u.Error(tr)
			return
		}

		if len(expanded) > 1 {
			newRes := make([]types.Event, len(expandedEvents)-1+len(expanded))
			copy(newRes, expandedEvents[:count])
			expandedEvents = newRes
		}

		for _, e := range expanded {
			expandedEvents[count] = e
			count++
		}
	}

	// Save in the database and apply overrides to the specific recurrence instances
	events, tr := u.Tx.Queries().OverrideEvents(expandedEvents[:count])
	if tr != nil {
		u.Error(tr)
		return
	}

	// Convert to exposed format
	convertedEvents := make([]exposedEvent, count)
	for i, event := range events {
		if event.GetName() == "" { // TODO: error handling
			continue
		}

		convertedEvents[i] = exposedEvent{
			Id:         event.GetId(),
			Calendar:   event.GetCalendar().GetId(),
			Name:       event.GetName(),
			Desc:       event.GetDesc(),
			Color:      event.GetColor(),
			Date:       event.GetDate(),
			Overridden: event.GetOverridden(),
			CanEdit:    event.CanEdit(),
			CanDelete:  event.CanDelete(),
			Settings:   event.GetSettings(),
		}
	}

	u.Success(&gin.H{"events": convertedEvents})
}

func GetEvent(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	eventId, tr := util.GetId(c, "event")
	if tr != nil {
		u.Error(tr)
		return
	}

	// Get event
	eventFromCal, err := u.Tx.Queries().GetEvent(userId, eventId, u.Context, u.Config)
	if err != nil {
		u.Error(err)
		return
	}

	event, err := u.Tx.Queries().OverrideEvent(eventFromCal)
	if err != nil {
		u.Error(err)
		return
	}

	// Convert to exposed format
	convertedCal := exposedEvent{
		Id:         event.GetId(),
		Calendar:   event.GetCalendar().GetId(),
		Name:       event.GetName(),
		Desc:       event.GetDesc(),
		Color:      event.GetColor(),
		Date:       event.GetDate(),
		Overridden: event.GetOverridden(),
		//Settings: event.GetSettings(),
		CanEdit:   event.CanEdit(),
		CanDelete: event.CanDelete(),
		Settings:  event.GetSettings(),
	}

	u.Success(&gin.H{"event": convertedCal})
}

func PutEvent(c *gin.Context, body *struct {
	Name     string         `json:"name" form:"name" binding:"required"`
	Desc     string         `json:"desc" form:"desc" binding:"omitempty"`
	Color    types.Color    `json:"color" form:"color" binding:"required"`
	AllDay   bool           `json:"date.all_day" form:"date_all_day"`
	Start    time.Time      `json:"date.start" form:"date_start" binding:"required"`
	End      *time.Time     `json:"date.end" form:"date_end"`
	Duration *time.Duration `json:"date.duration" form:"date_duration"`
	Rrule    *string        `json:"date.rrule" form:"date_rrule"`
	Rdate    *string        `json:"date.rdate" form:"date_rdate"`
	Exdate   *string        `json:"date.exdate" form:"date_exdate"`
}) {
	u := util.GetUtil(c)
	userId := util.GetUserId(c)

	// Find the associated calendar
	calendarId, tr := util.GetId(c, "calendar")
	if tr != nil {
		u.Error(tr)
		return
	}
	calendar, tr := cache.GetCached(u.Config.Cache, userId, calendarId, u.Context, func() (types.Calendar, *errors.ErrorTrace) {
		return u.Tx.Queries().GetCalendar(userId, calendarId, u.Context, u.Config)
	})
	if tr != nil {
		u.Error(tr)
		return
	}

	// Event date
	var date *types.EventDate
	if body.End == nil && body.Duration == nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Must specify end time or duration"))
		return
	} else if body.End != nil && body.Duration != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Cannot specify both end time and duration"))
		return
	} else if body.End != nil {
		date = types.NewEventDateFromEndTime(&body.Start, body.End, body.AllDay, nil)
	} else {
		date = types.NewEventDateFromDuration(&body.Start, body.Duration, body.AllDay, nil)
	}

	// Event recurrence
	rrule := ""
	if body.Rrule != nil {
		rrule = *body.Rrule
	}
	rdate := ""
	if body.Rdate != nil {
		rdate = *body.Rdate
	}
	exdate := ""
	if body.Exdate != nil {
		exdate = *body.Exdate
	}
	recurrence, err := types.EventRecurrenceFromStrings(rrule, rdate, exdate)
	if err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Invalid recurrence").
			AltStr(errors.LvlPlain, "Invalid date"),
		)
	}
	date.SetRecurrence(recurrence)

	// Add the event in the upstream
	event, tr := calendar.AddEvent(body.Name, body.Desc, &body.Color, date, u.Tx.Queries())
	if tr != nil {
		u.Error(tr)
		return
	}

	// Insert the event into the database
	tr = u.Tx.Queries().InsertEvent(event)
	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(&gin.H{"id": event.GetId().String()})
}

func PatchEvent(c *gin.Context, body *struct {
	Name       *string        `json:"name" form:"name" binding:"omitempty"`
	Desc       *string        `json:"desc" form:"desc" binding:"omitempty"`
	Color      *types.Color   `json:"color" form:"color"`
	AllDay     *bool          `json:"date.all_day" form:"date_all_day"`
	Start      *time.Time     `json:"date.start" form:"date_start"`
	End        *time.Time     `json:"date.end" form:"date_end"`
	Duration   *time.Duration `json:"date.duration" form:"date_duration"`
	Overridden bool           `json:"overridden" form:"overridden"`
	Rrule      *string        `json:"date.rrule" form:"date_rrule"`
	Rdate      *string        `json:"date.rdate" form:"date_rdate"`
	Exdate     *string        `json:"date.exdate" form:"date_exdate"`
}, query *struct {
	Affect string `form:"affect,default=this" binding:"required,oneof=this thisandfuture all"`
}) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	eventId, tr := util.GetId(c, "event")
	if tr != nil {
		u.Error(tr)
		return
	}

	event, tr := u.Tx.Queries().GetEvent(userId, eventId, u.Context, u.Config)
	if tr != nil {
		u.Error(tr)
		return
	}

	var newEventDate *types.EventDate
	if body.Start == nil {
		body.Start = event.GetDate().Start()
	}
	if body.End == nil && body.Duration == nil {
		body.End = event.GetDate().End()
	}
	if body.AllDay == nil {
		oldAllDay := event.GetDate().AllDay()
		body.AllDay = &oldAllDay
	}
	if body.End != nil && body.Duration != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Cannot specify both end time and duration"))
		return
	} else if body.End != nil {
		newEventDate = types.NewEventDateFromEndTime(body.Start, body.End, *body.AllDay, nil)
	} else {
		newEventDate = types.NewEventDateFromDuration(body.Start, body.Duration, *body.AllDay, nil)
	}

	// Event recurrence
	newRrule := event.GetDate().Recurrence().RruleString()
	if body.Rrule != nil {
		newRrule = *body.Rrule
	}
	newRdate := event.GetDate().Recurrence().RdateString()
	if body.Rdate != nil {
		newRdate = *body.Rdate
	}
	newExdate := event.GetDate().Recurrence().ExdateString()
	if body.Exdate != nil {
		newExdate = *body.Exdate
	}

	recurrence, err := types.EventRecurrenceFromStrings(newRrule, newRdate, newExdate)
	if err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Invalid recurrence").
			AltStr(errors.LvlPlain, "Invalid date"),
		)
	}
	newEventDate.SetRecurrence(recurrence)

	// If the event does not repeat, we automatically affect "all" instances
	if !event.GetDate().Recurrence().Repeats() {
		query.Affect = "all"
	}

	// Get parent event if needed
	var parentEvent types.Event
	if event.GetParentId() != nil && query.Affect != "this" {
		parentEvent, tr = u.Tx.Queries().GetEvent(userId, *event.GetParentId(), u.Context, u.Config)
		if tr != nil {
			u.Error(tr)
			return
		}
	} else {
		parentEvent = event
		if query.Affect == "thisandfuture" {
			query.Affect = "all"
		}
	}

	switch query.Affect {
	case "this":
		// Editing just one instance
		_, tr = event.GetCalendar().EditEvent(event, body.Name, body.Desc, body.Color, newEventDate, body.Overridden, "this", u.Tx.Queries())
		if tr != nil {
			u.Error(tr)
			return
		}
	case "thisandfuture":
		if event.CanDelete() {
			// Editing this and future instances means we transform this event into a master event and shorten the original one
			if newEventDate.Recurrence().Repeats() {
				ruleSet := parentEvent.GetDate().Recurrence().RuleSet()
				ruleSet.GetRRule().Options.Dtstart = parentEvent.GetDate().Recurrence().RuleSet().GetDTStart()
				parentEvent.GetDate().Recurrence().SetRuleSet(ruleSet)
			}
			if body.Name != nil {
				event.SetName(*body.Name)
			}
			if body.Desc != nil {
				event.SetDesc(*body.Desc)
			}
			if body.Color != nil {
				event.SetColor(body.Color)
			}
			event, tr := event.GetCalendar().AddEvent(event.GetName(), event.GetDesc(), event.GetColor(), newEventDate, u.Tx.Queries())
			if tr != nil {
				u.Error(tr)
				return
			}

			tr = u.Tx.Queries().InsertEvent(event)
			if tr != nil {
				u.Error(tr)
				return
			}

			ruleSet := parentEvent.GetDate().Recurrence().RuleSet()
			ruleSet.GetRRule().Options.Until = event.GetDate().Start().Add(-time.Second)
			ruleSet.GetRRule().Options.Count = 0
			parentEvent.GetDate().Recurrence().SetRuleSet(ruleSet)
			_, tr = parentEvent.GetCalendar().EditEvent(parentEvent, nil, nil, nil, parentEvent.GetDate(), false, "thisandfuture", u.Tx.Queries())
			if tr != nil {
				u.Error(tr)
				return
			}
		} else {
			// In case we cannot delete events (e.g., ical), we use overrides instead
			_, tr = event.GetCalendar().EditEvent(event, body.Name, body.Desc, body.Color, newEventDate, body.Overridden, "thisandfuture", u.Tx.Queries())
			if tr != nil {
				u.Error(tr)
				return
			}
		}
	case "all":
		// Editing all instances = editing parent event

		// If we edit all events, we have to be careful about how we treat changes to the time
		if parentEvent.GetId() != eventId {
			// First, get the relative shift in start/end timestamps
			deltaStart := newEventDate.Start().Sub(*event.GetDate().Start())
			deltaEnd := newEventDate.End().Sub(*event.GetDate().End())

			// Use these offsets to calculate when the master event should starte
			newStart := parentEvent.GetDate().Start().Add(deltaStart)
			newEnd := parentEvent.GetDate().End().Add(deltaEnd)
			newEventDate.SetStart(&newStart)
			newEventDate.SetEnd(&newEnd)
		}

		// Edit parent event
		_, tr = parentEvent.GetCalendar().EditEvent(parentEvent, body.Name, body.Desc, body.Color, newEventDate, body.Overridden, "all", u.Tx.Queries())
		if tr != nil {
			u.Error(tr)
			return
		}
	}

	u.Success(nil)
}

func DeleteEvent(c *gin.Context, query *struct {
	Affect string `form:"affect,default=this" binding:"required,oneof=this thisandfuture all"`
}) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	eventId, err := util.GetId(c, "event")
	if err != nil {
		u.Error(err)
		return
	}

	// Get event first
	event, err := u.Tx.Queries().GetEvent(userId, eventId, u.Context, u.Config)
	if err != nil {
		u.Error(err)
		return
	}

	var parentEvent types.Event
	if event.GetParentId() != nil {
		parentEvent, err = u.Tx.Queries().GetEvent(userId, *event.GetParentId(), u.Context, u.Config)
		if err != nil {
			u.Error(err)
			return
		}
	} else {
		parentEvent = event
	}

	// If the event does not repeat, we automatically affect "all" instances
	if !event.GetDate().Recurrence().Repeats() {
		query.Affect = "all"
	}

	switch query.Affect {
	case "this":
		// Removing just one instance of a recurrence equates to adding that event to EXDATE
		parentEvent.GetDate().Recurrence().AddException(event.GetDate().Start())
		_, err = parentEvent.GetCalendar().EditEvent(parentEvent, nil, nil, nil, parentEvent.GetDate(), false, "this", u.Tx.Queries())
		if err != nil {
			u.Error(err)
			return
		}
		// TODO: In CalDav we also need to delete the VEVENT of that instance so its information is forgotten when we later add it back
	case "thisandfuture":
		// Removing just one instance of a recurrence equates to changing the recurrence end date
		ruleSet := parentEvent.GetDate().Recurrence().RuleSet()
		ruleSet.GetRRule().Options.Until = event.GetDate().Start().Add(-time.Second)
		ruleSet.GetRRule().Options.Count = 0
		parentEvent.GetDate().Recurrence().SetRuleSet(ruleSet)
		_, err = parentEvent.GetCalendar().EditEvent(parentEvent, nil, nil, nil, parentEvent.GetDate(), false, "thisandfuture", u.Tx.Queries())
		if err != nil {
			u.Error(err)
			return
		}
	case "all":
		// Remove the event from the upstream source
		err = parentEvent.GetCalendar().DeleteEvent(parentEvent, u.Tx.Queries())
		if err != nil {
			u.Error(err)
			return
		}

		// Delete event entry from the database
		err = u.Tx.Queries().DeleteEvent(userId, parentEvent.GetId())
		if err != nil {
			u.Error(err)
			return
		}
	}

	u.Success(nil)
}
