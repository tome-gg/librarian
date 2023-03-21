package pkg

import (
	"fmt"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

// DSUReport ...
type DSUReport struct {
	ID          string `yaml:"id"`
	DatetimeRaw string `yaml:"datetime"`
	Datetime 		time.Time `yaml:"datetime_value"`
	Remarks     string `yaml:"remarks"`
	DoneYesterday   string `yaml:"done_yesterday"`
	DoingToday      string `yaml:"doing_today"`
	Blockers        string `yaml:"blockers"`
}

// UnmarshalYAML unmarshals the data from a YAML byte data.
func (r *DSUReport) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type dsuReportAlias DSUReport
	alias := dsuReportAlias(*r)
	if err := unmarshal(&alias); err != nil {
		return err
	}
	*r = DSUReport(alias)

	if dateStr := strings.TrimSpace(alias.DatetimeRaw); dateStr != "" {
		dateTime, err := dateparse.ParseAny(dateStr)
		if err != nil {
			return fmt.Errorf("failed to parse date: %s", err)
		}
		r.Datetime = dateTime
	}

	return nil
}