package steps

import (
	"time"

	"github.com/musicmash/notify/internal/db"
	"github.com/pkg/errors"
)

type ExcludeStep struct{}

func NewExcludeSubscribersStep() *ExcludeStep {
	return &ExcludeStep{}
}

func excludeSubscribers(allSubscribers, usersReceivedNotification []string) []string {
	receivedNotificationUsers := map[string]int8{}
	for _, user := range usersReceivedNotification {
		receivedNotificationUsers[user] = 1
	}

	result := []string{}
	for _, subscriber := range allSubscribers {
		if _, received := receivedNotificationUsers[subscriber]; received {
			continue
		}

		result = append(result, subscriber)
	}
	return result
}

func (s *ExcludeStep) Do(_ time.Time, items []*Item) ([]*Item, error) {
	result := []*Item{}
	for _, item := range items {
		for _, release := range item.Releases {
			usersReceivedNotification, err := db.DbMgr.FindUsersThatReceivedNotification(release.ID, item.Subscribers)
			if err != nil {
				return nil, errors.Wrapf(err, "tried to find users that received notification with id %d", release.ID)
			}

			item.Subscribers = excludeSubscribers(item.Subscribers, usersReceivedNotification)
			if len(item.Subscribers) == 0 {
				continue
			}

			result = append(result, item)
		}
	}

	return result, nil
}
