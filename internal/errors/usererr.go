package errorhandler

import "errors"

var (
	ErrUserAgrement = errors.New("Accept user agreement!")
	ErrShortUsername = errors.New("Username is too short!")
	ErrLongUsername = errors.New("Username is too long!")
	ErrUsernameStart = errors.New("Username can`t start with number!")

	ErrLongUsernameSymbols = errors.New("User name can`t be longer than 10 symobols")
	ErrLongEmailSymbols = errors.New("Email can`t be longer than 30 symobols")
	ErrQuery = errors.New("query error")
	ErrAlreadyExistUsername = errors.New("User name already exists")
	ErrAlreadyExistEmail = errors.New("Email already register")

	ErrEnterCorrectEmail = errors.New("Please enter correct email!")
	ErrPasswordMismatch = errors.New("Password mismatch!")
	ErrLowPassword = errors.New("Min size of password is 7 Use one of them [$#%!?.*]")
	ErrAuthServer = errors.New("Server error! Please try later")

	ErrEmailExist = errors.New("Email is not exists!")
	ErrIncorrectPassword = errors.New("Password is incorrect!")
)
