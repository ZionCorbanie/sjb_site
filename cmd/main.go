package main

import (
	"context"
	"errors"
	"sjb_site/internal/config"
	"sjb_site/internal/handlers"
	"sjb_site/internal/hash/passwordhash"
	database "sjb_site/internal/store/db"
	"sjb_site/internal/store/dbstore"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	m "sjb_site/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

/*
* Set to production at build time
* used to determine what assets to load
 */
var Environment = "development"

func init() {
	os.Setenv("env", Environment)
	// run generate script
	exec.Command("make", "tailwind-build").Run()
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	r := chi.NewRouter()

	cfg := config.MustLoadConfig()

	db := database.MustOpen(cfg.DatabaseName)
	passwordhash := passwordhash.NewHPasswordHash()

	userStore := dbstore.NewUserStore(
		dbstore.NewUserStoreParams{
			DB:           db,
			PasswordHash: passwordhash,
		},
	)

	groupStore := dbstore.NewGroupStore(
		dbstore.NewGroupStoreParams{
			DB:           db,
		},
	)

    groupUserStore := dbstore.NewGroupUserStore(
        dbstore.NewGroupUserStoreParams{
            DB: db,
        },
    )

	sessionStore := dbstore.NewSessionStore(
		dbstore.NewSessionStoreParams{
			DB: db,
		},
	)

	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	authMiddleware := m.NewAuthMiddleware(sessionStore, cfg.SessionCookieName)

	r.Group(func(r chi.Router) {
		r.Use(
			middleware.Logger,
			m.TextHTMLMiddleware,
			// m.CSPMiddleware,
			authMiddleware.AddUserToContext,
		)

		r.NotFound(handlers.NewNotFoundHandler().ServeHTTP)

		r.Get("/", handlers.NewHomeHandler().ServeHTTP)

        //Need to be logged in to access these routes
        r.Group(func(r chi.Router) {
            r.Use(authMiddleware.LoggedIn)
            r.Route("/webalmanak", func(r chi.Router) {
                r.Route("/leden", func(r chi.Router) {
                    r.Get("/", handlers.NewLedenSearchHandler().ServeHTTP)
                    r.Post("/", handlers.NewPostLedenSearchHandler(handlers.PostLedenSearchHandlerParams{
                        UserStore: userStore,
                    }).ServeHTTP)

                    r.Get("/{userId}", handlers.NewLedenHandler(handlers.GetLidHandlerParams{
                        UserStore: userStore,
                    }).ServeHTTP)
                    r.Get("/{userId}/edit", handlers.NewLidEditHandler(handlers.GetLidEditHandlerParams{
                        UserStore: userStore,
                    }).ServeHTTP)
                    r.Patch("/{userId}/edit", handlers.NewPatchtLidEditHandler(handlers.PatchLidEditHandlerParams{
                        UserStore: userStore,
                    }).ServeHTTP)
                })
                r.Route("/{groupType}", func(r chi.Router){
                    r.Get("/", handlers.NewGroupsHandler(handlers.GetGroupsHandlerParams{
                        GroupStore: groupStore,
                    }).ServeHTTP)
                })
                r.Get("/groep/{groupId}", handlers.NewGroupHandler(handlers.GetGroupHandlerParams{
                    GroupStore: groupStore,
                    GroupUserStore: groupUserStore,
                }).ServeHTTP)
            })
        })

        r.Route("/admin", func(r chi.Router) {
            r.Use(authMiddleware.IsAdmin)
            r.Get("/", handlers.NewAdminHandler().ServeHTTP)
        })

		r.Get("/about", handlers.NewAboutHandler().ServeHTTP)

		r.Get("/register", handlers.NewGetRegisterHandler().ServeHTTP)

		r.Post("/register", handlers.NewPostRegisterHandler(handlers.PostRegisterHandlerParams{
			UserStore: userStore,
		}).ServeHTTP)

		r.Get("/login", handlers.NewGetLoginHandler().ServeHTTP)

		r.Post("/login", handlers.NewPostLoginHandler(handlers.PostLoginHandlerParams{
			UserStore:         userStore,
			SessionStore:      sessionStore,
			PasswordHash:      passwordhash,
			SessionCookieName: cfg.SessionCookieName,
		}).ServeHTTP)

		r.Post("/logout", handlers.NewPostLogoutHandler(handlers.PostLogoutHandlerParams{
			SessionCookieName: cfg.SessionCookieName,
		}).ServeHTTP)
	})

	killSig := make(chan os.Signal, 1)

	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: r,
	}

	go func() {
		err := srv.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("Server shutdown complete")
		} else if err != nil {
			logger.Error("Server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	logger.Info("Server started", slog.String("port", cfg.Port), slog.String("env", Environment))
	<-killSig

	logger.Info("Shutting down server")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", slog.Any("err", err))
		os.Exit(1)
	}

	logger.Info("Server shutdown complete")
}
