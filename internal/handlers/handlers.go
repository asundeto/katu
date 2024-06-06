package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	errorhandler "yinyang/internal/errors"
	"yinyang/internal/models"
	"yinyang/internal/validator"
)

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.NotFound(w, r)
		return
	}
	session := app.SessionCheck(w, r)
	data := app.NewTemplateData(r)
	data = app.SetUserName(data, session)
	data, err := app.NotifCountSet(data)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	app.HomeCase(w, r, data, session)
}

func (app *Application) PostCategories(w http.ResponseWriter, r *http.Request) {
	category := validator.ToUpperFirst(strings.TrimPrefix(r.URL.Path, "/post/category/"))
	session := app.SessionCheck(w, r)
	data := app.NewTemplateData(r)
	data = app.SetUserName(data, session)
	app.FilterdPostForm(w, r, data, category)
}

func (app *Application) PostCreate(w http.ResponseWriter, r *http.Request, form *models.PostCreateForm) {
	session := app.SessionCheck(w, r)
	data := app.NewTemplateData(r)
	data.Categories = app.categories
	data.Form = form
	data = app.SetUserName(data, session)
	data, err := app.NotifCountSet(data)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	app.Render(w, http.StatusOK, "create.html", data, r)
}

func (app *Application) PostCreatePost(w http.ResponseWriter, r *http.Request) {
	session := app.SessionCheck(w, r)
	form := app.PostCreateFormFunc(r)
	if form.Error.TitleError != "" || form.Error.ContentError != "" || form.Error.ImageError != "" {
		jsonResponse, err := json.Marshal(form.Error)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Error(w, string(jsonResponse), http.StatusBadRequest)
		return
		// app.PostCreate(w, r, form)
	}
	id, err := app.Posts.Insert(form.Title, form.Content, form.Category, session.UserName, form.Image)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	err = app.Posts.ActivityInsert(session.UserName, "createpost", id, 0, 0)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}

func (app *Application) PostView(w http.ResponseWriter, r *http.Request, msg error) {
	id := validator.ValidPostId(r)
	if id == 0 {
		app.NotFound(w, r)
		return
	}
	post, err := app.Posts.Get(id)
	if err != nil {
		if errors.Is(err, errorhandler.ErrNoRecord) {
			app.NotFound(w, r)
			return
		} else {
			app.ServerError(w, err, r)
			return
		}
	}
	session := app.SessionCheck(w, r)
	comments, err := app.Posts.GetComments(post.ID)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	data := app.PostForm(r, post, comments, session, msg)
	data, err = app.NotifCountSet(data)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	app.Render(w, http.StatusOK, "view.html", data, r)
}

func (app *Application) CreateComment(w http.ResponseWriter, r *http.Request) {
	id := validator.ValidPostId(r)
	if id == 0 {
		app.NotFound(w, r)
		return
	}
	session := app.SessionCheck(w, r)
	comment := r.FormValue("comment")
	if comment == "" || strings.TrimSpace(comment) == "" || utf8.RuneCountInString(comment) > 300 || validator.CountLines(comment) > 15 {
		app.PostView(w, r, errorhandler.ErrComment)
		return
	}
	commentInput := validator.CommentForm(comment, session.UserName, id)
	commentId, err := app.Posts.PostComment(commentInput)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	err = app.Posts.ActivityInsert(session.UserName, "createcomment", id, int(commentId), 0)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	app.PostView(w, r, nil)
}

func (app *Application) Profile(w http.ResponseWriter, r *http.Request, msg error) {
	username := validator.ValidUserName(r)
	data := app.NewTemplateData(r)
	session := app.SessionCheck(w, r)
	data = app.SetUserName(data, session)
	email, profile_photo, err := app.Users.GetEmailAndPhotoByUserName(username)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	data.Form = ProfileForm(username, data.UserName, email, profile_photo, msg)
	data, err = app.NotifCountSet(data)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	app.Render(w, http.StatusOK, "profile.html", data, r)
}

func (app *Application) ProfileChange(w http.ResponseWriter, r *http.Request) {
	session := app.SessionCheck(w, r)
	file, header, err := r.FormFile("image")
	if err != nil {
		app.Profile(w, r, errorhandler.ErrUploadImage)
		return
	}
	image, err := app.uploadImage("profile_photo", file, header)
	if err != nil {
		app.Profile(w, r, err)
		return
	}
	err = app.RemoveOldPhoto(session.UserName)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	err = app.Users.ChangeUserProfilePhoto(session.UserName, image)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	app.Profile(w, r, nil)
}

func (app *Application) SessionCheck(w http.ResponseWriter, r *http.Request) *models.Session {
	session, err := app.CheckSession(w, r)
	if err != nil {
		app.ServerError(w, err, r)
		return nil
	}
	return session
}

func (app *Application) UserLogin(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = AuthError{}
	app.Render(w, http.StatusOK, "login.html", data, r)
}

func (app *Application) UserLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, r)
		return
	}
	email, password := strings.ToLower(r.FormValue("email")), r.FormValue("password")
	id, err := app.LoginCheckAuth(email, password)
	if err != nil {
		loginError := AuthError{
			Error: err.Error(),
		}
		jsonResponse, err := json.Marshal(loginError)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Error(w, string(jsonResponse), http.StatusBadRequest)
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
	app.Users.UserStatusOnline(userName)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) UserRegisterPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, r)
		return
	}
	email, username, password, passwordRepeat, checkbox := strings.ToLower(r.FormValue("email")), r.FormValue("username"), r.FormValue("password"), r.FormValue("password-repeat"), r.FormValue("checkbox")
	checkResult := app.CheckRegisterData(email, username, password, passwordRepeat, checkbox)
	if checkResult != nil {
		loginError := AuthError{
			Error: checkResult.Error(),
		}
		jsonResponse, err := json.Marshal(loginError)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Error(w, string(jsonResponse), http.StatusBadRequest)
		return
	}
	app.UserLogin(w, r)
}

func (app *Application) UserLogout(w http.ResponseWriter, r *http.Request) {
	session := app.SessionCheck(w, r)
	if session != nil {
		app.Users.UserStatusOffline(session.UserName)
	}
	if session != nil {
		err := app.Posts.DeleteSessionByUserId(session.UserID)
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

// Open messages page with render all info about it

func (app *Application) Messages(w http.ResponseWriter, r *http.Request) {
	session := app.SessionCheck(w, r)
	data := app.NewTemplateData(r)
	data = app.SetUserName(data, session)
	users, err := app.Posts.GetAllUsers(data.UserName)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	chats, err := app.Users.GetMyChats(data.UserName)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	form := MessagesStruct{
		StartedChats: chats,
		ChatWith:     ChatWith{},
		Users:        users,
	}
	data.Form = form
	data, err = app.NotifCountSet(data)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	app.Render(w, http.StatusOK, "messages.html", data, r)
}

// Open exist chat or creating new function

func (app *Application) OpenChat(w http.ResponseWriter, r *http.Request) {
	chatUserName := validator.ValidChatUserName(r)
	session := app.SessionCheck(w, r)
	data := app.NewTemplateData(r)
	data = app.SetUserName(data, session)
	chatWith := ChatWith{}
	history := []models.MessageStruct{}
	if chatUserName != "" {
		chatUserPhoto, err := app.Users.GetPhotoByUserName(chatUserName)
		if err != nil {
			app.ServerError(w, err, r)
			return
		}
		history, err = app.Users.GetHistoryOfChat(data.UserName, chatUserName)
		if err != nil {
			app.ServerError(w, err, r)
			return
		}
		isOnline := app.Users.GetUserStatus(chatUserName) == 1
		chatWith = ChatWith{
			With:       chatUserName,
			WithPhoto:  chatUserPhoto,
			WithStatus: isOnline,
			History:    history,
		}
	}
	users, err := app.Posts.GetAllUsers(data.UserName)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	chats, err := app.Users.GetMyChats(data.UserName)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}

	//-----------------------------------------------------------
	data.Form = MessagesStruct{
		StartedChats: chats,
		ChatWith:     chatWith,
		Users:        users,
	}
	data, err = app.NotifCountSet(data)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	app.Render(w, http.StatusOK, "messages.html", data, r)
	err = app.Users.UpdateChatHistory(data.UserName, chatUserName, history)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
}

// Sending message

func (app *Application) SendMessage(w http.ResponseWriter, r *http.Request) {
	chatUserName := validator.ValidChatUserName(r)
	message := validator.ChatMessageCorrector(r.FormValue("message"))
	typeImage := false
	imagePath := ""
	if message == "" {
		file, header, err := r.FormFile("chatImageInput")
		if err != nil {
			return
		}
		imagePath, err = app.uploadImage("chat_images", file, header)
		if err != nil {
			return
		}
		typeImage = true
	}
	session := app.SessionCheck(w, r)
	data := app.NewTemplateData(r)
	data = app.SetUserName(data, session)
	messageStruct := models.MessageStruct{}
	if typeImage {
		messageStruct = models.MessageStruct{
			Message: message,
			Time:    validator.GetCurrentTime(),
			Date:    validator.GetCurrentDate(),
			Type:    "image",
			Path:    imagePath,
			Author:  data.UserName,
		}
	} else {
		messageStruct = models.MessageStruct{
			Message: message,
			Time:    validator.GetCurrentTime(),
			Date:    validator.GetCurrentDate(),
			Type:    "message",
			Path:    "",
			Author:  data.UserName,
		}
	}
	err := app.Users.InsertChat(data.UserName, chatUserName, messageStruct)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	app.OpenChat(w, r)
}

// Activities page open function

func (app *Application) Activity(w http.ResponseWriter, r *http.Request) {
	session := app.SessionCheck(w, r)
	data := app.NewTemplateData(r)
	data = app.SetUserName(data, session)
	activities, err := app.Posts.ActivitiesGet(data.UserName)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	data.Form = activities
	err = app.Posts.MarkAllAsSeen(data.UserName)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	data, err = app.NotifCountSet(data)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	app.Render(w, http.StatusOK, "activity.html", data, r)
}

//Remove Activity Function

func (app *Application) RemoveActivity(w http.ResponseWriter, r *http.Request) {
	id := validator.ValidActivityID(r)
	activity, err := app.Posts.ActivityGet(id)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	app.RemoveActionUndo(activity)
	app.Activity(w, r)
}

func (app *Application) ErrorHandler(w http.ResponseWriter, errorNum int, r *http.Request) {
	data := app.NewTemplateData(r)
	Res := &models.ErrorStruct{
		Status: errorNum,
		Text:   http.StatusText(errorNum),
	}
	cookie, _ := r.Cookie("session_token")
	if cookie != nil {
		token := cookie.Value
		session, errCookie := app.Posts.GetSessionFromToken(token)
		if errCookie != nil {
			app.ServerError(w, nil, r)
			return
		}
		data = app.SetUserName(data, session)
	}
	data.ErrorStruct = Res
	data, err := app.NotifCountSet(data)
	if err != nil {
		app.ServerError(w, err, r)
		return
	}
	err = app.renderErr(w, http.StatusUnprocessableEntity, "error.html", data, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *Application) LikePostHandler(w http.ResponseWriter, r *http.Request) {
	app.ReactionPostHandler(w, r, true)
}
func (app *Application) DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	app.ReactionPostHandler(w, r, false)
}
func (app *Application) LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	app.ReactionCommentHandler(w, r, true)
}
func (app *Application) DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	app.ReactionCommentHandler(w, r, false)
}
