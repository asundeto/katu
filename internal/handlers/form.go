package handlers

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	temp "text/template"
	errorhandler "yinyang/internal/errors"
	"yinyang/internal/models"
	"yinyang/internal/template"
	"yinyang/internal/validator"
)

type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Posts         *models.Model
	TemplateCache map[string]*temp.Template
	Users         *models.UserModel
	Chats         *models.ChatModel
	Constructor   *validator.Constructor

	categories []string
	Reactions  *models.ReactionModel
}

type CommentCreateForm struct {
	CContent string
	validator.Validator
}

type AuthError struct {
	Error string
}

type UserInfo struct {
	UserInfoName  string
	UserInfoEmail string
}

type MessagesStruct struct {
	StartedChats []models.StartedChat
	ChatWith     ChatWith
	Users        []*models.User
}

type ChatWith struct {
	With       string
	WithPhoto  string
	WithStatus bool
	History    []models.MessageStruct
}

func (app *Application) ProfileForm(name, email, profilePhoto string, err error) models.User {
	form := models.User{
		Name:         name,
		Email:        email,
		ProfilePhoto: profilePhoto,
		Error:        err,
	}
	return form
}

func (app *Application) NotifCountSet(data *template.TemplateData) (*template.TemplateData, error) {
	actionCnt, err := app.Posts.GetUnseenActivityCount(data.UserName)
	if err != nil {
		return data, err
	}
	messagesCnt := app.Users.GetAllUnseenMessagesCount(data.UserName)
	if actionCnt != 0 || messagesCnt != 0 {
		not := models.Notifications{
			ActionCount:  actionCnt,
			MessageCount: messagesCnt,
		}
		data.Notifications = &not
	}
	return data, nil
}

func (app *Application) PostViewRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.PostView(w, r, nil)
	} else if r.Method == http.MethodPost {
		app.RequireAuthentication(app.CreateComment).ServeHTTP(w, r)
	}
}

func (app *Application) ProfileRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.Profile(w, r, nil)
	} else if r.Method == http.MethodPost {
		app.ProfileChange(w, r)
	}
}

func (app *Application) Chat(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.OpenChat(w, r)
	} else if r.Method == http.MethodPost {
		app.SendMessage(w, r)
	}
}

func (app *Application) CategoriesForm(r *http.Request) *models.CategoriesForm {
	r.ParseForm()
	form := &models.CategoriesForm{
		Game:        r.FormValue("Game"),
		Films:       r.FormValue("Films"),
		Programming: r.FormValue("Programming"),
		Anime:       r.FormValue("Anime"),
		Sport:       r.FormValue("Sport"),
	}
	return form
}

func ProfileForm(username, dataUsername, email, profile_photo string, msg error) *models.User {
	access := false
	if username == dataUsername {
		access = true
	}
	form := models.User{
		Name:         username,
		Email:        email,
		ProfilePhoto: profile_photo,
		Error:        msg,
		Access:       access,
	}
	return &form
}

func (app *Application) UserLoginForm(r *http.Request, form, email, password, passwordRepeat, username string, err error) *template.TemplateData {
	data := app.NewTemplateData(r)
	data.Form = AuthError{
		Error: err.Error(),
	}
	return data
}

func (app *Application) PostForm(r *http.Request, post *models.Post, comments []models.Comment, session *models.Session, msg error) *template.TemplateData {
	data := app.NewTemplateData(r)
	data.Post = post
	data.Comments = comments
	data.Form = msg
	data = validator.AddCategoriesToData(data)
	data = app.SetUserName(data, session)
	if session != nil {
		data.Post.IsAuthenticated = true
		for i := range data.Comments {
			data.Comments[i].IsAuthenticated = true
		}
	}
	return data
}

func (app *Application) FilteredPostsForm(posts []*models.Post, form *models.CategoriesForm) []*models.Post {
	var filteredPosts []*models.Post
	for _, post := range posts {
		if validator.MatchesCategory(post.Category, form.Game) ||
			validator.MatchesCategory(post.Category, form.Films) ||
			validator.MatchesCategory(post.Category, form.Programming) ||
			validator.MatchesCategory(post.Category, form.Anime) ||
			validator.MatchesCategory(post.Category, form.Sport) {
			filteredPosts = append(filteredPosts, post)
		}
	}
	return filteredPosts
}

func (app *Application) FilterdPostForm(w http.ResponseWriter, r *http.Request, data *template.TemplateData, category string) {
	var filteredPosts []*models.Post
	posts, err := app.Posts.Latest()
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	for i := range posts {
		if posts[i].Category != "" && strings.Contains(posts[i].Category, category) {
			filteredPosts = append(filteredPosts, posts[i])
		}
	}
	if filteredPosts == nil {
		app.Render(w, http.StatusOK, "home.html", data, r)
		return
	}
	data.Posts = filteredPosts
	app.Render(w, http.StatusOK, "home.html", data, r)
}

func (app *Application) SetUserName(data *template.TemplateData, session *models.Session) *template.TemplateData {
	data.IsAuthenticated = session != nil
	if data.IsAuthenticated {
		data.UserName = session.UserName
	}
	return data
}

func (app *Application) HomeCase(w http.ResponseWriter, r *http.Request, data *template.TemplateData, session *models.Session) {
	var posts []*models.Post
	var err error
	data.Posts = posts

	switch r.Method {
	case http.MethodGet:
		// Get the category from the request
		category := r.FormValue("category")

		// Check if the user is not logged in for certain categories
		if (category == "created" || category == "liked") && session == nil {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Determine which posts to fetch based on the category
		switch category {
		case "created":
			posts, err = app.Posts.GetPostsByUser(session.UserName)
		case "liked":
			posts, err = app.Posts.GetPostsByUserReaction(session.UserID)
		default:
			posts, err = app.Posts.Latest()
		}

	case http.MethodPost:
		// Get form data
		form := app.CategoriesForm(r)

		// Get latest posts
		posts, err = app.Posts.Latest()

		// Filter posts based on form categories
		if err == nil && form != nil {
			posts = app.FilteredPostsForm(posts, form)
		}
	}

	if err != nil {
		app.ServerError(w, err, r)
		return
	}

	data.Posts = posts
	app.Render(w, http.StatusOK, "home.html", data, r)
}

func (app *Application) PostCreateFormFunc(r *http.Request) *models.PostCreateForm {
	file, header, err := r.FormFile("image")
	var errImg error
	var image string
	catForm := app.CategoriesForm(r)
	form := &models.PostCreateForm{
		Title:    r.FormValue("title"),
		Content:  r.FormValue("content"),
		Category: validator.GetCats(catForm),
	}
	form = app.Constructor.PostCreateCheck(form)
	if err == nil && form.Error.TitleError == "" && form.Error.ContentError == "" {
		image, errImg = app.uploadImage("uploads", file, header)
	}
	if errImg != nil {
		form.Error.ImageError = errImg.Error()
	}
	if image != "" {
		form.Image = image
	}
	return form
}

func (app *Application) uploadImage(path string, file multipart.File, header *multipart.FileHeader) (string, error) {
	// Check image size
	if header.Size > 20000000 {
		return "", errorhandler.ErrPostImageSize
	}

	// Generate random image name
	image, err := validator.RandomStr(header.Filename)
	if err != nil {
		return "", err
	}

	err = validator.PathExists(path)
	if err != nil {
		return "", err
	}

	// Create full path for saving image
	filename := filepath.Join("ui/static", path, image)
	out, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Copy file contents to the created file
	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return image, nil
}

func (app *Application) RemoveOldPhoto(username string) error {
	_, profile_photo, err := app.Users.GetEmailAndPhotoByUserName(username)
	if err != nil {
		return err
	}
	if profile_photo != "" && profile_photo != "default.jpg" {
		if _, err = os.Stat("ui/static/profile_photo/" + profile_photo); err == nil {
			err = os.Remove("ui/static/profile_photo/" + profile_photo)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (app *Application) LoginCheckAuth(email, password string) (int, error) {
	if email == "" {
		return 0, errorhandler.ErrEmailExist
	}
	emailErr := app.Users.EmailExist(email)
	if emailErr != nil {
		return 0, emailErr
	}
	id, err := app.Users.Authenticate(email, password)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (app *Application) CheckRegisterData(email, username, password, passwordRepeat, checkbox string) error {
	if checkbox != "on" {
		return errorhandler.ErrUserAgrement
	}

	if err := validator.UserNameCheck(username); err != nil {
		return err
	}

	if err := app.Users.Duplicates(email, username); err != nil {
		return err
	}

	if !validator.IsEmailValid(email) {
		return errorhandler.ErrEnterCorrectEmail
	}

	if password != passwordRepeat {
		return errorhandler.ErrPasswordMismatch
	}

	if !validator.IsPasswordValid(password) {
		return errorhandler.ErrLowPassword
	}

	if err := app.Users.Insert(username, email, password, "default.jpg"); err != nil {
		return errorhandler.ErrAuthServer
	}

	return nil
}

func (app *Application) RemoveActionUndo(activity *models.Activity) error {
	if activity.Type == "createpost" {
		activitiesAll, err := app.Posts.AllActivitiesGet()
		if err != nil {
			return err
		}
		for i := 0; i < len(activitiesAll); i++ {
			if activitiesAll[i].Post.ID == activity.Post.ID {
				app.RemoveAction(&activitiesAll[i])
			}
		}
		err = app.Posts.DeletePost(activity.Post.ID)
		if err != nil {
			return err
		}
	} else {
		err := app.RemoveAction(activity)
		if err != nil {
			return err
		}
		return nil
	}
	err := app.Posts.DeleteActivity(activity.ID)
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) RemoveAction(activity *models.Activity) error {
	if activity.Type == "createcomment" {
		err := app.Posts.RemoveComment(activity.Comment.Id)
		if err != nil {
			return err
		}
	} else if activity.Type == "reactionpost" {
		userId, err := app.Users.GetIDByUserName(activity.Username)
		if err != nil {
			return err
		}
		err = app.Reactions.RemovePostReaction(userId, activity.Post.ID)
		if err != nil {
			return err
		}
	} else if activity.Type == "reactioncomment" {
		userId, err := app.Users.GetIDByUserName(activity.Username)
		if err != nil {
			return err
		}
		err = app.Reactions.RemoveCommentReaction(userId, activity.Comment.Id)
		if err != nil {
			return err
		}
	}
	err := app.Posts.DeleteActivity(activity.ID)
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) ReactionPostHandler(w http.ResponseWriter, r *http.Request, boly bool) {
	session := app.SessionCheck(w, r)
	userID, err := app.Users.GetIDByUserName(session.UserName)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.NotFound(w, r)
		return
	}
	if boly {
		boly, _ = app.Reactions.LikePost(userID, id)
	} else {
		boly, _ = app.Reactions.DislikePost(userID, id)
	}
	if boly {
		err = app.Posts.ActivityInsert(session.UserName, "reactionpost", id, 0, 0)
		if err != nil {
			app.ServerError(w, err, r)
			return
		}
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func (app *Application) ReactionCommentHandler(w http.ResponseWriter, r *http.Request, boly bool) {
	session := app.SessionCheck(w, r)
	userID, err := app.Users.GetIDByUserName(session.UserName)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.NotFound(w, r)
		return
	}
	postID, err := app.Reactions.GetCommentPostID(id)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	if boly {
		boly, _ = app.Reactions.LikeComment(userID, id, postID)
	} else {
		boly, _ = app.Reactions.DislikeComment(userID, id, postID)
	}
	if boly {
		err = app.Posts.ActivityInsert(session.UserName, "reactioncomment", postID, id, 0)
		if err != nil {
			app.ServerError(w, err, r)
			return
		}
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}
