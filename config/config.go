package config

import (
	"encoding/json"
	"os"

	"github.com/simon04/aur-out-of-date/upstream"
)

// Config contains options for running aur-out-of-date
type Config struct {
	Ignore map[string]([]upstream.Version) `json:"ignore"`
}

// FromFile reads the config from the given filename
func FromFile(filename string) (*Config, error) {
	var config Config
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return &config, nil
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(f)
	err = dec.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// IsIgnored determines whether the package in version is to be ignored
func (conf *Config) IsIgnored(pkg string, version upstream.Version) bool {
	ignoredVersions, ok := conf.Ignore[pkg]
	if !ok {
		return false
	}
	for _, v := range ignoredVersions {
		if v == "*" || v.String() == version.String() {
			return true
		}
	}
	return false
}
