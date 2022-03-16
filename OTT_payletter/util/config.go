package util

import (
	"encoding/json"
	"io/ioutil"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	systemEnv map[string]string
}

var ServerConfig = Config{}

func (c *Config) LoadConfig() bool {
	file, err := os.Open("util/setting.json")
	if err != nil {
		return false
	}
	defer file.Close()
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return false
	}
	err = json.Unmarshal(byteValue, &c.systemEnv)
	return err == nil
}

func (c *Config) GetStringData(key string) string {
	val, exists := c.systemEnv[key]
	if !exists {
		return ""
	}
	return val
}

// systemEnv의 data 조회
func (c *Config) GetData() map[string]string {
	retData := make(map[string]string)
	for key, val := range c.systemEnv {
		retData[key] = val
	}
	return retData
}
