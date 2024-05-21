package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func server(apiConfig *Config) {

	// Define CORS options
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"}, // You can customize this based on your needs
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // You can customize this based on your needs
		AllowCredentials: true,
		MaxAge:           300, // Maximum age for cache, in seconds
	}
	router := chi.NewRouter()
	apiRoute := chi.NewRouter()
	router.Use(cors.Handler(corsOptions))
	apiRoute.Get("/hello", helloReady)
	apiRoute.Get("/error", errorReady)
	// HANDLE POSTS
	apiRoute.Post("/posts", apiConfig.postPosttHandler)
	apiRoute.Get("/posts", apiConfig.getPostsHandler)
	apiRoute.Put("/post/{postID}", apiConfig.updatePosttHandler)
	apiRoute.Delete("/post/{postID}", apiConfig.deletePosttHandler)
	// HANDLE PLAYLISTS

	apiRoute.Post("/playlists", apiConfig.postPlaylistHandler)
	apiRoute.Get("/playlists", apiConfig.getPlaylistsHandler)
	// HANDLE TUTORIALS
	apiRoute.Post("/tutorials", apiConfig.postTutorialHandler)
	apiRoute.Get("/tutorials/{playlistID}", apiConfig.getPlaylistTutorials)
	apiRoute.Get("/tutorials", apiConfig.getTutorials)
	// HANLDE IMAGES
	apiRoute.Get("/images", apiConfig.getImagesHandler)
	apiRoute.Post("/images", apiConfig.postImageHandler)

	router.Mount("/api", apiRoute)
	router.Get("/", renderHome)
	srv := &http.Server{
		Addr:              ":" + apiConfig.PORT,
		Handler:           router,
		ReadHeaderTimeout: time.Minute,
	}
	log.Printf("serving server on port %v", apiConfig.PORT)

	log.Printf("Serving on port: %s\n", apiConfig.PORT)
	log.Fatal(srv.ListenAndServe())
}
