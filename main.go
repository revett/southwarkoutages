package main

import (
	"context"

	commonlog "github.com/revett/common/log"
	"github.com/revett/southwarkoutages/internal/web"
	"github.com/rs/zerolog/log"
)

const uri = "https://www.southwark.gov.uk/housing/repairs/communal-breakdowns"

func main() {
	log.Logger = commonlog.New()

	ctx := context.Background()

	html, err := web.Fetch(ctx, uri)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	outages, err := web.Scrape(ctx, html)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Info().Msgf("found %d outages", len(outages))

	for _, outage := range outages {
		fields := map[string]any{
			"reportedAt": outage.ReportedAt.Format("2006-01-02"),
		}

		if outage.ReturnOfServiceAt != nil {
			fields["returnOfServiceAt"] = outage.ReturnOfServiceAt.Format("2006-01-02")
		}

		log.Info().Fields(fields).Msg(outage.Location)
	}
}
