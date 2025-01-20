package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"sjb_site/internal/config"
	"sjb_site/internal/handlers"
	"sjb_site/internal/hash/passwordhash"
	database "sjb_site/internal/store/db"
	"sjb_site/internal/store/dbstore"
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

		r.Get("/", handlers.NewHomeHandler(&handlers.HomeHandlerParams{
			PostStore: postStore,
		}).ServeHTTP)

		r.Get("/post/{postId}", handlers.NewPostHandler(handlers.PostHandlerParams{
			PostStore: postStore,
		}).ServeHTTP)

		r.Get("/posts", handlers.NewPostsHandler(handlers.PostsHandlerParams{
			PostsStore: postStore,
		}).ServeHTTP)
		r.Get("/posts/{page}", handlers.NewPostsHandler(handlers.PostsHandlerParams{
			PostsStore: postStore,
		}).ServeHTTP)

		r.Get("/menu/{menuId}", handlers.NewMenuHandler(handlers.GetMenuHandlerParams{
			MenuStore: menuStore,
		}).ServeHTTP)

		r.Route("/comments/{postId}", func(r chi.Router) {
			r.Get("/", handlers.NewCommentsHandler(handlers.CommentsHandlerParams{
				CommentStore: commentStore,
			}).ServeHTTP)
			r.Post("/", handlers.NewPostCommentHandler(handlers.PostCommentHandlerParams{
				CommentStore: commentStore,
			}).ServeHTTP)
			r.Delete("/{commentId}", handlers.NewDeleteCommentHandler(handlers.DeleteCommentHandlerParams{
				CommentStore: commentStore,
			}).ServeHTTP)
		})

		//Need to be logged in to access these routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.LoggedIn)
			r.Route("/webalmanak", func(r chi.Router) {
				r.Route("/leden", func(r chi.Router) {
					r.Get("/", handlers.NewUserSearchHandler().ServeHTTP)
					r.Post("/", handlers.NewPostUserSearchHandler(handlers.PostUserSearchHandlerParams{
						UserStore: userStore,
					}).ServeHTTP)

					r.Get("/{userId}", handlers.NewUserHandler(handlers.GetUserHandlerParams{
						UserStore: userStore,
					}).ServeHTTP)
					r.Get("/{userId}/edit", handlers.NewUserEditHandler(handlers.GetUserEditHandlerParams{
						UserStore: userStore,
					}).ServeHTTP)
					r.Patch("/{userId}/edit", handlers.NewPatchtUserEditHandler(handlers.PatchUserEditHandlerParams{
						UserStore: userStore,
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
			r.Get("/", handlers.NewAdminHandler().ServeHTTP)
			r.Get("/menu", handlers.NewGetCreateMenuHandler().ServeHTTP)
			r.Post("/menu", handlers.NewPostCreateMenuHandler(handlers.PostCreateMenuHandlerParams{
				MenuStore: menuStore,
			}).ServeHTTP)
			r.Get("/post", handlers.NewGetCreatePostHandler().ServeHTTP)
			r.Post("/post", handlers.NewPostCreatePostHandler(handlers.PostCreatePostHandlerParams{
				PostStore: postStore,
			}).ServeHTTP)
			r.Route("/leden", func(r chi.Router) {
				r.Get("/", handlers.NewGetUserManagementHandler().ServeHTTP)
				r.Post("/", handlers.NewPostUserManagementHandler(handlers.PostUserManagementHandlerParams{
					UserStore: userStore,
				}).ServeHTTP)
				r.Route("/{userId}", func(r chi.Router) {
					r.Get("/", handlers.NewAdminUserEditHandler(handlers.GetAdminUserEditHandlerParams{
						UserStore: userStore,
					}).ServeHTTP)
					r.Patch("/", handlers.NewPatchAdminUserEditHandler(handlers.PatchAdminUserEditHandlerParams{
						UserStore: userStore,
					}).ServeHTTP)
					r.Delete("/delete", handlers.NewDeleteUserHandler(handlers.DeleteUserHandlerParams{
						UserStore: userStore,
					}).ServeHTTP)
				})
			})
			r.Get("/groep", handlers.NewGetCreateGroupHandler().ServeHTTP)
			r.Post("/groep", handlers.NewPostCreateGroupHandler(handlers.PostCreateGroupHandlerParams{
				GroupStore: groupStore,
			}).ServeHTTP)
			r.Get("/groep/{groupId}", handlers.NewAdminGroupEditHandler(handlers.GetAdminGroupEditHandlerParams{
				GroupStore: groupStore,
			}).ServeHTTP)
			r.Patch("/groep/{groupId}", handlers.NewPatchAdminGroupEditHandler(handlers.PatchAdminGroupEditHandlerParams{
				GroupStore: groupStore,
			}).ServeHTTP)
			r.Delete("/groep/{groupId}", handlers.NewDeleteGroupHandler(handlers.DeleteGroupHandlerParams{
				GroupStore: groupStore,
			}).ServeHTTP)
			r.Route("/groepen", func(r chi.Router) {
				r.Get("/{groupType}", handlers.NewGroupManagementHandler(handlers.GetGroupManagementHandlerParams{
					GroupStore: groupStore,
				}).ServeHTTP)
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
