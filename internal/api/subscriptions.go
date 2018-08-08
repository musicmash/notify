package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/log"
	tasks "github.com/objque/musicmash/internal/tasks/subscriptions"
)

func validateUser(userID string, w http.ResponseWriter) error {
	_, err := db.DbMgr.FindUserByID(userID)
	return err
}

func createSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	if err := validateUser(userID, w); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
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