package investigo

import (
	"errors"
	"time"
)

type ForecastType string

const (
	Positive ForecastType = "positive"
	Negative ForecastType = "negative"
	Neutral  ForecastType = "neutral"
)

type EventType string

const (
	Speech EventType = "speech"
)

type Event struct {
	Category         string    `json:"category"`
	CountryId        Country   `json:"country_id"`
	Currency         string    `json:"currency"`
	Description      string    `json:"description"`
	EventCycleSuffix string    `json:"event_cycle_suffix,omitempty"`
	EventId          int       `json:"event_id"`
	EventMetaTitle   string    `json:"event_meta_title"`
	EventTranslated  string    `json:"event_translated"`
	Importance       string    `json:"importance"`
	LongName         string    `json:"long_name"`
	PageLink         string    `json:"page_link"`
	ShortName        string    `json:"short_name"`
	Source           string    `json:"source"`
	SourceUrl        string    `json:"source_url"`
	EventType        EventType `json:"event_type,omitempty"`
}

type Occurrence struct {
	Actual              float64      `json:"actual,omitempty"`
	ActualToForecast    ForecastType `json:"actual_to_forecast"`
	EventId             int          `json:"event_id"`
	Forecast            float64      `json:"forecast,omitempty"`
	OccurrenceId        int          `json:"occurrence_id"`
	OccurrenceTime      *time.Time   `json:"occurrence_time"`
	Precision           int          `json:"precision"`
	Preliminary         bool         `json:"preliminary"`
	Previous            float64      `json:"previous,omitempty"`
	PreviousRevisedFrom float64      `json:"previous_revised_from,omitempty"`
	ReferencePeriod     string       `json:"reference_period,omitempty"`
	RevisedToPrevious   ForecastType `json:"revised_to_previous"`
	Unit                string       `json:"unit,omitempty"`
}

type EconomicCalendarResponse struct {
	Events         []Event      `json:"events"`
	Occurrences    []Occurrence `json:"occurrences"`
	NextPageCursor string       `json:"next_page_cursor"`
}

func (i *Investigo) EconomicCalendar(filter *Filter) (*EconomicCalendarResponse, error) {
	result := new(EconomicCalendarResponse)
	_, err := i.rb.GetToStruct("https://endpoints.investing.com/pd-instruments/v1/calendars/economic/events/occurrences"+filter.String(), result)
	if err != nil {
		return nil, errors.New("fail to get economic calendar: " + err.Error())
	}
	return result, nil
}
