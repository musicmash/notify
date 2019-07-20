package steps

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/musicmash/notify/internal/db"
)

// Finds subscriber chat and puts it into []*Item.Chats
type ChatStep struct{}

func NewFindSubscriberChatsStep() *ChatStep {
	return &ChatStep{}
}

func removeAt(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func (c *ChatStep) Do(_ time.Time, items []*Item) ([]*Item, error) {
	for i, _ := range items {
		for subscriberID, subscriber := range items[i].Subscribers {
			chatID, err := db.DbMgr.FindChatByUserName(subscriber)
			if err != nil {
				if !gorm.IsRecordNotFoundError(err) {
					return nil, err
				}

				// handle chat for user not found error
				items[i].Subscribers = removeAt(items[i].Subscribers, subscriberID)
				continue
			}

			if items[i].Chats == nil {
				items[i].Chats = []*db.Chat{{UserName: subscriber, ID: *chatID}}
				continue
			}

			items[i].Chats = append(
				items[i].Chats, &db.Chat{UserName: subscriber, ID: *chatID},
			)
		}
	}
	return items, nil
}
