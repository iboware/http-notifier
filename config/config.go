package config

import (
	"github.com/spf13/pflag"
)

type Config struct {
	Interval int
	Url      string
}

// RegisterFlags adds the configuration flags to the given flag set.
func (c *Config) RegisterFlags(f *pflag.FlagSet) {
	f.IntVarP(&c.Interval, "interval", "i", 5, "Notification Interval")
	f.StringVarP(&c.Url, "url", "u", "", "Notification URL")
}
