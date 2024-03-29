package scrapper

import (
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/unluckythoughts/go-microservice/tools/web"
)

type APIQueryData struct {
	URL         string
	Method      string
	Body        []byte
	Headers     http.Header
	PageParam   string
	QueryParams map[string]string
	Response    interface{}
	HasNextPage HasNextPage
	Transport   http.RoundTripper
}

type HasNextPage func(resp interface{}) bool

func GetAPIResponse(ctx web.Context, q APIQueryData) error {
	c := web.NewClientWithTransport(q.URL, q.Transport, q.Headers)
	params := url.Values{}
	for k, v := range q.QueryParams {
		params.Add(k, v)
	}

	path := "/"
	if len(params) > 0 {
		path = "?" + params.Encode()
	}
	if q.Method == "" {
		q.Method = http.MethodGet
	}

	status, err := c.Send(q.Method, path, q.Body, q.Response)
	if err != nil {
		return errors.Wrapf(err, "error while get data from %s", q.URL+path)
	}

	if status != 200 {
		return errors.Errorf("unexpected status %d when get data from %s", status, q.URL+path)
	}

	return nil
}
