package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
		log.Fatal(err)
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
