package steps

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	mashapi "github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/musicmash/notify/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNotifierSteps_Fetch(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	client := mashapi.NewProvider(server.URL, 1)
	step := NewFetchStep(client)
	mux.HandleFunc("/v1/releases", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		newReleases := []*releases.Release{
			{
				ArtistID:  testutil.StoreIDQ,
				Title:     testutil.ReleaseArchitectsHollyHell,
				Poster:    testutil.PosterSimple,
				StoreName: "itunes",
				StoreID:   "1473811207",
			},
		}

		b, _ := json.Marshal(&newReleases)
		_, _ = w.Write(b)
	})

	// action
	items, err := step.Do(time.Now(), nil)

	// assert
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, int64(testutil.StoreIDQ), items[0].ArtistID)
	assert.Empty(t, items[0].ArtistName)
	assert.Empty(t, items[0].Subscribers)
	assert.Empty(t, items[0].Chats)
	assert.Len(t, items[0].Releases, 1)
	assert.Equal(t, int64(testutil.StoreIDQ), items[0].ArtistID)
	assert.Equal(t, testutil.ReleaseArchitectsHollyHell, items[0].Releases[0].Title)
	assert.Equal(t, testutil.PosterSimple, items[0].Releases[0].Poster)
	assert.Equal(t, "itunes", items[0].Releases[0].StoreName)
	assert.Equal(t, "1473811207", items[0].Releases[0].StoreID)
}
