package validation

import (
	"flag"
	"time"

	"github.com/cortexproject/cortex/pkg/util"
)

// Limits describe all the limits for users; can be used to describe global default
// limits via flags, or per-user limits via yaml config.
type Limits struct {
	// Distributor enforced limits.
	IngestionRate          float64       `yaml:"ingestion_rate"`
	IngestionBurstSize     int           `yaml:"ingestion_burst_size"`
	MaxLabelNameLength     int           `yaml:"max_label_name_length"`
	MaxLabelValueLength    int           `yaml:"max_label_value_length"`
	MaxLabelNamesPerSeries int           `yaml:"max_label_names_per_series"`
	RejectOldSamples       bool          `yaml:"reject_old_samples"`
	RejectOldSamplesMaxAge time.Duration `yaml:"reject_old_samples_max_age"`
	CreationGracePeriod    time.Duration `yaml:"creation_grace_period"`

	// Ingester enforces limits.
	MaxSeriesPerQuery  int `yaml:"max_series_per_query"`
	MaxSamplesPerQuery int `yaml:"max_samples_per_query"`
	MaxSeriesPerUser   int `yaml:"max_series_per_user"`
	MaxSeriesPerMetric int `yaml:"max_series_per_metric"`

	// Config for overrides, convenient if it goes here.
	PerTenantOverrideConfig string
	PerTenantOverridePeriod time.Duration
}

// RegisterFlags adds the flags required to config this to the given FlagSet
func (l *Limits) RegisterFlags(f *flag.FlagSet) {
	f.Float64Var(&l.IngestionRate, "distributor.ingestion-rate-limit", 25000, "Per-user ingestion rate limit in samples per second.")
	f.IntVar(&l.IngestionBurstSize, "distributor.ingestion-burst-size", 50000, "Per-user allowed ingestion burst size (in number of samples).")
	f.IntVar(&l.MaxLabelNameLength, "validation.max-length-label-name", 1024, "Maximum length accepted for label names")
	f.IntVar(&l.MaxLabelValueLength, "validation.max-length-label-value", 2048, "Maximum length accepted for label value. This setting also applies to the metric name")
	f.IntVar(&l.MaxLabelNamesPerSeries, "validation.max-label-names-per-series", 30, "Maximum number of label names per series.")
	f.BoolVar(&l.RejectOldSamples, "validation.reject-old-samples", false, "Reject old samples.")
	f.DurationVar(&l.RejectOldSamplesMaxAge, "validation.reject-old-samples.max-age", 14*24*time.Hour, "Maximum accepted sample age before rejecting.")
	f.DurationVar(&l.CreationGracePeriod, "validation.create-grace-period", 10*time.Minute, "Duration which table will be created/deleted before/after it's needed; we won't accept sample from before this time.")

	f.IntVar(&l.MaxSeriesPerQuery, "ingester.max-series-per-query", 100000, "The maximum number of series that a query can return.")
	f.IntVar(&l.MaxSamplesPerQuery, "ingester.max-samples-per-query", 1000000, "The maximum number of samples that a query can return.")
	f.IntVar(&l.MaxSeriesPerUser, "ingester.max-series-per-user", 5000000, "Maximum number of active series per user.")
	f.IntVar(&l.MaxSeriesPerMetric, "ingester.max-series-per-metric", 50000, "Maximum number of active series per metric name.")

	f.StringVar(&l.PerTenantOverrideConfig, "limits.per-user-override-config", "", "File name of per-user overrides.")
	f.DurationVar(&l.PerTenantOverridePeriod, "limits.per-user-override-period", 10*time.Second, "Period with this to reload the overrides.")
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (l *Limits) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// We want to set c to the defaults and then overwrite it with the input.
	// To make unmarshal fill the plain data struct rather than calling UnmarshalYAML
	// again, we have to hide it using a type indirection.  See prometheus/config.
	util.DefaultValues(l)
	type plain Limits
	return unmarshal((*plain)(l))
}
