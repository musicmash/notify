package main

import (
	"flag"

	raven "github.com/getsentry/raven-go"
	"github.com/musicmash/notify/internal/api"
	"github.com/musicmash/notify/internal/config"
	"github.com/musicmash/notify/internal/cron"
	"github.com/musicmash/notify/internal/db"
	"github.com/musicmash/notify/internal/log"
	"github.com/musicmash/notify/internal/notifier"
	"github.com/musicmash/notify/internal/notifier/telegram"
	"github.com/pkg/errors"
)

func main() {
	configPath := flag.String("config", "/etc/musicmash/notify.yaml", "Path to notify.yaml config")
	flag.Parse()

	if err := config.InitConfig(*configPath); err != nil {
		panic(err)
	}
	if config.Config.Log.Level == "" {
		config.Config.Log.Level = "info"
	}

	log.SetLogFormatter(&log.DefaultFormatter)
	log.ConfigureStdLogger(config.Config.Log.Level)

	db.DbMgr = db.NewMainDatabaseMgr()
	telegram.New(config.Config.Notifier.TelegramToken)
	if config.Config.Sentry.Enabled {
		if err := raven.SetDSN(config.Config.Sentry.Key); err != nil {
			panic(errors.Wrap(err, "tried to setup sentry client"))
		}
	}

	log.Info("Running musicmash..")
	go cron.Run(db.ActionNotify, config.Config.Notifier.CountOfSkippedHours, notifier.Notify)
	log.Panic(api.ListenAndServe(config.Config.HTTP.IP, config.Config.HTTP.Port))
}
