package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode/utf8"

	"yinyang/internal/models"
	"yinyang/internal/validator"
)

type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Posts         *models.Model
	TemplateCache map[string]*template.Template
	Users         *models.UserModel

	categories []string
	Reactions  *models.ReactionModel
}

type CategoriesForm struct {
	Game        string
	Films       string
	Programming string
	Anime       string
	Sport       string
	validator.Validator
}

type CommentCreateForm struct {
	CContent string
	validator.Validator
}

type PostCreateForm struct {
	Title      string
	Content    string
	Image      string
	Category   string
	Categories []string
	validator.Validator
}

type ErrorStruct struct {
	Status int
	Text   string
}

// TO DO -----------------------------------------------
type PostMessage struct {
	Title   string
	Content string
	Error   string
}

type Msg struct {
	Form             string
	Error            string
	ReturnedEmail    string
	ReturnedPass     string
	ReturnedPass2    string
	ReturnedUsername string
}

type UserInfo struct {
	UserInfoName  string
	UserInfoEmail string
}

type PostCategoriesForm struct {
	Id             int
	CategoriesName string
}

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.NotFound(w, r)
		return
	}
	session, err := app.CheckSession(w, r)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	fmt.Println("SESSION", session)
	switch r.Method {
	case http.MethodGet:
		category := r.URL.Query().Get("category")

		var posts []*models.Post

		switch category {
		case "created":
			if session == nil {
				http.Redirect(w, r, "/user/login", http.StatusSeeOther)
				return
			}
			posts, err = app.Posts.GetPostsByUser(session.UserName)

		case "liked":
			if session == nil {
				http.Redirect(w, r, "/user/login", http.StatusSeeOther)
				return
			}
			posts, err = app.Posts.GetPostsByUserReaction(session.UserID)

		case "latest":
			posts, err = app.Posts.Latest()

		default:
			posts, err = app.Posts.Latest()
		}

		if err != nil {
			app.ServerError(w, err, r)
			return
		}

		data := app.NewTemplateData(r)
		data.Posts = posts
		data.IsAuthenticated = session != nil
		if data.IsAuthenticated {
			data.UserName = session.UserName
		}

		app.Render(w, http.StatusOK, "home.html", data, r)
	case http.MethodPost:

		var posts []*models.Post
		r.ParseForm()

		form := &CategoriesForm{
			Game:        r.FormValue("Game"),
			Films:       r.FormValue("Films"),
			Programming: r.FormValue("Programming"),
			Anime:       r.FormValue("Anime"),
			Sport:       r.FormValue("Sport"),
		}

		posts, err = app.Posts.Latest()
		if err != nil {
			app.ServerError(w, err, r)
			return
		}

		var filteredPosts []*models.Post

		for i := range posts {
			switch {
			case form.Game != "" && posts[i].Category != "" && strings.Contains(posts[i].Category, form.Game):
				filteredPosts = append(filteredPosts, posts[i])
			case form.Films != "" && posts[i].Category != "" && strings.Contains(posts[i].Category, form.Films):
				filteredPosts = append(filteredPosts, posts[i])
			case form.Programming != "" && posts[i].Category != "" && strings.Contains(posts[i].Category, form.Programming):
				filteredPosts = append(filteredPosts, posts[i])
			case form.Anime != "" && posts[i].Category != "" && strings.Contains(posts[i].Category, form.Anime):
				filteredPosts = append(filteredPosts, posts[i])
			case form.Sport != "" && posts[i].Category != "" && strings.Contains(posts[i].Category, form.Sport):
				filteredPosts = append(filteredPosts, posts[i])
			}
		}

		if filteredPosts == nil {
			data := app.NewTemplateData(r)
			data.Posts = []*models.Post{}
			data.IsAuthenticated = session != nil

			app.Render(w, http.StatusOK, "home.html", data, r)
			return
		}

		data := app.NewTemplateData(r)
		data.Posts = filteredPosts
		data.IsAuthenticated = session != nil
		if data.IsAuthenticated {
			data.UserName = session.UserName
		}

		app.Render(w, http.StatusOK, "home.html", data, r)
	}
}

func (app *Application) PostCategories(w http.ResponseWriter, r *http.Request, category string) {
	session, err := app.CheckSession(w, r)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	var posts []*models.Post
	posts, err = app.Posts.Latest()
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	var filteredPosts []*models.Post
	for i := range posts {
		if posts[i].Category != "" && strings.Contains(posts[i].Category, category) {
			filteredPosts = append(filteredPosts, posts[i])
		}
	}
	if filteredPosts == nil {
		data := app.NewTemplateData(r)
		data.Posts = []*models.Post{}
		data.IsAuthenticated = session != nil
		if data.IsAuthenticated {
			data.UserName = session.UserName
		}
		app.Render(w, http.StatusOK, "home.html", data, r)
		return
	}
	data := app.NewTemplateData(r)
	data.Posts = filteredPosts
	data.IsAuthenticated = session != nil
	if data.IsAuthenticated {
		data.UserName = session.UserName
	}
	app.Render(w, http.StatusOK, "home.html", data, r)
}

func (app *Application) PostView(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/post/view/")
	if idStr[0] == '0' {
		app.NotFound(w, r)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		app.NotFound(w, r)
		return
	}

	post, err := app.Posts.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.NotFound(w, r)
		} else {
			app.ServerError(w, err, r)
		}
		return
	}

	session, err := app.CheckSession(w, r)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}

	comments, err := app.Posts.GetComments(post.ID)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	data := app.NewTemplateData(r)
	data.Post = post
	data.Comments = comments
	data = addCategoriesToData(data)
	data.IsAuthenticated = session != nil
	if data.IsAuthenticated {
		data.UserName = session.UserName
	}

	if session != nil {
		data.Post.IsAuthenticated = true
		for i := range data.Comments {
			data.Comments[i].IsAuthenticated = true
		}
	}

	app.Render(w, http.StatusOK, "view.html", data, r)
}

func (app *Application) CreateComment(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/post/view/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		app.NotFound(w, r)
		return
	}
	session, err := app.CheckSession(w, r)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	if session == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	comment := r.FormValue("comment")

	if comment == "" || strings.TrimSpace(comment) == "" || utf8.RuneCountInString(comment) > 300 || countLines(comment) > 15 {
		post, err := app.Posts.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.NotFound(w, r)
			} else {
				app.ServerError(w, err, r)
			}
			return
		}

		comments, err := app.Posts.GetComments(post.ID)
		if err != nil {
			app.ServerError(w, err, r)
			return
		}

		data := &TemplateData{
			CurrentYear:     time.Now().Year(),
			Post:            post,
			IsAuthenticated: true,
			Comments:        comments,
			CommentError:    true,
			UserName:        session.UserName,
		}
		data = addCategoriesToData(data)

		app.Render(w, http.StatusOK, "view.html", data, r)
		return
	}

	comment = strings.Replace(comment, "\n", "<br>", -1)

	commentInput := models.Comment{
		Author:   session.UserName,
		CContent: comment,
		PostID:   id,
	}

	err = app.Posts.PostComment(commentInput)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}

func (app *Application) PostCreate(w http.ResponseWriter, r *http.Request) {
	session, err := app.CheckSession(w, r)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	data := app.NewTemplateData(r)
	data.Categories = app.categories
	data.Form = PostCreateForm{}
	data.IsAuthenticated = session != nil
	if data.IsAuthenticated {
		data.UserName = session.UserName
	}
	app.Render(w, http.StatusOK, "create.html", data, r)
}

func (app *Application) PostCreatePost(w http.ResponseWriter, r *http.Request) {
	session, err := app.CheckSession(w, r)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	if session == nil {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	err = r.ParseForm()
	if err != nil {
		app.ClientError(w, r)
		return
	}

	// image upload
	withImg := false
	file, header, err := r.FormFile("image")
	if err == nil {
		withImg = true
		defer file.Close()
		filename := filepath.Join("ui/static/uploads", header.Filename)
		out, err := os.Create(filename)
		if err != nil {
			app.ServerError(w, err, r)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			app.ServerError(w, err, r)
			return
		}
	}

	catForm := &CategoriesForm{
		Game:        r.FormValue("Game"),
		Films:       r.FormValue("Films"),
		Programming: r.FormValue("Programming"),
		Anime:       r.FormValue("Anime"),
		Sport:       r.FormValue("Sport"),
	}
	form := &PostCreateForm{}
	if withImg {
		form = &PostCreateForm{
			Title:    r.PostForm.Get("title"),
			Content:  r.PostForm.Get("content"),
			Category: getCats(catForm),
			Image:    header.Filename,
		}
	} else {
		form = &PostCreateForm{
			Title:    r.PostForm.Get("title"),
			Content:  r.PostForm.Get("content"),
			Category: getCats(catForm),
		}
	}

	userName := session.UserName
	id, err := app.Posts.Insert(form.Title, form.Content, form.Category, userName, form.Image)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}

func (app *Application) UserLogin(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = Msg{}
	app.Render(w, http.StatusOK, "login.html", data, r)
}

func (app *Application) Profile(w http.ResponseWriter, r *http.Request) {
	session, err := app.CheckSession(w, r)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	userName := session.UserName
	email, profile_photo, err := app.Users.GetEmailAndPhotoByUserName(userName)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	data := app.NewTemplateData(r)
	data.Form = models.User{
		Name:         userName,
		Email:        email,
		ProfilePhoto: profile_photo,
	}
	data.IsAuthenticated = session != nil
	if data.IsAuthenticated {
		data.UserName = session.UserName
	}
	app.Render(w, http.StatusOK, "profile.html", data, r)
}

func (app *Application) ProfileChange(w http.ResponseWriter, r *http.Request) {
	session, err := app.CheckSession(w, r)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	userName := session.UserName
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		filename := filepath.Join("ui/static/profile_photo", header.Filename)
		out, err := os.Create(filename)
		if err != nil {
			app.ServerError(w, err, r)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			app.ServerError(w, err, r)
			return
		}
	}
	err = app.Users.ChangeUserProfilePhoto(userName, header.Filename)
	if err != nil {
		app.ServerError(w, err, r)
	}
	email, profile_photo, err := app.Users.GetEmailAndPhotoByUserName(userName)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	data := app.NewTemplateData(r)
	data.Form = models.User{
		Name:         userName,
		Email:        email,
		ProfilePhoto: profile_photo,
	}
	data.IsAuthenticated = session != nil
	if data.IsAuthenticated {
		data.UserName = session.UserName
	}
	app.Render(w, http.StatusOK, "profile.html", data, r)
}

// func (app *Application) createErrorMessage()

func (app *Application) UserLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, r)
		return
	}
	email, password := strings.ToLower(r.FormValue("email")), r.FormValue("password")
	loginCheck, id := app.loginCheckAuth(email, password)
	if loginCheck != "" {
		data := app.NewTemplateData(r)
		data.Form = Msg{
			Form:          "Log",
			Error:         loginCheck,
			ReturnedEmail: email,
			ReturnedPass:  password,
		}
		app.Render(w, http.StatusUnprocessableEntity, "login.html", data, r)
		return
	}
	userName, err := app.Users.GetUserNameByEmail(email)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	token, expiration, err := app.Posts.CreateSession(id, userName)
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

func (app *Application) UserRegisterPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, r)
		return
	}
	data := app.NewTemplateData(r)
	email, username, password, passwordRepeat, checkbox := strings.ToLower(r.FormValue("email")), r.FormValue("username"), r.FormValue("password"), r.FormValue("password-repeat"), r.FormValue("checkbox")
	checkResult := app.checkRegisterData(email, username, password, passwordRepeat, checkbox)
	if checkResult != "" {
		data.Form = Msg{
			Form:             "Reg",
			ReturnedEmail:    email,
			ReturnedPass:     password,
			ReturnedPass2:    passwordRepeat,
			ReturnedUsername: username,
			Error:            checkResult,
		}
		app.Render(w, http.StatusOK, "login.html", data, r)
		return
	}
	data.Form = Msg{
		Form: "Log",
	}
	app.Render(w, http.StatusOK, "login.html", data, r)
}

func (app *Application) UserLogout(w http.ResponseWriter, r *http.Request) {
	session, err := app.CheckSession(w, r)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}

	if session != nil {
		err = app.Posts.DeleteSessionByUserId(session.UserID)
		if err != nil {
			app.ServerError(w, err, r)
			return
		}
	}
	cookie := &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now().AddDate(-1, 0, 0),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false
	}
	token := cookie.Value

	session, err := app.Posts.GetSessionFromToken(token)
	if err != nil {
		return false
	}

	if session.ExpirationDate.Before(time.Now()) {
		return false
	}
	return true
}

func (app *Application) ErrorHandler(w http.ResponseWriter, errorNum int, r *http.Request) {
	data := app.NewTemplateData(r)
	Res := &ErrorStruct{
		Status: errorNum,
		Text:   http.StatusText(errorNum),
	}
	cookie, _ := r.Cookie("session_token")
	if cookie != nil {
		token := cookie.Value
		session, errCookie := app.Posts.GetSessionFromToken(token)
		if errCookie != nil {
			app.ServerError(w, errCookie, r)
			return
		}
		data.UserName = session.UserName
	}
	data.ErrorStruct = Res
	err := app.renderErr(w, http.StatusUnprocessableEntity, "error.html", data, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *Application) checkRegisterData(email, username, password, passwordRepeat, checkbox string) string {
	if checkbox != "on" {
		return "Accept user agreement!"
	}
	userNameChecked := userNameCheck(username)
	if userNameChecked != "" {
		return userNameChecked
	}
	checkUserDublicate := app.Users.Duplicates(email, username)
	if checkUserDublicate != "" {
		return checkUserDublicate
	}
	if !IsEmailValid(email) {
		return "Please enter correct email!"
	}
	if password != passwordRepeat {
		return "Password mismatch!"
	}
	if !isPasswordValid(password) {
		return "Low password complexity Use one of them [$#%!?.*]"
	}
	var profilePhoto = "default.jpg"
	err := app.Users.Insert(username, email, password, profilePhoto)
	if err != nil {
		return "Server error! Please try later"
	}
	return ""
}

func (app *Application) loginCheckAuth(email, password string) (string, int) {
	emailErr := app.Users.EmailExist(email)
	if emailErr != nil {
		if errors.Is(emailErr, models.ErrInvalidCredentials) {
			return "Email non registered yet!", 0
		} else {
			return "Server error!", 0
		}
	}
	id, err := app.Users.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			return "Password is incorrect!", 0
		} else {
			return "Server error!", 0
		}
	}
	return "", id
}
