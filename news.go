package investigo

import (
	"errors"
	"fmt"
	"time"
)

type Instrument struct {
	Id         int  `json:"id"`
	PrimaryTag bool `json:"primary_tag"`
}

type Media struct {
	Caption   string `json:"caption"`
	Copyright string `json:"copyright"`
	Purpose   string `json:"purpose"`
	Url       string `json:"url"`
}

type Article struct {
	ArticleType         string       `json:"article_type"`
	AutomatedType       string       `json:"automated_type"`
	Body                string       `json:"body"`
	CategoryIds         []int        `json:"category_ids"`
	CountryIds          []int        `json:"country_ids"`
	CreatedBy           int          `json:"created_by"`
	DomainId            int          `json:"domain_id"`
	EconomicEventIds    []int        `json:"economic_event_ids"`
	EditedBy            int          `json:"edited_by"`
	EditorsPick         bool         `json:"editors_pick"`
	Featured            bool         `json:"featured"`
	Id                  int          `json:"id"`
	Important           bool         `json:"important"`
	Instruments         []Instrument `json:"instruments"`
	Keywords            string       `json:"keywords"`
	Link                string       `json:"link"`
	Media               []Media      `json:"media"`
	MultiPublishHeaders []string     `json:"multi_publish_headers"`
	NewsType            string       `json:"news_type"`
	PublishCompanyId    int          `json:"publish_company_id"`
	PublishedAt         time.Time    `json:"published_at"`
	SourceExternalLink  string       `json:"source_external_link"`
	SourceId            int          `json:"source_id"`
	SourceName          string       `json:"source_name"`
	Title               string       `json:"title"`
	UpdatedAt           time.Time    `json:"updated_at"`
}

type BreakingNewsResponse struct {
	Articles       []Article `json:"articles"`
	NextPageCursor string    `json:"next_page_cursor"`
	Total          int       `json:"total"`
}

func (i *Investigo) BreakingNews(filter *Filter) (*BreakingNewsResponse, error) {
	domainId := 1
	if filter.DomainId != 0 {
		domainId = int(filter.DomainId)
	}
	limit := 100
	if filter.Limit != 0 {
		limit = filter.Limit
	}
	pageCursor := ""
	if filter.PageCursor != "" {
		pageCursor = "&page_cursor=" + filter.PageCursor
	}
	result := new(BreakingNewsResponse)
	_, err := i.rb.GetToStruct(fmt.Sprintf("https://endpoints.investing.com/news-delivery/api/v2/articles/delivery/domains/%d/news/lists/breaking-news?limit=%d%s", domainId, limit, pageCursor), result)
	if err != nil {
		return nil, errors.New("fail to get economic calendar: " + err.Error())
	}
	return result, nil
}
