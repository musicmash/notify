package steps

import (
	"time"

	mashapi "github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/pkg/errors"
)

// Fetches new releases from musicmash and prepares []*Item
type FetchStep struct {
	client *mashapi.Provider
}

func NewFetchStep(client *mashapi.Provider) *FetchStep {
	return &FetchStep{client: client}
}

func (f *FetchStep) Do(period time.Time, items []*Item) ([]*Item, error) {
	newReleases, err := releases.Get(f.client, period)
	if err != nil {
		return nil, errors.Wrapf(err, "tried to get releases since %s", period.String())
	}

	return f.groupReleasesByArtist(newReleases), nil
}

func (f *FetchStep) groupReleasesByArtist(newReleases []*releases.Release) []*Item {
	groupedReleases := map[int64]*Item{}
	for _, release := range newReleases {
		if _, exist := groupedReleases[release.ArtistID]; !exist {
			groupedReleases[release.ArtistID] = &Item{
				ArtistID: release.ArtistID,
				Releases: []*releases.Release{release},
			}
			continue
		}

		groupedReleases[release.ArtistID].Releases = append(
			groupedReleases[release.ArtistID].Releases, release,
		)
	}

	i := 0
	items := make([]*Item, len(groupedReleases))
	for _, item := range groupedReleases {
		items[i] = item
		i++
	}
	return items
}
