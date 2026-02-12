package investigo

import (
	"errors"
	"time"
)

type Exchange struct {
	Country    string  `json:"country"`
	CountryId  Country `json:"country_id"`
	Crypto     bool    `json:"crypto"`
	Id         int     `json:"id"`
	LongName   string  `json:"long_name"`
	MarketLink string  `json:"market_link"`
	ShortName  string  `json:"short_name"`
	TimeZone   string  `json:"time_zone"`
}

type Holiday struct {
	Exchange       *Exchange  `json:"exchange"`
	ExchangeClosed bool       `json:"exchange_closed"`
	ExchangeId     int        `json:"exchange_id"`
	HolidayEnd     *time.Time `json:"holiday_end"`
	HolidayId      int        `json:"holiday_id"`
	HolidayName    string     `json:"holiday_name"`
	HolidayStart   *time.Time `json:"holiday_start"`
}

type HolidaysResponse struct {
	Holidays       []Holiday `json:"holidays"`
	NextPageCursor string    `json:"next_page_cursor"`
}

func (i *Investigo) Holidays(filter *Filter) (*HolidaysResponse, error) {
	result := new(HolidaysResponse)
	_, err := i.rb.GetToStruct("https://endpoints.investing.com/pd-instruments/v1/calendars/holidays"+filter.String(), result)
	if err != nil {
		return nil, errors.New("fail to get economic calendar: " + err.Error())
	}
	return result, nil
}
