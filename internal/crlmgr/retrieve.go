package crlmgr

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// ErrNotModified returned when a source URL status code is 304 not modified.
var ErrNotModified = errors.New("source is not modified")

// Retrieve the url and return a reader for the body along with the last modified time.
func Retrieve(ctx context.Context, u *url.URL, ts time.Time) (io.ReadCloser, time.Time, error) {
	c := &http.Client{
		Timeout: 10 * time.Second, //nolint:gomnd // TODO pull this from config.
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("unable to create request: %w", err)
	}

	if !ts.Equal(time.Time{}) {
		req.Header.Add("If-Modified-Since", ts.UTC().Format(http.TimeFormat))
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("unable to retrieve CRL: %w", err)
	}

	lm, _ := http.ParseTime(resp.Header.Get("Last-Modified"))

	if resp.StatusCode == http.StatusNotModified {
		return nil, lm, ErrNotModified
	}

	return resp.Body, lm, nil
}
