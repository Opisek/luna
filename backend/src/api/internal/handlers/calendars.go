package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/cache"
	"luna-backend/errors"
	"luna-backend/types"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var _ binding.BindUnmarshaler = (*types.Color)(nil) // necessary for https://gin-gonic.com/en/docs/binding/bind-custom-unmarshaler/

type exposedCalendar struct {
	Id           types.ID     `json:"id"`
	Source       types.ID     `json:"source"`
	Name         string       `json:"name"`
	Desc         string       `json:"desc,omitempty"`
	Color        *types.Color `json:"color"`
	Overridden   bool         `json:"overridden"`
	CanEdit      bool         `json:"can_edit"` // TODO: might exclude from here and add to "detailed" view instead
	CanDelete    bool         `json:"can_delete"`
	CanAddEvents bool         `json:"can_add_events"`
}

func GetCalendars(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	sourceId, err := util.GetId(c, "source")
	if err != nil {
		u.Error(err)
		return
	}

	// Get the specified source
	source, err := cache.GetCached(u.Config.Cache, userId, sourceId, u.Context, func() (types.Source, *errors.ErrorTrace) {
		return u.Tx.Queries().GetSource(userId, sourceId, u.Context, u.Config)
	})
	if err != nil {
		u.Error(err)
		return
	}

	// Get the associated calendars
	calsFromSource, err := source.GetCalendars(u.Tx.Queries())

	if err != nil {
		u.Error(err)
		return
	}

	cals, err := u.Tx.Queries().OverrideCalendars(calsFromSource)
	if err != nil {
		u.Error(err)
		return
	}

	// Convert to exposed format
	convertedCals := make([]exposedCalendar, len(cals))
	for i, cal := range cals {
		u.Config.Cache.Cache(userId, cal)

		convertedCals[i] = exposedCalendar{
			Id:         cal.GetId(),
			Source:     cal.GetSource().GetId(),
			Name:       cal.GetName(),
			Desc:       cal.GetDesc(),
			Color:      cal.GetColor(),
			Overridden: cal.GetOverridden(),
			//Settings: cal.GetSettings(),
			CanEdit:      cal.CanEdit(),
			CanDelete:    cal.CanDelete(),
			CanAddEvents: cal.CanAddEvents(),
		}
	}

	u.Success(&gin.H{"calendars": convertedCals})
}

func GetCalendar(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	calendarId, err := util.GetId(c, "calendar")
	if err != nil {
		u.Error(err)
		return
	}

	// Get calendar
	calFromSource, err := u.Tx.Queries().GetCalendar(userId, calendarId, u.Context, u.Config)
	if err != nil {
		u.Error(err)
		return
	}

	cal, err := u.Tx.Queries().OverrideCalendar(calFromSource)
	if err != nil {
		u.Error(err)
		return
	}

	u.Config.Cache.Cache(userId, cal)

	// Convert to exposed format
	convertedCal := exposedCalendar{
		Id:         cal.GetId(),
		Source:     cal.GetSource().GetId(),
		Name:       cal.GetName(),
		Desc:       cal.GetDesc(),
		Color:      cal.GetColor(),
		Overridden: cal.GetOverridden(),
		//Settings: cal.GetSettings(),
		CanEdit:      cal.CanEdit(),
		CanDelete:    cal.CanDelete(),
		CanAddEvents: cal.CanAddEvents(),
	}

	u.Success(&gin.H{"calendar": convertedCal})
}

func PutCalendar(c *gin.Context, body *struct {
	Name  string      `json:"name" form:"name" binding:"required"`
	Desc  string      `json:"desc" form:"desc" binding:"omitempty"`
	Color types.Color `json:"color" form:"color" binding:"required"`
}) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	sourceId, tr := util.GetId(c, "source")
	if tr != nil {
		u.Error(tr)
		return
	}

	source, tr := cache.GetCached(u.Config.Cache, userId, sourceId, u.Context, func() (types.Source, *errors.ErrorTrace) {
		return u.Tx.Queries().GetSource(userId, sourceId, u.Context, u.Config)
	})
	if tr != nil {
		u.Error(tr)
		return
	}

	cal, tr := source.AddCalendar(body.Name, body.Desc, &body.Color, u.Tx.Queries())
	if tr != nil {
		u.Error(tr)
		return
	}

	tr = u.Tx.Queries().InsertCalendar(cal)
	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(&gin.H{"id": cal.GetId().String()})
}

func PatchCalendar(c *gin.Context, body *struct {
	Name       *string      `json:"name" form:"name" binding:"omitempty"`
	Desc       *string      `json:"desc" form:"desc" binding:"omitempty"`
	Color      *types.Color `json:"color" form:"color" binding:"required"`
	Overridden bool         `json:"overridden" form:"overridden"`
}) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	calendarId, err := util.GetId(c, "calendar")
	if err != nil {
		u.Error(err)
		return
	}

	calendar, err := u.Tx.Queries().GetCalendar(userId, calendarId, u.Context, u.Config)
	if err != nil {
		u.Error(err)
		return
	}

	if body.Name == nil {
		oldName := calendar.GetName()
		body.Name = &oldName
	}

	if body.Desc == nil {
		oldDesc := calendar.GetDesc()
		body.Desc = &oldDesc
	}

	if body.Color == nil {
		oldColor := calendar.GetColor()
		body.Color = oldColor
	}

	_, err = calendar.GetSource().EditCalendar(calendar, *body.Name, *body.Desc, body.Color, body.Overridden, u.Tx.Queries())
	if err != nil {
		u.Error(err)
		return
	}

	u.Success(nil)
}

func DeleteCalendar(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	calendarId, err := util.GetId(c, "calendar")
	if err != nil {
		u.Error(err)
		return
	}

	calendar, err := u.Tx.Queries().GetCalendar(userId, calendarId, u.Context, u.Config)
	if err != nil {
		u.Error(err)
		return
	}

	err = calendar.GetSource().DeleteCalendar(calendar, u.Tx.Queries())
	if err != nil {
		u.Error(err)
		return
	}

	err = u.Tx.Queries().DeleteCalendar(userId, calendarId)
	if err != nil {
		u.Error(err)
		return
	}

	u.Success(nil)
}

func ChangeCalendarDisplayOrder(c *gin.Context, body *struct {
	Index *uint16 `json:"index" form:"index" binding:"required"`
}) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	calendarId, tr := util.GetId(c, "calendar")
	if tr != nil {
		u.Error(tr)
		return
	}

	tr = u.Tx.Queries().UpdateCalendarDisplayOrder(userId, calendarId, *body.Index)
	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(nil)
}
