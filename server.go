package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Middleware to check for the API key in the authorization header for all POST, PUT, DELETE, and OPTIONS requests
func (apiConfig *Config) userAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken := r.Header.Get("Authorization")
			// Check if the request is to get user
			if r.Method == "GET" && r.URL.Path == "/users/me" && accessToken == apiConfig.API_KEY {
				// Allow the request if the admin API key is valid
				next.ServeHTTP(w, r)
				return
			}

			if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" || r.Method == "OPTIONS" {

				if accessToken == "" {
					respondWithError(w, http.StatusForbidden, "No access token")
					return
				}
				// Check if the request is for creating a user
				if r.URL.Path == "/api/signup" && accessToken == apiConfig.API_KEY {
					// Allow the request if the admin API key is valid
					next.ServeHTTP(w, r)
					return
				}
				// Check if the request is for creating a login
				if r.URL.Path == "/api/signin" && accessToken == apiConfig.API_KEY {
					// Allow the request if the admin API key is valid
					next.ServeHTTP(w, r)
					return
				}
				// Check if the request is refreshing
				if r.URL.Path == "/api/refresh" && accessToken == apiConfig.API_KEY {
					// Allow the request if the admin API key is valid
					next.ServeHTTP(w, r)
					return
				}

				// Check for normal users
				accessTokenExit, err := apiConfig.DB.AccessTokenExists(r.Context(), sql.NullString{
					Valid:  true,
					String: accessToken,
				})
				if err != nil {
					respondWithError(w, http.StatusForbidden, fmt.Sprintf("Error validating access token err: %v", err))
					return
				}
				if !accessTokenExit {
					respondWithError(w, http.StatusForbidden, "Invalid access token")
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
		AllowedOrigins: []string{"http://localhost:3000", "https://muhammaddev.com", "http://192.168.246.175:3000"}, // You can customize this based on your needs

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
	router.Use(apiConfig.userAuth())

	// ADD ROUTES
	apiRoute.Get("/hello", helloReady)
	apiRoute.Get("/error", errorReady)

	// Handle Auth
	apiRoute.Post("/signup", apiConfig.signupHandler)
	apiRoute.Post("/signin", apiConfig.signinHandler)
	apiRoute.Post("/refresh", apiConfig.refreshTokens)
	apiRoute.Post("/validate", apiConfig.validate)

	apiRoute.Put("/user", apiConfig.passwordChangeHandler)
	apiRoute.Get("/users/me", apiConfig.getUserHandler)

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

	router.Mount("/api", apiRoute)
	router.Get("/", renderHome)
	srv := &http.Server{
		Addr:              ":" + apiConfig.PORT,
		Handler:           router,
		ReadHeaderTimeout: time.Minute,
	}

	log.Printf("Serving on port: %s\n", apiConfig.PORT)
	log.Fatal(srv.ListenAndServe())
}
