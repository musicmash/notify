package steps

import (
	"time"

	"github.com/musicmash/subscriptions/pkg/api"
	"github.com/musicmash/subscriptions/pkg/api/subscriptions"
	"github.com/pkg/errors"
)

// Finds subscribers for []*Item.ArtistID and puts their into []*Item.Subscribers
type SubscriptionStep struct {
	client *api.Provider
}

func NewFindSubscribersStep(client *api.Provider) *SubscriptionStep {
	return &SubscriptionStep{client: client}
}

func (s *SubscriptionStep) Do(period time.Time, items []*Item) ([]*Item, error) {
	ids := make([]int64, len(items))
	for i, item := range items {
		ids[i] = item.ArtistID
	}

	subs, err := subscriptions.GetArtistsSubscriptions(s.client, ids)
	if err != nil {
		return nil, errors.Wrapf(err, "tried to get releases since %s", period)
	}

	itemsMap := make(map[int64]*Item, len(items))
	for _, item := range items {
		itemsMap[item.ArtistID] = item
	}

	for _, subscription := range subs {
		if itemsMap[subscription.ArtistID].Subscribers == nil {
			itemsMap[subscription.ArtistID].Subscribers = []string{subscription.UserName}
			continue
		}

		itemsMap[subscription.ArtistID].Subscribers = append(
			itemsMap[subscription.ArtistID].Subscribers, subscription.UserName,
		)
	}

	for i, item := range items {
		items[i] = itemsMap[item.ArtistID]
	}
	return items, nil
}
