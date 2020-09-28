package config

import "github.com/jinzhu/configor"

// LoadConfigs conducts sequential of files or it will return an error if something goes wrong.
// See configor.Load for more details.
func LoadConfigs(config interface{}, files ...string) error {
	for _, file := range files {
		if err := configor.Load(config, file); err != nil {
			return err
		}
	}
	return nil
}
