package models

import (
	"database/sql"
	"errors"
	"time"
	errorhandler "yinyang/internal/errors"

	"github.com/google/uuid"
)

type ErrorStruct struct {
	Status int
	Text   string
}

type Session struct {
	UserID         int
	UserName       string
	Token          string
	ExpirationDate time.Time
}

func (m *Model) CreateSession(userId int, userName string) (string, time.Time, error) {
	token := uuid.NewString()
	date := time.Now().Add(24 * time.Hour)
	stmt := `INSERT INTO Sessions ( user_id, user_name, token, expiration_date) 
			VALUES(?,?,?,?)`
	_, err := m.DB.Exec(stmt, userId, userName, token, date)
	if err != nil {
		return "", date, err
	}
	return token, date, nil
}

func (m *Model) GetSessionFromToken(token string) (*Session, error) {
	stmt := `SELECT user_id, user_name, token, expiration_date FROM Sessions  where token = ?`

	row := m.DB.QueryRow(stmt, token)
	session := &Session{}
	err := row.Scan(&session.UserID, &session.UserName, &session.Token, &session.ExpirationDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorhandler.ErrNoRecord
		}
		return nil, err
	}
	return session, nil
}

func (m *Model) DeleteSessionByUserId(userId int) error {
	_, err := m.DB.Exec("DELETE FROM Sessions WHERE user_id = $1", userId)
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) DeleteExpiredSessions() error {
	query := "SELECT user_name FROM Sessions WHERE expiration_date < $1"
	rows, err := m.DB.Query(query, time.Now())
	if err != nil {
		return err
	}
	defer rows.Close()

	var usersToDelete []string
	for rows.Next() {
		var userName string
		if err := rows.Scan(&userName); err != nil {
			return err
		}
		usersToDelete = append(usersToDelete, userName)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Now that we have all users to delete, we can close the previous rows
	if err := rows.Close(); err != nil {
		return err
	}

	// Perform offline status update for all users
	for _, userName := range usersToDelete {
		if err := m.UserStatusOffline(userName); err != nil {
			return err
		}
	}

	// Now delete the expired sessions
	_, err = m.DB.Exec("DELETE FROM Sessions WHERE expiration_date < $1", time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) UserStatusOffline(username string) error {
	updateStmt := "UPDATE Users SET online = ? WHERE name = ?"
	_, err := m.DB.Exec(updateStmt, 0, username)
	if err != nil {
		return err
	}
	return nil
}

// func (m *Model) DeleteExpiredSessions() error {
// 	query := "SELECT user_name FROM Sessions WHERE expiration_date < $1"
// 	fmt.Println("1")
// 	rows, err := m.DB.Query(query, time.Now())
// 	if err != nil {
// 		return err
// 	}

// 	defer rows.Close()
// 	fmt.Println("2")
// 	for rows.Next() {
// 		var userName string
// 		if err := rows.Scan(&userName); err != nil {
// 			return err
// 		}
// 		err = m.UserStatusOffline(userName)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	if err := rows.Err(); err != nil {
// 		return err
// 	}
// 	fmt.Println("3")
// 	fmt.Println("4")
// 	_, err = m.DB.Exec("DELETE FROM Sessions WHERE expiration_date < $1", time.Now())
// 	if err != nil {
// 		fmt.Println("THIRD")
// 		return err
// 	}
// 	fmt.Println("5")
// 	return nil
// }

// func (m *Model) DeleteExpiredSessions() error {
// 	query := "SELECT user_name FROM Sessions WHERE expiration_date < $1"
// 	fmt.Println("1")
// 	rows, err := m.DB.Query(query, time.Now())
// 	if err != nil {
// 		return err
// 	}

// 	defer rows.Close()
// 	fmt.Println("2")
// 	for rows.Next() {
// 		var userName string
// 		if err := rows.Scan(&userName); err != nil {
// 			return err
// 		}
// 		err = m.UserStatusOffline(userName)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	if err := rows.Err(); err != nil {
// 		return err
// 	}
// 	fmt.Println("3")
// 	// err := m.DB.QueryRow(query, time.Now()).Scan(&username)
// 	// if err != nil {
// 	// 	fmt.Println("FIRST")
// 	// 	return err
// 	// }
// 	// err = m.UserStatusOffline(username)
// 	// if err != nil {
// 	// 	fmt.Println("SECOND")
// 	// 	return err
// 	// }
// 	fmt.Println("4")
// 	_, err = m.DB.Exec("DELETE FROM Sessions WHERE expiration_date < $1", time.Now())
// 	if err != nil {
// 		fmt.Println("THIRD")
// 		return err
// 	}
// 	fmt.Println("5")
// 	return nil
// }
