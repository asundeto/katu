package errorhandler

import (
	"errors"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEntry     = errors.New("duplicate entry")
	ErrNoComments         = errors.New("models: no comments found for post")
	ErrZeroCode           = errors.New("Temporary token is invalid")
	ErrGoogleInfo         = errors.New("The information received from Google")
	ErrGitHUBInfo         = errors.New("The information received from Github")
	ErrServerError        = errors.New("Server Error")
	ErrComment            = errors.New("Please enter correct value")
	ErrIncorrectValue     = errors.New("Incorrect format!")
	ErrPostTitle          = errors.New("The title length must be from 3 to 20 characters!")
	ErrPostContent        = errors.New("The content length must be from 5 to 400 characters!")
	ErrPostImageExtension = errors.New("The image can have the following formats: .png .jpg .gif")
	ErrPostImageSize      = errors.New("Image size is more than 20mb!")
	ErrUploadImage        = errors.New("Please upload a correct image!")
)

var (
	ErrTooManyRequests = errors.New("Too many request")
)
