package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type GoogleLoginUserData struct {
	ID         int
	Name       string
	Email      string
	Password   string
	FirstName  string
	SecondName string
	Provider   string
}

type GoogleUserResult struct {
	Id             string
	Email          string
	Verified_email bool
	Name           string
	Given_name     string
	Family_name    string
	Picture        string
	Locale         string
	Password       string
}

type GoogleResponseToken struct {
	AccessToken string `json:"access_token"`
	TokenID     string `json:"id_token"`
}

type Session struct {
	UserID         int
	UserName       string
	Token          string
	ExpirationDate time.Time
}

type UserService struct {
	// Any dependencies or configurations needed by the service
}

func (app *Application) GoogleAuthorization(data *GoogleLoginUserData) (*Session, error) {
	session := &Session{}
	return session, nil
}

const (
	GoogleAuthURL      = "https://accounts.google.com/o/oauth2/auth"                                // const URL
	GoogleClientID     = "722031461724-dnvp1cl4hngcs1kgt0a2qi9j86a3dr1n.apps.googleusercontent.com" // my google account
	GoogleRedirectURL  = "https://localhost:8080/auth/google/callback"                              // callback endpoint
	GoogleClientSecret = "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"                                      // my google account
)

func (app *Application) GoogleAuthHandler(w http.ResponseWriter, r *http.Request) {
	scope := url.QueryEscape("email profile https://www.googleapis.com/auth/drive.file")
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&prompt=select_account", GoogleAuthURL, GoogleClientID, GoogleRedirectURL, scope)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *Application) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code") // temporary token given by Google

	// fmt.Println("Code: ", code) // 4/0AfJohXnLtvZf6XVZjrQRMyaH_CdPg6yB7XoUnrdBqv0wd6RkJDo361ff3yc6qzaMdS6oKQ

	if code == "" {
		app.ServerError(w, errors.New("Temporary token is invalid"), r)
		return
	}

	tokenRes, err := getGoogleOauthToken(code)
	if err != nil {
		app.ServerError(w, errors.New("The information received from Google"), r)
		return
	}

	googleUser, err := getGoogleUser(tokenRes.AccessToken, tokenRes.TokenID)
	if err != nil {
		app.ServerError(w, errors.New("The information received from Google"), r)
		return
	}

	// creating the struct type of User after Google Auth
	googleData := GoogleLoginUserData{
		Name:       googleUser.Name,
		Email:      googleUser.Email,
		FirstName:  googleUser.Given_name,
		SecondName: googleUser.Family_name,
		Provider:   "Google",
		Password:   googleUser.Password,
	}
	fmt.Println("&&&", googleData)

	loginCheck, id := app.loginCheckAuth(googleUser.Email, googleUser.Password)
	if loginCheck == "Email non registered yet!" {
		var profilePhoto = "default.jpg"
		err = app.Users.Insert(googleUser.Name, googleUser.Email, googleUser.Password, profilePhoto)
		if err != nil {
			app.ServerError(w, errors.New("Server Error"), r)
			return
		}
	}
	_, _, err = app.Posts.CreateSession(id, googleUser.Name)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	token, expiration, err := app.Posts.CreateSession(id, googleUser.Name)
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

func getGoogleOauthToken(code string) (*GoogleResponseToken, error) {
	// The URL to which the POST request is sent. This is the token endpoint of the OAuth 2.0 provider (in this case, Google).
	const rootURL = "https://oauth2.googleapis.com/token" // URL for getting the access token using code

	// map is convenient for encoding form values in the request body
	values := url.Values{}
	values.Set("code", code)
	values.Set("client_id", GoogleClientID)
	values.Set("client_secret", GoogleClientSecret)
	values.Set("redirect_uri", GoogleRedirectURL)
	values.Set("grant_type", "authorization_code")

	// Make a POST request to the Google token endpoint
	response, err := http.Post(rootURL, "application/x-www-form-urlencoded", strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var tokenRes GoogleResponseToken
	err = json.NewDecoder(response.Body).Decode(&tokenRes)
	if err != nil {
		return nil, err
	}

	return &tokenRes, nil
}

func getGoogleUser(AccessToken string, TokenID string) (*GoogleUserResult, error) {
	rootURL := fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=%s", AccessToken)

	// Prepare the request to the Google People API
	req, err := http.NewRequest("GET", rootURL, nil)
	if err != nil {
		return nil, err
	}

	// Set the Authorization header with the access token
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", TokenID))

	// Make the request to the Google People API
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var userRes GoogleUserResult
	err = json.NewDecoder(res.Body).Decode(&userRes)
	if err != nil {
		return nil, err
	}

	return &userRes, nil
}

// import (
// 	"fmt"
// 	"net/http"
// 	"net/url"
// )

// const (
// 	GoogleAuthURL      = "https://accounts.google.com/o/oauth2/auth"                                // const URL
// 	GoogleClientID     = "722031461724-dnvp1cl4hngcs1kgt0a2qi9j86a3dr1n.apps.googleusercontent.com" // my google account
// 	GoogleRedirectURL  = "https://localhost:7000/auth/google/callback"                              // callback endpoint
// 	GoogleClientSecret = "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"                                      // my google account
// )

// func (app *Application) GoogleAuthHandler(w http.ResponseWriter, r *http.Request) {
// 	scope := url.QueryEscape("email profile https://www.googleapis.com/auth/drive.file")
// 	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&prompt=select_account", GoogleAuthURL, GoogleClientID, GoogleRedirectURL, scope)
// 	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
// }

// // var (
// // 	oauthConfig = &oauth2.Config{
// // 		ClientID:     GoogleClientID,
// // 		ClientSecret: GoogleClientSecret,
// // 		RedirectURL:  GoogleRedirectURL,
// // 		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
// // 		Endpoint: oauth2.Endpoint{
// // 			AuthURL:  GoogleAuthURL,
// // 			TokenURL: "https://oauth2.googleapis.com/token",
// // 		},
// // 	}
// // 	oauthStateString = "random"
// // )

// func (app *Application) handleMain(w http.ResponseWriter, r *http.Request) {
// 	var htmlIndex = `<html><body><a href="/auth/google/login">Google Log In</a></body></html>`
// 	fmt.Fprintf(w, htmlIndex)
// }

// // func (app *Application) handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
// // 	url := oauthConfig.AuthCodeURL(oauthStateString)
// // 	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
// // }

// func (app *Application) handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
// 	state := r.FormValue("state")
// 	if state != "random" {
// 		fmt.Println("Invalid oauth state")
// 		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
// 		return
// 	}

// 	// code := r.FormValue("code")
// 	// token, err := oauthConfig.Exchange(r.Context(), code)
// 	// if err != nil {
// 	// 	fmt.Println("Code exchange failed with error:", err)
// 	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
// 	// 	return
// 	// }

// 	// You can use the token to get user information from Google
// 	// For example, you can make a request to:
// 	// https://www.googleapis.com/oauth2/v2/userinfo?access_token=TOKEN.AccessToken
// 	// fmt.Println("Access token:", token.AccessToken)
// 	// fmt.Println("Refresh token:", token.RefreshToken)
// 	// fmt.Println("Token type:", token.TokenType)
// 	// fmt.Println("Expiry:", token.Expiry)

// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
// }
