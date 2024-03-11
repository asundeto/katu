package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GitHubResponseToken struct {
	AccessToken string `json:"access_token"`
	TokenID     string `json:"id_token"`
	Scope       string `json:"scope"`
}

type GitHubLoginUserData struct {
	ID         int
	Name       string
	Email      string
	Password   string
	FirstName  string
	SecondName string
	Login      string
	Provider   string
}

const (
	GitHubAuthURL      = "https://github.com/login/oauth/authorize"
	GitHubClientID     = "7204d1f96b4db7e5d453"
	GitHubRedirectURL  = "https://localhost:8080/auth/github/callback"
	GitHubClientSecret = "2a4621170475143853a9752bf405fb1d2f781051"
)

func (app *Application) GitHUBAuthorization(data *GitHubLoginUserData) (*Session, error) {
	session := &Session{}
	return session, nil
}

func (app *Application) GithubAuthHandler(w http.ResponseWriter, r *http.Request) {
	// prompt=consent    ---- this is needed for prompting the user to confirm authorisation
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=offline&prompt=consent", GitHubAuthURL, GitHubClientID, GitHubRedirectURL)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *Application) GithubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code") // temporary token given by Github
	if code == "" {
		// helpers.ErrorHandler(w, http.StatusUnauthorized, errors.New("Temporary token is invalid"))
		app.ServerError(w, errors.New("Temporary token is invalid"), r)
	}

	tokenRes, err := getGithubOauthToken(code)
	if err != nil {
		// helpers.ErrorHandler(w, http.StatusBadGateway, errors.New("The information received from Github"))
		app.ServerError(w, errors.New("The information received from Github"), r)
		return
	}

	githubData, err := getGithubData(tokenRes.AccessToken)
	if err != nil {
		// helpers.ErrorHandler(w, http.StatusBadGateway, errors.New("The information received from GitHub"))
		app.ServerError(w, errors.New("The information received from Github"), r)
		return
	}
	fmt.Println("&&&", githubData)

	userData, err := getUserData(githubData)
	if err != nil {
		app.ServerError(w, errors.New("The information received from Github"), r)
		return
	}
	fmt.Println("-------------------", userData.Login)
	loginCheck, id := app.loginCheckAuth(userData.Login, userData.Password)
	if loginCheck == "Email non registered yet!" {
		fmt.Println("First time Registration")
		var profilePhoto = "default.jpg"
		err = app.Users.Insert(userData.Login, userData.Login, userData.Password, profilePhoto)
		if err != nil {
			app.ServerError(w, errors.New("Server Error"), r)
			return
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
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getUserData(data string) (GitHubLoginUserData, error) {
	userData := GitHubLoginUserData{}
	if err := json.Unmarshal([]byte(data), &userData); err != nil {
		return GitHubLoginUserData{}, err
	}

	return userData, nil
}

func getGithubOauthToken(code string) (*GitHubResponseToken, error) {
	requestBodyMap := map[string]string{
		"client_id":     GitHubClientID,
		"client_secret": GitHubClientSecret,
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

	var ghresp GitHubResponseToken
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
