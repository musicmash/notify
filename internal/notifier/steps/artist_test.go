package steps

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	artsapi "github.com/musicmash/artists/pkg/api"
	"github.com/musicmash/notify/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNotifierSteps_Artist(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	items := []*Item{{ArtistID: testutil.StoreIDQ}, {ArtistID: testutil.StoreIDW}}
	client := artsapi.NewProvider(server.URL, 1)
	step := NewSetArtistNameStep(client)
	mux.HandleFunc("/v1/artists", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		_, _ = fmt.Fprintf(w,
			`[{"id":%d,"name":"%s"},{"id":%d,"name":"%s"}]`,
			testutil.StoreIDQ, testutil.ArtistSkrillex,
			testutil.StoreIDW, testutil.ArtistArchitects)
	})

	// action
	items, err := step.Do(time.Now(), items)

	// assert
	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, int64(testutil.StoreIDQ), items[0].ArtistID)
	assert.Equal(t, testutil.ArtistSkrillex, items[0].ArtistName)
	assert.Empty(t, items[0].Chats)
	assert.Empty(t, items[0].Subscribers)
	assert.Empty(t, items[0].Releases)
	assert.Equal(t, int64(testutil.StoreIDW), items[1].ArtistID)
	assert.Equal(t, testutil.ArtistArchitects, items[1].ArtistName)
	assert.Empty(t, items[1].Chats)
	assert.Empty(t, items[1].Subscribers)
	assert.Empty(t, items[1].Releases)
}
