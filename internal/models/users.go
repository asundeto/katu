package models

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	ProfilePhoto   string
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password, profile_photo string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, profile_photo, created)
		VALUES(?, ?, ?, ?, datetime('now'))`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword), profile_photo)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) EmailExist(email string) error {
	var count int
	stmt := "SELECT COUNT(*) FROM users WHERE email = ?"
	err := m.DB.QueryRow(stmt, email).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrInvalidCredentials
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (u *UserModel) GetUserNameByEmail(email string) (string, error) {
	stmt := `SELECT name FROM Users WHERE email = ?`

	var name string
	err := u.DB.QueryRow(stmt, email).Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (u *UserModel) GetEmailAndPhotoByUserName(username string) (string, string, error) {
	stmt := `SELECT email FROM Users WHERE name = ?`
	stmt2 := `SELECT profile_photo FROM Users WHERE name = ?`
	var email string
	err := u.DB.QueryRow(stmt, username).Scan(&email)
	if err != nil {
		return "", "", err
	}
	var profilePhoto string
	err2 := u.DB.QueryRow(stmt2, username).Scan(&profilePhoto)
	if err2 != nil {
		return "", "", err2
	}

	return email, profilePhoto, nil
}

func (m *UserModel) ChangeUserProfilePhoto(username, profilePhoto string) error {
	updateStmt := "UPDATE Users SET profile_photo = ? WHERE name = ?"
	_, err := m.DB.Exec(updateStmt, profilePhoto, username)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) Duplicates(email, name string) string {
	if len(name) > 10 {
		return "User name can`t be longer than 10 symobols"
	}
	if len(email) > 30 {
		return "Email can`t be longer than 30 symobols"
	}

	user := User{Email: email, Name: name}
	var count, count2 int

	stmt2 := "SELECT COUNT(*) FROM users WHERE name = ?"
	stmt := "SELECT COUNT(*) FROM users WHERE email = ?"

	errUsername := m.DB.QueryRow(stmt2, user.Name).Scan(&count2)
	if errUsername != nil {
		return "query error"
	}
	errEmail := m.DB.QueryRow(stmt, user.Email).Scan(&count)
	if errEmail != nil {
		return "query error"
	}

	if count2 > 0 {
		return "User name already exists"
	}
	if count > 0 {
		return "Email already register"
	}

	return ""
}
