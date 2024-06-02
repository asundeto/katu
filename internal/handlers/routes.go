package handlers

import (
	"net/http"
	"strings"
	"time"
)

func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			remainingPath := r.URL.Path[len("/static/"):]
			r.URL.Path = "/" + remainingPath
			fileServer.ServeHTTP(w, r)
		} else {
			MethodNotAllowedHandler(w, r, []string{http.MethodGet})
		}
	})
	mux.Handle("/", NewRateLimiter(50, time.Minute).LimitMiddleware(app.Home))
	mux.Handle("/auth/google/login", NewRateLimiter(10, time.Minute).LimitMiddleware(app.GoogleAuthHandler))
	mux.Handle("/auth/google/callback", NewRateLimiter(10, time.Minute).LimitMiddleware(app.GoogleCallback))
	mux.Handle("/auth/github/login", NewRateLimiter(10, time.Minute).LimitMiddleware(app.GithubAuthHandler))
	mux.Handle("/auth/github/callback", NewRateLimiter(10, time.Minute).LimitMiddleware(app.GithubCallback))
	mux.Handle("/post/category/", NewRateLimiter(30, time.Minute).LimitMiddleware(app.PostCategories))
	mux.Handle("/user/login", NewRateLimiter(10, time.Minute).LimitMiddleware(app.UserLogin))
	mux.Handle("/user/login/post", NewRateLimiter(10, time.Minute).LimitMiddleware(app.UserLoginPost))
	mux.Handle("/user/register", NewRateLimiter(10, time.Minute).LimitMiddleware(app.UserRegisterPost))
	mux.Handle("/chat/", app.RequireAuthentication(app.Chat))
	mux.Handle("/user/profile/", NewRateLimiter(10, time.Minute).LimitMiddleware(app.RequireAuthentication(app.ProfileRoute)))
	mux.Handle("/messages", NewRateLimiter(1000, time.Minute).LimitMiddleware(app.RequireAuthentication(app.Messages)))
	mux.Handle("/activity", NewRateLimiter(50, time.Minute).LimitMiddleware(app.RequireAuthentication(app.Activity)))
	mux.Handle("/remove/activity/", NewRateLimiter(50, time.Minute).LimitMiddleware(app.RequireAuthentication(app.RemoveActivity)))
	mux.Handle("/post/create/post", NewRateLimiter(10, time.Minute).LimitMiddleware(app.RequireAuthentication(app.PostCreatePost)))
	mux.Handle("/likePost", NewRateLimiter(50, time.Minute).LimitMiddleware(app.RequireAuthentication(app.LikePostHandler)))
	mux.Handle("/dislikePost", NewRateLimiter(50, time.Minute).LimitMiddleware(app.RequireAuthentication(app.DislikePostHandler)))
	mux.Handle("/likeComment", NewRateLimiter(50, time.Minute).LimitMiddleware(app.RequireAuthentication(app.LikeCommentHandler)))
	mux.Handle("/dislikeComment", NewRateLimiter(50, time.Minute).LimitMiddleware(app.RequireAuthentication(app.DislikeCommentHandler)))
	mux.Handle("/user/logout/", NewRateLimiter(10, time.Minute).LimitMiddleware(app.RequireAuthentication(app.UserLogout)))
	mux.Handle("/post/create", NewRateLimiter(10, time.Minute).LimitMiddleware(app.RequireAuthentication(func(w http.ResponseWriter, r *http.Request) {
		app.PostCreate(w, r, nil)
	})))
	mux.Handle("/post/view/", NewRateLimiter(30, time.Minute).LimitMiddleware(app.PostViewRoute))
	return mux
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request, allowedMethods []string) {
	w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
