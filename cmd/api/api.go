package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/robertgouveia/social/docs"
	"github.com/robertgouveia/social/internal/auth"
	"github.com/robertgouveia/social/internal/mail"
	"github.com/robertgouveia/social/internal/store"
	"github.com/robertgouveia/social/internal/store/cache"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

// DI - Dependency Injection
type application struct {
	config        config
	store         store.Storage
	logger        *zap.SugaredLogger
	mailer        *mail.MailHog
	authenticator auth.Authenticator
	cacheStorage  cache.Storage
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type basicConfig struct {
	user     string
	password string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	mail        mailConfig
	frontendURL string
	auth        authConfig
	redisCfg    redisCfg
}

type redisCfg struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

type mailConfig struct {
	exp       time.Duration
	fromEmail string
	sendGrid  sendGridConfig
}

type sendGridConfig struct {
	apiKey string
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.With(app.BasicAuthMiddleware()).Get("/health", app.healthCheckHandler)

		// multiple documents, configuration
		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)

				r.Get("/", app.checkPostOwnership("moderator", app.getPostHandler))
				r.Delete("/", app.checkPostOwnership("admin", app.deletePostHandler))
				r.Patch("/", app.updatePostHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)

			r.With(app.AuthTokenMiddleware).Route("/{userID}", func(r chi.Router) {
				r.Get("/", app.getUserHandler)

				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
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
