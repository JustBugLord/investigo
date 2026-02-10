package investigo

import (
	"errors"
	"regexp"

	"github.com/JustBugLord/reqtango/v2"
)

var base = "https://www.investing.com/"
var regex = regexp.MustCompile("\"accessToken\":\"([^\"]+)\"")

type Investigo struct {
	rb    *reqtango.RequestBuilder
	token string
}

func NewInvestigo(defaultHeaders map[string]string) (*Investigo, error) {
	rb := reqtango.NewRequestBuilder(defaultHeaders)
	site, err := rb.Get(base)
	if err != nil {
		return nil, errors.New("fail get base site: " + err.Error())
	}

	var token string
	if len(site.Body) > 0 && regex.MatchString(string(site.Body)) {
		allMatches := regex.FindAllStringSubmatch(string(site.Body), -1)
		token = allMatches[0][1]
	}

	if token == "" {
		return nil, errors.New("fail parse token")
	}

	rb.SetHeaders(map[string]string{
		"Authorization": "Bearer " + token,
		"Referer":       "https://www.investigo.com/",
		"User-Agent":    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.7444.265 Safari/537.36",
	})

	return &Investigo{
		rb:    reqtango.NewRequestBuilder(defaultHeaders),
		token: token,
	}, nil
}

func NewInvestigoSimple() (*Investigo, error) {
	return NewInvestigo(nil)
}
