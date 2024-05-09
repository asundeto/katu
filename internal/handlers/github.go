package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	errorhandler "yinyang/internal/errors"
	"yinyang/internal/models"
)

func (app *Application) GitHUBAuthorization(data *models.GitHubLoginUserData) (*Session, error) {
	session := &Session{}
	return session, nil
}

func (app *Application) GithubAuthHandler(w http.ResponseWriter, r *http.Request) {
	// prompt=consent    ---- this is needed for prompting the user to confirm authorisation
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=offline&prompt=consent", models.GitHubAuthURL, models.GitHubClientID, models.GitHubRedirectURL)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *Application) GithubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code") // temporary token given by Github
	if code == "" {
		app.ServerError(w, errorhandler.ErrZeroCode, r)
	}

	tokenRes, err := getGithubOauthToken(code)
	if err != nil {
		app.ServerError(w, errorhandler.ErrGitHUBInfo, r)
		return
	}

	githubData, err := getGithubData(tokenRes.AccessToken)
	if err != nil {
		app.ServerError(w, errorhandler.ErrGitHUBInfo, r)
		return
	}
	fmt.Println("&&&", githubData)

	userData, err := getUserData(githubData)
	if err != nil {
		app.ServerError(w, errorhandler.ErrGitHUBInfo, r)
		return
	}
	id, err := app.LoginCheckAuth(userData.Email, userData.Password)
	if err != nil {
		if errors.Is(err, errorhandler.ErrInvalidCredentials) {
			var profilePhoto = "default.jpg"
			err = app.Users.Insert(userData.Name, userData.Email, userData.Password, profilePhoto)
			if err != nil {
				app.ServerError(w, errorhandler.ErrServerError, r)
				return
			}
		} else {
			app.ServerError(w, err, r)
		}
	}
	_, _, err = app.Posts.CreateSession(id, userData.Login)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	token, expiration, err := app.Posts.CreateSession(id, userData.Login)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	cookie := &http.Cookie{
		Name:    "session_token",
		Value:   token,
		Expires: expiration,
		Path:    "/",
	}
	app.Users.UserStatusOnline(userData.Name)
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getUserData(data string) (models.GitHubLoginUserData, error) {
	userData := models.GitHubLoginUserData{}
	if err := json.Unmarshal([]byte(data), &userData); err != nil {
		return models.GitHubLoginUserData{}, err
	}

	return userData, nil
}

func getGithubOauthToken(code string) (*models.GitHubResponseToken, error) {
	requestBodyMap := map[string]string{
		"client_id":     models.GitHubClientID,
		"client_secret": models.GitHubClientSecret,
		"code":          code,
	}
	requestJSON, err := json.Marshal(requestBodyMap)
	if err != nil {
		return nil, err
	}

	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		return nil, reqerr
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return nil, resperr
	}

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ghresp models.GitHubResponseToken
	if err := json.Unmarshal(respbody, &ghresp); err != nil {
		return nil, err
	}

	return &ghresp, nil
}

func getGithubData(accessToken string) (string, error) {
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if reqerr != nil {
		return "", reqerr
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return "", resperr
	}

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respbody), nil
}
