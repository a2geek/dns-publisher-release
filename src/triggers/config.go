package triggers

import "time"

type TriggerConfig struct {
	FileWatcher string // preferential, defaults to '/var/vcap/instance/dns/records.json'
	Refresh     string // use duration notation
}

func (c *TriggerConfig) Validate() error {
	if c.Refresh != "" {
		_, err := time.ParseDuration(c.Refresh)
		if err != nil {
			return err
		}
	}
	return nil
}
