package steps

import (
	"time"

	"github.com/musicmash/artists/pkg/api"
	"github.com/musicmash/artists/pkg/api/artists"
	"github.com/pkg/errors"
)

// Updates ArtistName in whole []*Item
type ArtistStep struct {
	client *api.Provider
}

func NewSetArtistNameStep(client *api.Provider) *ArtistStep {
	return &ArtistStep{client: client}
}

func (a *ArtistStep) Do(_ time.Time, items []*Item) ([]*Item, error) {
	ids := make([]int64, len(items))
	for i, item := range items {
		ids[i] = item.ArtistID
	}

	artistsInfo, err := artists.GetFullInfo(a.client, ids)
	if err != nil {
		return nil, errors.Wrap(err, "tried to get artist details")
	}

	itemsMap := make(map[int64]*Item, len(items))
	for _, item := range items {
		itemsMap[item.ArtistID] = item
	}

	for _, artist := range artistsInfo {
		itemsMap[artist.ID].ArtistName = artist.Name
	}

	for i, item := range items {
		items[i] = itemsMap[item.ArtistID]
	}
	return items, nil
}
