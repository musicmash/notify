package config

import (
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

var Config *AppConfig

type AppConfig struct {
	HTTP     HTTPConfig `yaml:"http"`
	DB       DBConfig   `yaml:"db"`
	Log      LogConfig  `yaml:"log"`
	Notifier Notifier   `yaml:"notifier"`
	Sentry   Sentry     `yaml:"sentry"`
}

type HTTPConfig struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

type LogConfig struct {
	File  string `yaml:"file"`
	Level string `yaml:"level"`
}

type DBConfig struct {
	Type  string `yaml:"type"`
	Host  string `yaml:"host"`
	Name  string `yaml:"name"`
	Login string `yaml:"login"`
	Pass  string `yaml:"pass"`
	Log   bool   `yaml:"log"`
}

type Notifier struct {
	TelegramToken       string  `yaml:"telegram_token"`
	CountOfSkippedHours float64 `yaml:"count_of_skipped_hours"`
}

type Sentry struct {
	Enabled bool   `yaml:"enabled"`
	Key     string `yaml:"key"`
}

func InitConfig(filepath string) error {
	data, err := ioutil.ReadFile(filepath)

	if err != nil {
		return err
	}

	if err := Load(data); err != nil {
		return err
	}

	log.Infof("Config loaded from %v.", filepath)
	return nil
}

func Load(data []byte) error {
	cfg := AppConfig{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	Config = &cfg
	return nil
}

func (db *DBConfig) GetConnString() (dialect, connString string) {
	if db.Type != "mysql" {
		panic("Only mysql is currently supported")
	}
	connString = fmt.Sprintf(
		"%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=UTC",
		db.Login, db.Pass, db.Host, db.Name)
	return db.Type, connString
}
