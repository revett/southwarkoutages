package domain

import (
	"crypto/md5" //nolint:gosec
	"fmt"
	"time"

	"github.com/go-playground/validator"
	"github.com/pkg/errors"
)

// Outage represents information about a communal outage in Southwark.
type Outage struct {
	Location          string     `validate:"required"`
	Description       string     `validate:"required"`
	ReportedAt        time.Time  `validate:"required"`
	ReturnOfServiceAt *time.Time `validate:""`
	Status            string     `validate:"required"`
}

// Hash returns a consisnt hash for identifying the Outage, using ReportedAt and Location.
func (o Outage) Hash() string {
	hash := fmt.Sprintf("%s_%s", o.ReportedAt.Format("2006-01-02"), o.Location)
	return fmt.Sprintf("%x", md5.Sum([]byte(hash))) //nolint:gosec
}

// IsValid verifies that the data within the Outage looks correct, useful when scraping.
func (o Outage) IsValid() (bool, error) {
	if err := validator.New().Struct(o); err != nil {
		return false, errors.Wrap(err, "invalid outage data")
	}

	return true, nil
}
