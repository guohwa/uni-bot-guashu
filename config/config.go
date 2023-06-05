package config

import (
	"encoding/json"
	"os"

	"bot/log"
)

const (
	FILE = "config.json"
)

type database struct {
	Host string `json:"host"`
	Name string `json:"name"`
}

type app struct {
	Title     string   `json:"title"`
	Mode      string   `json:"mode"`
	Listen    string   `json:"listen"`
	Trust     []string `json:"trust"`
	WhiteList string   `json:"whiteList"`
}

type config struct {
	App      app      `json:"app"`
	Database database `json:"database"`
}

var (
	App      = app{}
	Database = database{}
)

func init() {
	var c = config{}
	if err := c.Load(); err != nil {
		log.Fatal(err)
	}
	App = c.App
	Database = c.Database
}

func (c *config) Load() error {
	body, err := os.ReadFile(FILE)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, c)
	return err
}

func (c *config) Save() error {
	body, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(FILE, body, 0644)
}
