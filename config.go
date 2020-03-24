package proxy

import (
	"io/ioutil"

	"launchpad.net/goyaml"
)

// ProxyConfig Type
type ProxyConfig struct {
	Bind         string    `yaml:"bind"`
	WaitQueueLen int       `yaml:"wait_queue_len"`
	MaxConn      int       `yaml:"max_conn"`
	Timeout      int       `yaml:"timeout"`
	FailOver     int       `yaml:"failover"`
	Backend      []string  `yaml:"backend"`
	Log          LogConfig `yaml:"log"`
	Stats        string    `yaml:"stats"`
}

// LogConfig Type
type LogConfig struct {
	Level       int8   `yaml:"level"`
	Path        string `yaml:"path"`
	MaxSize     int    `yaml:"max_size"`
	MaxBackup   int    `yaml:"max_backup"`
	MaxAge      int    `yaml:"max_age"`
	Compress    bool   `yaml:"compress"`
	ServiceName string `yaml:"servicename"`
}

func parseConfigFile(filepath string) error {
	if config, err := ioutil.ReadFile(filepath); err == nil {
		if err = goyaml.Unmarshal(config, &config); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}
