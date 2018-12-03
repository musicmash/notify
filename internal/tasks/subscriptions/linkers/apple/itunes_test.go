package apple

import (
	"net/http"
	"testing"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_AppleLinker_Reserve(t *testing.T) {
	task := NewLinker("http://url.mock", "xxx")

	// action
	task.reserveArtists([]string{testutil.ArtistSkrillex, testutil.ArtistArchitects})

	// assert
	assert.Len(t, task.reservedArtists, 2)
}

func Test_AppleLinker_Free(t *testing.T) {
	// arrange
	task := NewLinker("http://url.mock", "xxx")
	artists := []string{testutil.ArtistSkrillex, testutil.ArtistArchitects}
	task.reserveArtists(artists)
	assert.Len(t, task.reservedArtists, 2)

	// action
	task.freeReservedArtists(artists)

	// assert
	assert.Len(t, task.reservedArtists, 0)
}

func Test_AppleLinker_Search_AlreadyExists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	task := NewLinker("http://url.mock", "xxx")
	artists := []string{testutil.ArtistSkrillex, testutil.ArtistArchitects}
	assert.NoError(t, db.DbMgr.EnsureArtistExistsInStore(testutil.ArtistSkrillex, testutil.StoreApple, "xyz"))
	assert.NoError(t, db.DbMgr.EnsureArtistExistsInStore(testutil.ArtistArchitects, testutil.StoreApple, "zyx"))

	// action
	task.SearchArtists(artists)
}

func Test_AppleLinker_Search(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	task := NewLinker(server.URL, "xxx")
	mux.HandleFunc("/v1/catalog/us/search", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
{
  "results": {
    "artists": {
      "data": [
        {
          "attributes": {
            "name": "Architects"
          },
          "id": "182821355"
        }
      ]
    }
  }
}
		`))
	})

	// action
	task.SearchArtists([]string{testutil.ArtistArchitects})

	// assert
	artists, err := db.DbMgr.GetArtistsForStore(testutil.StoreApple)
	assert.NoError(t, err)
	assert.Len(t, artists, 1)
}
