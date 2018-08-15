package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/log"
	tasks "github.com/objque/musicmash/internal/tasks/subscriptions"
	"github.com/pkg/errors"
	"github.com/objque/musicmash/internal/api/validators"
)

func createSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	if err := validators.IsUserExits(w, userID); err != nil {
		return
	}

	userArtists := []string{}
	if err := json.NewDecoder(r.Body).Decode(&userArtists); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, stateID := tasks.FindArtistsAndSubscribeUserTask(userID, userArtists)
	body := map[string]interface{}{
		"state_id": stateID,
	}
	buffer, _ := json.Marshal(&body)

	w.WriteHeader(http.StatusAccepted)
	w.Write(buffer)
}

func deleteSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	if err := validators.IsUserExits(w, userID); err != nil {
		return
	}

	artists := []string{}
	if err := json.NewDecoder(r.Body).Decode(&artists); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := db.DbMgr.UnsubscribeUserFromArtists(userID, artists); err != nil {
		log.Error(errors.Wrapf(err, "tried to unsubscribe user '%s' from artists '%v'", userID, artists))
	}

	w.WriteHeader(http.StatusOK)
}
