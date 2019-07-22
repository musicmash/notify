package main

import (
	artsapi "github.com/musicmash/artists/pkg/api"
	mashapi "github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/notify/internal/config"
	"github.com/musicmash/notify/internal/notifier"
	subsapi "github.com/musicmash/subscriptions/pkg/api"
)

func makeNotifier() *notifier.Notifier {
	return notifier.New(
		mashapi.NewProvider(config.Config.Musicmash, 1),
		artsapi.NewProvider(config.Config.Artists, 1),
		subsapi.NewProvider(config.Config.Subscriptions, 1),
	)
}
