package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	errorhandler "yinyang/internal/errors"
	"yinyang/internal/models"
)

func (app *Application) CheckSession(w http.ResponseWriter, r *http.Request) (*models.Session, error) {
	token, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, nil
		}
		return nil, err
	}
	session, err := app.Posts.GetSessionFromToken(token.Value)
	if err != nil {
		if errors.Is(err, errorhandler.ErrNoRecord) {
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   "",
				Expires: time.Now().Add(-1 * time.Minute),
			})
			fmt.Println()
			return nil, nil
		} else {
			return nil, err
		}
	}
	return session, nil
}

func (app *Application) DeleteExpiredSessions() error {
	err := app.Posts.DeleteExpiredSessions()
	if err != nil {
		app.ErrorLog.Println(err)
	}
	return err
}
