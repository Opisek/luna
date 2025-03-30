package errors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Custom error-handling system for this project created to solve the following problems with using normal errors:
// - Status codes end up being too generic due to the HTTP response being handled by the "top level" handler function.
// - Loss of detail or detail overload when using many nested functions.
// - Inability to dynamically choose the verbosity of returned responses.
// - Limited formatting flexibility

const (
	// Debug Level
	// - Full stack trace of all errors
	// - MIGHT reveal sensitive information
	// - MUST NOT be used in production
	// - MUST be used when returning any internal IDs or errors
	// - e.g. User tries to get calendars => Could not get calendars from file <ID> belonging to source <ID>
	LvlDebug = iota

	// Wordy (verbose) level
	// - Detailed logs for advanced users
	// - Does not pose a security threat in production
	// - MUST use when referring to technical terms (CalDAV, iCal, HTTP Request, ...)
	// - MUST NOT contain internal IDs or errors
	// - e.g. User tries to get calendars => Could not get calendars from iCal file
	LvlWordy

	// Plain (normal) level
	// - Balanced details for standard users
	// - MAY be used with infrastructure terms (Frontend, Backend, Database, File, ...)
	// - MUST NOT be used with technical terms
	// - e.g. User tries to get calendars => Could not get calendars from source <name>
	LvlPlain

	// Broad (minimal) level
	// - Very generic and minimal level of detail
	// - For users with very limited technical knowledge
	// - MAY be used in high security scenarios
	// - Error SHOULD only reflect the just attempted action
	// - e.g. User tries to get calendars => Could not get calendars
	LvlBroad
)

// If no error occurs, nil MUST be returned.
//
// At the first place that an error takes place, a new trace MUST be created:
// errors.New().Status(http.StatusSomeCode)
//
// If the status code is not yet known, it MAY be added later.
//
// If a function call returns a non-nil error trace tr != nil, that trace
// instance SHOULD be returned until the highest level function is reached
// (typically the api handler).
//
// There are limited cases in which an instantiated error trace MAY be ignored
// or replaced with a new one.
//
// Every stack frame including the first one MAY add its own details to the
// error trace, respecting the verbosity levels from above:
// .Append(errors.LvlSomething, "Something happened")
// .AddErr(errors.LvlSomething, err)
//
// The verbosity level MAY be increased in relation to this guide, for example
// making authentication-related endpoints return more generic errors for
// security reasons.
//
// If multiple things are to be reported, that are on the same hierarchal level
// detail-wise, a conjunction can be created:
// .AndStr("Something else also happened")
// .AndErr(err)
//
// If there are multiple levels of verbosity available for same detail level,
// for which there exists a total subset relation a disjunction can be created:
// .Append(errors.LvlDebug, "Very detailed description of something that happened")
// .AltStr(errors.LvlWordy, "Somewhat detailed description of something that happened")
// .AltStr(errors.LvlPlain, "Description of something that happened")
// .AltStr(errors.LvlBroad, "Something happened")
//
// All messages MUST have proper capitalization but MUST NOT end with punctuation.
//
// The top-level function (like the api handler) finally passes a complete error trace
// to the error channel:
// util.GetUtil(c).Error(tr)

type errEntry struct {
	detailLevel int
	message     string
}

type ErrorTrace struct {
	httpCode int
	trace    [][]*errEntry
}

func New() *ErrorTrace {
	tr := &ErrorTrace{
		httpCode: http.StatusInternalServerError,
		trace:    make([][]*errEntry, 0),
	}

	return tr
}

func (tr *ErrorTrace) Status(httpCode int) *ErrorTrace {
	if tr == nil {
		tr = New()
	}

	tr.httpCode = httpCode

	return tr
}

// Add higher level of abstraction
func (tr *ErrorTrace) Append(detailLevel int, msg string, args ...any) *ErrorTrace {
	disjunction := make([]*errEntry, 1)

	disjunction[0] = &errEntry{
		detailLevel: detailLevel,
		message:     fmt.Sprintf(msg, args...),
	}

	tr.trace = append([][]*errEntry{disjunction}, tr.trace...)

	return tr
}

func (tr *ErrorTrace) AddErr(detailLevel int, err error) *ErrorTrace {
	if err != nil {
		tr.Append(detailLevel, err.Error())
	} else {
		tr.Append(detailLevel, "")
	}
	return tr
}

// Add more errors to the last level of abstraction
func (tr *ErrorTrace) AndStr(msg string, args ...any) *ErrorTrace {
	lastDisjunction := tr.trace[len(tr.trace)-1]
	lastErr := lastDisjunction[len(lastDisjunction)-1]
	newMsg := fmt.Sprintf(msg, args...)

	if lastErr.message == "" {
		lastErr.message = newMsg
	} else if newMsg != "" {
		newMsg = strings.ToLower(newMsg[:1]) + newMsg[1:]
		lastErr.message = lastErr.message + " and " + newMsg
	}

	return tr
}

func (tr *ErrorTrace) AndErr(err error) *ErrorTrace {
	if err != nil {
		return tr.AndStr(err.Error())
	} else {
		return tr.AndStr("")
	}
}

// Add more verbosity options to the last level of abstraction
func (tr *ErrorTrace) AltStr(detailLevel int, msg string, args ...any) *ErrorTrace {
	if len(tr.trace) == 0 {
		return tr.Append(detailLevel, msg, args...)
	}

	disjunction := tr.trace[len(tr.trace)-1]

	disjunction = append(disjunction, &errEntry{
		detailLevel: detailLevel,
		message:     fmt.Sprintf(msg, args...),
	})

	tr.trace[len(tr.trace)-1] = disjunction

	return tr
}

func (tr *ErrorTrace) AltErr(detailLevel int, err error) *ErrorTrace {
	return tr.AltStr(detailLevel, err.Error())
}

// For lowering the exposed level of details
func (tr *ErrorTrace) ConvertTo(detailLevel int) *ErrorTrace {
	for _, disjunction := range tr.trace {
		for _, option := range disjunction {
			if option.detailLevel > detailLevel {
				option.detailLevel = detailLevel
			}
		}
	}
	return tr
}

// Serializing will output a suitable error string corresponding to the
func (tr *ErrorTrace) Serialize(detailLevel int) string {
	alreadyAdded := make(map[string]bool)

	count := 0
	msgs := make([]string, len(tr.trace))

	for _, entry := range tr.trace {
		lowest := int(^uint(0) >> 1)
		lowestMsg := ""

		for _, option := range entry {
			if option.detailLevel >= detailLevel && option.detailLevel < lowest {
				lowest = option.detailLevel
				lowestMsg = option.message
			}
		}

		if lowestMsg == "" {
			continue
		}

		lowestMsg = strings.ToUpper(lowestMsg[:1]) + lowestMsg[1:]

		_, msgExists := alreadyAdded[lowestMsg]
		if !msgExists {
			alreadyAdded[lowestMsg] = true
			msgs[count] = lowestMsg
			count++
		}
	}

	if count == 0 {
		msgs = make([]string, 1)
		msgs[0] = http.StatusText(tr.httpCode)
		count = 1
	}

	return strings.Join(msgs[:count], ": ")
}

func (tr *ErrorTrace) SerializeError(detailLevel int) error {
	return errors.New(tr.Serialize(detailLevel))
}

func (tr *ErrorTrace) GetStatus() int {
	return tr.httpCode
}
