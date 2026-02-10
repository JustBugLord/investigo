package investigo

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Country int

const (
	Albania   Country = 95
	Austria   Country = 54
	Belgium   Country = 34
	EuroZone  Country = 72
	Finland   Country = 71
	France    Country = 22
	Germany   Country = 17
	Russia    Country = 56
	UK        Country = 4
	Canada    Country = 6
	US        Country = 5
	Australia Country = 25
	China     Country = 37
	Japan     Country = 35
)

type Category string

const (
	Employment       Category = "employment"
	EconomicActivity Category = "economic_activity"
	Inflation        Category = "inflation"
	Credit           Category = "credit"
	CentralBanks     Category = "central_banks"
	ConfidenceIndex  Category = "confidence_index"
	Balance          Category = "balance"
	Bonds            Category = "bonds"
)

type Importance string

const (
	High   Importance = "high"
	Medium Importance = "medium"
	Low    Importance = "low"
)

type DomainId int

const (
	RU DomainId = 7
	EN DomainId = 1
)

type Filter struct {
	DomainId   DomainId
	Limit      int
	StartDate  time.Time
	EndDate    time.Time
	Importance Importance
	Countries  []Country
	Categories []Category
	PageCursor string
}

func NewFilter() *Filter {
	return &Filter{
		Countries:  make([]Country, 0),
		Categories: make([]Category, 0),
		Limit:      100,
		DomainId:   EN,
	}
}

func (f *Filter) WithLimit(limit int) *Filter {
	f.Limit = limit
	return f
}

func (f *Filter) WithDomainId(domainId DomainId) *Filter {
	f.DomainId = domainId
	return f
}

func (f *Filter) WithCountries(countries ...Country) *Filter {
	f.Countries = countries
	return f
}

func (f *Filter) WithCategories(categories ...Category) *Filter {
	f.Categories = categories
	return f
}

func (f *Filter) WithStartDate(startDate time.Time) *Filter {
	f.StartDate = startDate
	return f
}

func (f *Filter) WithEndDate(endDate time.Time) *Filter {
	f.EndDate = endDate
	return f
}

func (f *Filter) SetToday() *Filter {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999, time.UTC)
	return f.WithStartDate(startDate).WithEndDate(endDate)
}

func (f *Filter) SetTomorrow() *Filter {
	now := time.Now().AddDate(0, 0, 1)
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999, time.UTC)
	return f.WithStartDate(startDate).WithEndDate(endDate)
}

func (f *Filter) String() string {
	var vars []string
	if f.DomainId != 0 {
		vars = append(vars, fmt.Sprintf("domain_id=%d", f.DomainId))
	}
	if f.Limit != 0 {
		vars = append(vars, fmt.Sprintf("limit=%d", f.Limit))
	}
	if !f.StartDate.IsZero() {
		vars = append(vars, fmt.Sprintf("start_date=%s", f.StartDate.Format("2006-01-02T15:04:05.000Z")))
	}
	if !f.EndDate.IsZero() {
		vars = append(vars, fmt.Sprintf("end_date=%s", f.EndDate.Format("2006-01-02T15:04:05.000Z")))
	}
	if f.Importance != "" {
		vars = append(vars, fmt.Sprintf("importance=%s", f.Importance))
	}
	if f.Countries != nil {
		countries := make([]string, 0)
		for _, country := range f.Countries {
			countries = append(countries, strconv.Itoa(int(country)))
		}
		if len(countries) > 0 {
			vars = append(vars, fmt.Sprintf("country_ids=%s", strings.Join(countries, ",")))
		}
	}
	if f.Categories != nil {
		categories := make([]string, 0)
		for _, category := range f.Categories {
			categories = append(categories, string(category))
		}
		if len(categories) > 0 {
			vars = append(vars, fmt.Sprintf("categories=%s", strings.Join(categories, ",")))
		}
	}
	if f.PageCursor != "" {
		vars = append(vars, fmt.Sprintf("page_cursor=%s", f.PageCursor))
	}
	if len(vars) > 0 {
		return "?" + strings.Join(vars, "&")
	}
	return ""
}
