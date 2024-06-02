package models

import (
	"database/sql"
	"errors"
	"time"

	errorhandler "yinyang/internal/errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	ProfilePhoto   string
	Created        time.Time
	Error          error
	Access         bool
	Online         bool
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password, profile_photo string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, profile_photo, online, created)
		VALUES(?, ?, ?, ?, ?, datetime('now'))`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword), profile_photo, 0)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) EmailExist(email string) error {
	var count int
	stmt := "SELECT COUNT(*) FROM users WHERE email = ? OR name = ?"
	err := m.DB.QueryRow(stmt, email, email).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errorhandler.ErrEmailExist
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = ? OR name = ?"
	err := m.DB.QueryRow(stmt, email, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errorhandler.ErrEmailExist
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, errorhandler.ErrIncorrectPassword
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (u *UserModel) GetUserNameByEmail(email string) (string, error) {
	stmt := `SELECT name FROM Users WHERE email = ? OR name = ?`

	var name string
	err := u.DB.QueryRow(stmt, email, email).Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (u *UserModel) GetUserNameById(id int) (string, error) {
	stmt := `SELECT name FROM Users WHERE ID = ?`

	var name string
	err := u.DB.QueryRow(stmt, id).Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (u *UserModel) GetIDByUserName(name string) (int, error) {
	stmt := `SELECT ID FROM Users WHERE name = ?`

	var id int
	err := u.DB.QueryRow(stmt, name).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *Model) GetAllUsers(username string) ([]*User, error) {
	stmt := `SELECT id, name, profile_photo, online FROM Users`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		s := &User{}
		var online int
		err = rows.Scan(&s.ID, &s.Name, &s.ProfilePhoto, &online)
		if err != nil {
			return nil, err
		}
		s.Online = online == 1
		if s.Name == username {
			continue
		}
		users = append(users, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
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

func (m *UserModel) UserStatusOffline(username string) error {
	updateStmt := "UPDATE Users SET online = ? WHERE name = ?"
	_, err := m.DB.Exec(updateStmt, 0, username)
	if err != nil {
		return err
	}
	return nil
}
func (m *UserModel) UserStatusOnline(username string) error {
	updateStmt := "UPDATE Users SET online = ? WHERE name = ?"
	_, err := m.DB.Exec(updateStmt, 1, username)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) GetUserStatus(username string) int {
	var online int
	query := "SELECT online FROM Users WHERE name = ?"
	err := m.DB.QueryRow(query, username).Scan(&online)
	if err != nil {
		return 0
	}
	return online
}

func (m *UserModel) Duplicates(email, name string) error {
	if len(name) > 10 {
		return errorhandler.ErrLongUsernameSymbols
	}
	if len(email) > 30 {
		return errorhandler.ErrLongEmailSymbols
	}

	user := User{Email: email, Name: name}
	var count, count2 int

	stmt2 := "SELECT COUNT(*) FROM users WHERE name = ?"
	stmt := "SELECT COUNT(*) FROM users WHERE email = ?"

	errUsername := m.DB.QueryRow(stmt2, user.Name).Scan(&count2)
	if errUsername != nil {
		return errorhandler.ErrQuery
	}
	errEmail := m.DB.QueryRow(stmt, user.Email).Scan(&count)
	if errEmail != nil {
		return errorhandler.ErrQuery
	}

	if count2 > 0 {
		return errorhandler.ErrAlreadyExistUsername
	}
	if count > 0 {
		return errorhandler.ErrAlreadyExistEmail
	}

	return nil
}

func (u *UserModel) GetPhotoByUserName(username string) (string, error) {
	stmt := `SELECT profile_photo FROM Users WHERE name = ?`
	var profilePhoto string
	err := u.DB.QueryRow(stmt, username).Scan(&profilePhoto)
	if err != nil {
		return "", nil
	}
	return profilePhoto, nil
}
