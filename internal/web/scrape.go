package web

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/revett/southwarkoutages/internal/domain"
	"github.com/rs/zerolog/log"
)

// Scrape extracts outage information from the given HTML.
func Scrape(ctx context.Context, html string) ([]*domain.Outage, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(html))
	if err != nil {
		return nil, errors.Wrap(err, "creating goquery document from reader")
	}

	outages := []*domain.Outage{}
	var scrapeErr error

	doc.Find("table.table tbody tr").Each(func(_ int, selection *goquery.Selection) {
		if scrapeErr != nil {
			return
		}

		outage, err := extractDataFromSelection(selection)
		if err != nil {
			scrapeErr = err
			return
		}

		outages = append(outages, outage)
	})
	if scrapeErr != nil {
		return nil, scrapeErr
	}

	validOutages := []*domain.Outage{}

	for _, outage := range outages {
		ok, err := outage.IsValid()
		if ok {
			validOutages = append(validOutages, outage)
		}

		if !ok {
			log.Warn().Err(err).Msgf("invalid outage: %+v", outage)
		}
	}

	return validOutages, nil
}

func extractDataFromSelection(selection *goquery.Selection) (*domain.Outage, error) {
	reportedAt, err := parseTime(selection.Find("td:nth-child(3)").Text())
	if err != nil {
		return nil, errors.Wrap(err, "parsing reported at string as time")
	}

	returnOfServiceAt, err := parseTime(selection.Find("td:nth-child(4)").Text())
	if err != nil {
		return nil, errors.Wrap(err, "parsing estimated return to service string as time")
	}

	return &domain.Outage{
		Location:          removeNewlines(selection.Find("td:nth-child(1)").Text()),
		Description:       selection.Find("td:nth-child(2)").Text(),
		ReportedAt:        *reportedAt,
		ReturnOfServiceAt: returnOfServiceAt,
		Status:            selection.Find("td:nth-child(5)").Text(),
	}, nil
}

func parseTime(timeStr string) (*time.Time, error) {
	if timeStr == "TBC" {
		return nil, nil //nolint:nilnil
	}

	timestamp, err := time.Parse("02/01/2006", timeStr)
	if err != nil {
		return nil, errors.Wrap(err, "parsing time string")
	}

	return &timestamp, nil
}

func removeNewlines(input string) string {
	parts := strings.Split(input, "\n")

	for i, part := range parts {
		if strings.TrimSpace(part) == "" {
			continue
		}

		parts[i] = strings.TrimSpace(part)
	}

	return strings.Join(parts, ", ")
}
