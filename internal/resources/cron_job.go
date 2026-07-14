package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_cron_job",
		Filename: "cron_job.tf",
		Fetch:    fetchCronJob,
	})
}

func fetchCronJob(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/cron/settings/searchJobs")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Cron().GetJob(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("cron_job %s: %w", uuid, err)
		}
		label := tracker.Derive("", d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_cron_job" %s {
  enabled     = %s
  minutes     = %s
  hours       = %s
  days        = %s
  months      = %s
  weekdays    = %s
  who         = %s
  command     = %s
  parameters  = %s
  description = %s
}
`, hclString(label),
			hclBool(stringToBool(d.Enabled)),
			hclString(d.Minutes),
			hclString(d.Hours),
			hclString(d.Days),
			hclString(d.Months),
			hclString(d.Weekdays),
			hclString(d.Who),
			hclString(d.Command.String()),
			hclStringOrNull(d.Parameters),
			hclString(d.Description))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}
