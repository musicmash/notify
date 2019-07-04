package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/musicmash/notify/internal/log"
)

func getMux() *chi.Mux {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/{user_name}/chats", addUserChat)
	r.Post("/users", createUser)
	r.Post("/{user_name}/chats", addUserChat)
	r.Get("/{user_name}/feed", getUserFeed)
	r.Post("/{user_name}/subscriptions", createSubscriptions)
	r.Delete("/{user_name}/subscriptions", deleteSubscriptions)
	r.Get("/{user_name}/subscriptions", getUserSubscriptions)
	r.Get("/{user_name}/artists", searchArtist)
	r.Get("/{user_name}/artists/{artist_name}", getArtistDetails)
	r.Get("/users/{user_name}", getUser)
	r.Get("/healthz", healthz)
	return r
}

func ListenAndServe(ip string, port int) error {
	addr := fmt.Sprintf("%s:%d", ip, port)
	log.Infof("Listening API on '%s'", addr)
	return http.ListenAndServe(addr, getMux())
}
