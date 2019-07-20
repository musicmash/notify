package steps

import (
	"time"

	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/musicmash/notify/internal/db"
)

type Step interface {
	Do(period time.Time, items []*Item) ([]*Item, error)
}

// ArtistID id of artist equals to id from service index
// ArtistName name of artist equals name from artists service index
// Subscribers users that must recieve notification
// Chats array of subscribers chat_ids
// Releases releases from musicmash service that must be delivered to user
type Item struct {
	ArtistID    int64
	ArtistName  string
	Subscribers []string
	Chats       []*db.Chat
	Releases    []*releases.Release
}
