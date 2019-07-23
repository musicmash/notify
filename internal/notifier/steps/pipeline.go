package steps

import (
	"time"

	artsapi "github.com/musicmash/artists/pkg/api"
	mashapi "github.com/musicmash/musicmash/pkg/api"
	subsapi "github.com/musicmash/subscriptions/pkg/api"
)

type Pipeline struct {
	steps []Step
}

func NewPipeline(mashClient *mashapi.Provider, artsClient *artsapi.Provider, subsClient *subsapi.Provider) *Pipeline {
	return &Pipeline{
		steps: []Step{
			// fetch new releases
			NewFetchStep(mashClient),

			// find subscribers
			NewFindSubscribersStep(subsClient),

			// exclude subscribers that already received notifications
			NewExcludeSubscribersStep(),

			// find subscribers chats
			NewFindSubscriberChatsStep(),

			// set artist name
			NewSetArtistNameStep(artsClient),
		},
	}
}

func (p *Pipeline) Do(period time.Time) ([]*Item, error) {
	// items will be inited after first (fetch) step
	items := []*Item{}
	var err error
	for _, step := range p.steps {
		items, err = step.Do(period, items)
		if err != nil {
			return nil, err
		}
	}
	return items, nil
}
