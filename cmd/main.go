package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sjb_site/internal/config"
	"sjb_site/internal/handlers"
	"sjb_site/internal/handlers/admin"
	"sjb_site/internal/hash/passwordhash"
	database "sjb_site/internal/store/db"
	"sjb_site/internal/store/dbstore"
	"syscall"
	"time"

	m "sjb_site/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*
* Set to production at build time
* used to determine what assets to load
 */

func init() {
	os.Setenv("env", os.Getenv("ENVIRONMENT"))
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	r := chi.NewRouter()
	go startMetricsServer()

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
			DB: db,
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

	menuStore := dbstore.NewMenuStore(
		dbstore.NewMenuStoreParams{
			DB: db,
		},
	)

    postStore := dbstore.NewPostStore(
        dbstore.NewPostStoreParams{
            DB: db,
        },
    )

    commentStore := dbstore.NewCommentStore(
        dbstore.NewCommentStoreParams{
            DB: db,
        },
    )

    pollStore := dbstore.NewPollStore(
        dbstore.NewPollStoreParams{
            DB: db,
        },
    )

	calendarStore := dbstore.NewCalendarStore(
		dbstore.NewCalendarStoreParams{
			DB: db,
		},
	)

	promoStore := dbstore.NewPromoStore(
		dbstore.NewPromoStoreParams{
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
			m.MetricsMiddleware,
		)

		r.NotFound(handlers.NewNotFoundHandler().ServeHTTP)

		r.Get("/", handlers.NewHomeHandler(&handlers.HomeHandlerParams{
            PostStore: postStore,
            MenuStore: menuStore,
        }).ServeHTTP)

		r.Get("/info", handlers.NewInfoHandler().ServeHTTP)

		r.Get("/post/{postId}", handlers.NewPostHandler(handlers.PostHandlerParams{
			PostStore: postStore,
		}).Show)

		r.Route("/posts", func(r chi.Router){
			handler := handlers.NewPostHandler(handlers.PostHandlerParams{
				PostStore: postStore,
			})

			r.Get("/", handler.List)
			//r.Get("/{page}", handler.List)
		})

        r.Get("/menu/{menuId}", handlers.NewMenuHandler(handlers.MenuHandlerParams{
            MenuStore: menuStore,
        }).Day)

        r.Get("/eettafel", handlers.NewEettafelHandler().ServeHTTP)
        r.Get("/eettafel/{pageId}", handlers.NewMenuHandler(handlers.MenuHandlerParams{
            MenuStore: menuStore,
        }).Week)

		//Need to be logged in to access these routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.LoggedIn)

			r.Route("/comments/{postId}", func(r chi.Router) {
				handler := handlers.NewCommentHandler(handlers.CommentHandlerParams{
					CommentStore: commentStore,
				})
				r.Get("/", handler.List)
				r.Post("/", handler.Create)
				r.Delete("/{commentId}", handler.Delete)
			})

            r.Route("/poll", func(r chi.Router) {
				handler := handlers.NewPollHandler(handlers.PollHandlerParams{
					PollStore: pollStore,
				})
				r.Get("/", handler.Show)

				r.Route("/{pollId}", func(r chi.Router) {
					r.Post("/", handler.Vote)
					r.Delete("/", handler.RemoveVote)
				})
			})

			r.Route("/promo", func(r chi.Router) {
				handler := handlers.NewPromoHandler(handlers.PromoHandlerParams{
					PromoStore: promoStore,
				})
				r.Get("/", handler.Show)
				r.Get("/{promoId}", handler.Popup)
			})

			r.Route("/agenda", func(r chi.Router) {
				handler := handlers.NewCalendarHandler(handlers.CalendarHandlerParams{
					CalendarStore: calendarStore,
				})
				r.Get("/{day}", handler.Day)
				r.Get("/{eventId}/popup", handler.Popup)
			})

			r.Route("/webalmanak", func(r chi.Router) {
				r.Route("/leden", func(r chi.Router) {
					handler := handlers.NewUserHandler(handlers.UserHandlerParams{
						UserStore: userStore,
						GroupUserStore: groupUserStore,
					})
					r.Get("/", handler.List)
					r.Post("/", handler.Search)

					r.Get("/{userId}", handler.Show)
					r.Get("/{userId}/edit", handler.Edit)
					r.Patch("/{userId}/edit", handler.Patch)
				})
                r.Route("/jaarclubs", func(r chi.Router) {
                    r.Get("/", handlers.NewJaarclubsHandler(handlers.GetJaarclubsHandlerParams{
                        GroupStore: groupStore,
                    }).ServeHTTP)
                    r.Get("/{jaarlaag}", handlers.NewJaarclubsHandler(handlers.GetJaarclubsHandlerParams{
                        GroupStore: groupStore,
                    }).ServeHTTP)
                })
                r.Route("/commissies", func(r chi.Router) {
                    r.Get("/", handlers.NewCommissiesHandler(handlers.GetCommissiesHandlerParams{
                        GroupStore: groupStore,
                    }).ServeHTTP)
				})
				r.Route("/{groupType}", func(r chi.Router) {
					r.Get("/", handlers.NewGroupsHandler(handlers.GetGroupsHandlerParams{
						GroupStore: groupStore,
					}).ServeHTTP)
				})
				r.Get("/groep/{groupId}", handlers.NewGroupHandler(handlers.GetGroupHandlerParams{
					GroupStore:     groupStore,
					GroupUserStore: groupUserStore,
				}).ServeHTTP)
			})
		})

		r.Route("/admin", func(r chi.Router) {
			r.Use(authMiddleware.IsAdmin)
			r.Get("/", admin.NewAdminHandler().ServeHTTP)
			r.Route("/menu", func(r chi.Router) {
				handler := admin.NewMenuHandler(admin.MenuHandlerParams{
					MenuStore: menuStore,
				})
				r.Get("/", handler.New)
				r.Post("/", handler.Create)
			})

            r.Route("/post", func(r chi.Router) {
				handler := admin.NewPostHandler(admin.PostHandlerParam{
					PostStore: postStore,
				})
                r.Get("/", handler.New)
                r.Post("/", handler.Create) 
                r.Get("/{postId}", handler.Edit)
                r.Patch("/{postId}", handler.Patch)
                r.Delete("/{postId}", handler.Delete)
            })
			r.Route("/promo", func(r chi.Router) {
				handler := admin.NewPromoHandler(admin.PromoHandlerParams{
					PromoStore: promoStore,
				})
				r.Get("/", handler.New)
				r.Post("/", handler.Put) //Same logic as put
				r.Delete("/inactive", handler.DeleteInactive)

				r.Get("/{promoId}", handler.Edit)
				r.Put("/{promoId}", handler.Put)
				r.Delete("/{promoId}", handler.Delete)
			})
            r.Post("/upload", handlers.NewPostUploadHandler().ServeHTTP)
			r.Route("/leden", func(r chi.Router) {
				handler := admin.NewUserHandler(admin.UserHandlerParams{
					UserStore: userStore,
				})				
				r.Get("/", handler.List)
				r.Post("/", handler.Search)
				r.Route("/{userId}", func(r chi.Router) {
					r.Get("/", handler.Edit)
					r.Patch("/", handler.Patch)
					r.Delete("/", handler.Delete)
				})
			})
            r.Route("/poll", func(r chi.Router) {
				handler := admin.NewPollHandler(admin.PollHandlerParams{
                    PollStore: pollStore,
                })
                r.Get("/", handler.New)
                r.Post("/", handler.Create)
                r.Route("/{pollId}", func(r chi.Router) {
                    r.Get("/", handler.Edit)
                    r.Put("/", handler.Put)
                    r.Delete("/", handler.Delete)
                    r.Post("/activate", handler.Activate)
                })
            })
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

	logger.Info("Server started", slog.String("port", cfg.Port), slog.String("env", os.Getenv("env")))
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

func startMetricsServer() {

    mux := http.NewServeMux()
    mux.Handle("/metrics", promhttp.Handler())

    addr := ":9091"
    fmt.Printf("Metrics server listening on %s\n", addr)
    if err := http.ListenAndServe(addr, mux); err != nil {
        fmt.Printf("Metrics server failed: %v\n", err)
    }
}
