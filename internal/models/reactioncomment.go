package models

import (
	"database/sql"
)

type ReactionCommentModel struct {
	DB *sql.DB
}

func (r *ReactionCommentModel) LikeComment(user_id, comment_id int) error {
	stmt := `SELECT name FROM Users WHERE id =?`
	var name string
	err := r.DB.QueryRow(stmt, user_id).Scan(&name)
	if err != nil {
		return err
	}
	return nil
}
