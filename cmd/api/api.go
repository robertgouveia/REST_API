package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/robertgouveia/social/docs"
	"github.com/robertgouveia/social/internal/store"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

// DI - Dependency Injection
type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type config struct {
	addr   string
	db     dbConfig
	env    string
	apiURL string
	mail   mailConfig
}

type mailConfig struct {
	exp time.Duration
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		// multiple documents, configuration
		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)

				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)

				r.Get("/", app.getUserHandler)

				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})
			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}

	app.logger.Infow("Server has started", "addr", app.config.addr, "env", app.config.env)
	return srv.ListenAndServe()
}
