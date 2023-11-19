package web

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// Fetch fetches HTML content from a URL with timeout.
func Fetch(ctx context.Context, uri string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return "", errors.Wrap(err, "creating http request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "fetching page")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(errors.Wrap(err, "closing response body"))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "reading response body of page")
	}

	return string(bodyBytes), nil
}
