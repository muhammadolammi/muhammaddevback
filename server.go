package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Middleware to check for the API key in the authorization header for all POST, PUT, DELETE, and OPTIONS requests
func apiKeyAuth(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" || r.Method == "OPTIONS" {
				authHeader := r.Header.Get("Authorization")
				if authHeader != apiKey {
					http.Error(w, "Action Not Permitted ", http.StatusForbidden)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
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
	// ADD MIDDLREWARE
	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(cors.Handler(corsOptions))
	router.Use(apiKeyAuth(apiConfig.API_KEY))

	// ADD ROUTES
	apiRoute.Get("/hello", helloReady)
	apiRoute.Get("/error", errorReady)
	// HANDLE POSTS
	apiRoute.Post("/posts", apiConfig.postPosttHandler)
	apiRoute.Get("/posts", apiConfig.getPostsHandler)
	apiRoute.Get("/post/{postID}", apiConfig.getPostWithIdHandler)
	apiRoute.Get("/getpost/{postTitle}", apiConfig.getPostWithTitleHandler)

	apiRoute.Put("/post/{postID}", apiConfig.updatePostHandler)
	apiRoute.Delete("/post/{postID}", apiConfig.deletePostHandler)

	// HANDLE PLAYLISTS

	apiRoute.Post("/playlists", apiConfig.postPlaylistHandler)
	apiRoute.Get("/playlists", apiConfig.getPlaylistsHandler)
	// HANDLE TUTORIALS
	apiRoute.Post("/tutorials", apiConfig.postTutorialHandler)
	apiRoute.Get("/tutorials/{playlistID}", apiConfig.getPlaylistTutorialsHandler)
	apiRoute.Get("/tutorials", apiConfig.getTutorialsHandler)
	apiRoute.Get("/gettutorial/{tutorialTitle}", apiConfig.getTutorialWithTitleHandler)
	apiRoute.Get("/tutorial/{tutorialID}", apiConfig.getTutorialWithIdHandler)
	apiRoute.Delete("/tutorial/{tutorialID}", apiConfig.deleteTutorialHandler)
	apiRoute.Put("/tutorial/{tutorialID}", apiConfig.updateTutorialHandler)

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
