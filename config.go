package proxy

import (
	"io/ioutil"
	"os"

	"launchpad.net/goyaml"
)

var config ProxyConfig

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

// pathExists
// 判断文件路径是不是存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 解析配置文件
func parseConfigFile(filepath string) error {
	exist, err := pathExists(filepath)
	if exist {
		if configTemp, err := ioutil.ReadFile(filepath); err == nil {

			if err = goyaml.Unmarshal(configTemp, &config); err != nil {
				return err
			}
		} else {
			return err
		}
		return nil
	}

	if err == nil {
		err := os.Mkdir(filepath, os.ModePerm)
		if err != nil {
			logger.Error("mkdir failed![%s]" + err.Error())
			return err
		} else {
			logger.Error("mkdir success!")
			return nil
		}
	}

	return err
}
