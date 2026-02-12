package socket

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var separator = ":%%"

type EventType string

const (
	Update        EventType = "update"
	Heartbeat     EventType = "heartbeat"
	BulkSubscribe EventType = "bulk-subscribe"
)

type WSRequest struct {
	Event   EventType `json:"_event"`
	Message string    `json:"message"`
	TzId    int       `json:"tzID"`
}

func (req *WSRequest) WithMessageFromEvents(occurenceIds []int) *WSRequest {
	events := make([]string, 0)
	for _, occurenceId := range occurenceIds {
		events = append(events, "event-"+strconv.Itoa(occurenceId))
	}
	msg := strings.Join(events, separator) + ":"
	req.Message = msg
	return req
}

func (req *WSRequest) String() string {
	return fmt.Sprintf("[\"{\\\"_event\\\":\\\"%s\\\",\\\"tzID\\\":%d,\\\"message\\\":\\\"%s\\\"}\"]", req.Event, req.TzId, req.Message)
}

type WSResponse struct {
	Event EventType `json:"event"`
	Data  string    `json:"data"`
}

func WSResponseFromRaw(rawData []byte) *WSResponse {
	if len(rawData) == 0 {
		return nil
	}
	return &WSResponse{
		Event: Update,
		Data:  string(rawData),
	}
}

func (resp *WSResponse) DataAsEconomicNewsUpdate() (*WSEconomicNewsUpdate, error) {
	cleaned := strings.ReplaceAll(resp.Data, `\"}"]`, ``)
	index := strings.Index(cleaned, "::")
	if index != -1 {
		cleaned = cleaned[index+2:]
	}
	cleaned = strings.ReplaceAll(cleaned, `\\\"`, `"`)
	cleaned = strings.ReplaceAll(cleaned, `\\\\"`, `"`)
	cleaned = strings.ReplaceAll(cleaned, `\"`, `"`)
	object := new(WSEconomicNewsUpdate)
	if err := json.Unmarshal([]byte(cleaned), object); err != nil {
		return nil, errors.New("error unmarshalling ws response: " + err.Error())
	}
	return object, nil
}

type Color string

const (
	Black Color = "blackFont"
	Green Color = "greenFont"
	Red   Color = "redFont"
)

func (c Color) PaintString(source string) string {
	var builder strings.Builder
	switch c {
	case Black:
		builder.WriteString("\033[36m")
	case Green:
		builder.WriteString("\033[32m")
	case Red:
		builder.WriteString("\033[31m")
	}
	builder.WriteString(source)
	builder.WriteString("\033[0m")
	return builder.String()
}

type WSEconomicNewsUpdate struct {
	EventID     string `json:"event_ID"`
	ActualColor Color  `json:"actual_color"`
	RevFromCol  Color  `json:"rev_from_col"`
	Previous    string `json:"previous"`
	Forecast    string `json:"forecast"`
	Actual      string `json:"actual"`
	RevFrom     string `json:"rev_from"`
}
