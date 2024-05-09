package models

import (
	"database/sql"
	"errors"
)

type ReactionModel struct {
	DB *sql.DB
}

func (r *ReactionModel) LikePost(userID, postID int) (bool, error) {
	stmt := `SELECT like, dislike FROM post_reactions WHERE post_id = ? AND user_id = ?`
	var like, dislike int
	var boly bool
	err := r.DB.QueryRow(stmt, postID, userID).Scan(&like, &dislike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = r.DB.Exec(`INSERT INTO post_reactions (post_id, user_id, like, dislike) VALUES (?, ?, ?, ?)`, postID, userID, 1, 0)
			boly = true
			if err != nil {
				return boly, err
			}
		} else {
			return boly, err
		}
	}
	if like == 1 {
		boly = false
		_, err := r.DB.Exec(`DELETE FROM post_reactions WHERE post_id = ? AND user_id = ?`, postID, userID)
		if err != nil {
			return boly, err
		}
		err = r.DeleteRectionAcivity(userID, postID, "reactionpost")
		if err != nil {
			return boly, err
		}
	} else if dislike == 1 {
		boly = true
		_, err := r.DB.Exec(`UPDATE post_reactions SET dislike = ? WHERE post_id = ? AND user_id = ?`, 0, postID, userID)
		if err != nil {
			return boly, err
		}
		_, err = r.DB.Exec(`UPDATE post_reactions SET like = ? WHERE post_id = ? AND user_id = ?`, 1, postID, userID)
		if err != nil {
			return boly, err
		}
		err = r.DeleteRectionAcivity(userID, postID, "reactionpost")
		if err != nil {
			return boly, err
		}
	}
	return boly, nil
}

func (r *ReactionModel) DislikePost(userID, postID int) (bool, error) {
	stmt := `SELECT like, dislike FROM post_reactions WHERE post_id = ? AND user_id = ?`
	var like, dislike int
	var boly bool
	err := r.DB.QueryRow(stmt, postID, userID).Scan(&like, &dislike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			boly = true
			_, err = r.DB.Exec(`INSERT INTO post_reactions (post_id, user_id, like, dislike) VALUES (?, ?, ?, ?)`, postID, userID, 0, 1)
			if err != nil {
				return boly, err
			}
		} else {
			return boly, err
		}
	}
	if dislike == 1 {
		boly = false
		_, err := r.DB.Exec(`DELETE FROM post_reactions WHERE post_id = ? AND user_id = ?`, postID, userID)
		if err != nil {
			return boly, err
		}
		err = r.DeleteRectionAcivity(userID, postID, "reactionpost")
		if err != nil {
			return boly, err
		}
	} else if like == 1 {
		boly = true
		_, err := r.DB.Exec(`UPDATE post_reactions SET like = ? WHERE post_id = ? AND user_id = ?`, 0, postID, userID)
		if err != nil {
			return boly, err
		}
		_, err = r.DB.Exec(`UPDATE post_reactions SET dislike = ? WHERE post_id = ? AND user_id = ?`, 1, postID, userID)
		if err != nil {
			return boly, err
		}
		err = r.DeleteRectionAcivity(userID, postID, "reactionpost")
		if err != nil {
			return boly, err
		}
	}
	return boly, nil
}

func (r *ReactionModel) LikeComment(userID, commentID, postID int) (bool, error) {
	stmt := `SELECT like, dislike FROM comment_reactions WHERE comment_id = ? AND user_id = ?`
	var like, dislike int
	var boly bool
	err := r.DB.QueryRow(stmt, commentID, userID).Scan(&like, &dislike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			boly = true
			_, err = r.DB.Exec(`INSERT INTO comment_reactions (comment_id, user_id, like, dislike) VALUES (?, ?, ?, ?)`, commentID, userID, 1, 0)
			if err != nil {
				return boly, err
			}
		} else {
			return boly, err
		}
	}
	if like == 1 {
		boly = false
		_, err := r.DB.Exec(`DELETE FROM comment_reactions WHERE comment_id = ? AND user_id = ?`, commentID, userID)
		if err != nil {
			return boly, err
		}
		err = r.DeleteRectionAcivity(userID, postID, "reactioncomment")
		if err != nil {
			return boly, err
		}
	} else if dislike == 1 {
		boly = true
		_, err := r.DB.Exec(`UPDATE comment_reactions SET dislike = ? WHERE comment_id = ? AND user_id = ?`, 0, commentID, userID)
		if err != nil {
			return boly, err
		}
		_, err = r.DB.Exec(`UPDATE comment_reactions SET like = ? WHERE comment_id = ? AND user_id = ?`, 1, commentID, userID)
		if err != nil {
			return boly, err
		}
		err = r.DeleteRectionAcivity(userID, postID, "reactioncomment")
		if err != nil {
			return boly, err
		}
	}
	return boly, nil
}

func (r *ReactionModel) DislikeComment(userID, commentID, postID int) (bool, error) {
	stmt := `SELECT like, dislike FROM comment_reactions WHERE comment_id = ? AND user_id = ?`
	var like, dislike int
	var boly bool
	err := r.DB.QueryRow(stmt, commentID, userID).Scan(&like, &dislike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			boly = true
			_, err = r.DB.Exec(`INSERT INTO comment_reactions (comment_id, user_id, like, dislike) VALUES (?, ?, ?, ?)`, commentID, userID, 0, 1)
			if err != nil {
				return boly, err
			}
		} else {
			return boly, err
		}
	}
	if dislike == 1 {
		boly = false
		_, err := r.DB.Exec(`DELETE FROM comment_reactions WHERE comment_id = ? AND user_id = ?`, commentID, userID)
		if err != nil {
			return boly, err
		}
		err = r.DeleteRectionAcivity(userID, postID, "reactioncomment")
		if err != nil {
			return boly, err
		}
	} else if like == 1 {
		boly = true
		_, err := r.DB.Exec(`UPDATE comment_reactions SET like = ? WHERE comment_id = ? AND user_id = ?`, 0, commentID, userID)
		if err != nil {
			return boly, err
		}
		_, err = r.DB.Exec(`UPDATE comment_reactions SET dislike = ? WHERE comment_id = ? AND user_id = ?`, 1, commentID, userID)
		if err != nil {
			return boly, err
		}
		err = r.DeleteRectionAcivity(userID, postID, "reactioncomment")
		if err != nil {
			return boly, err
		}
	}
	return boly, nil
}

func (r *ReactionModel) DeleteRectionAcivity(userID, postID int, typeOf string) error {
	username, err := r.GetUserNameById(userID)
	if err != nil {
		return err
	}
	actID, err := r.ActivityGetID(username, typeOf, postID)
	if err != nil {
		return err
	}
	err = r.DeleteActivity(actID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReactionModel) RemovePostReaction(userID, postID int) error {
	stmt := `SELECT like, dislike FROM post_reactions WHERE post_id = ? AND user_id = ?`
	var like, dislike int
	err := r.DB.QueryRow(stmt, postID, userID).Scan(&like, &dislike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = r.DB.Exec(`INSERT INTO post_reactions (post_id, user_id, like, dislike) VALUES (?, ?, ?, ?)`, postID, userID, 0, 1)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if like == 1 || dislike == 1 {
		_, err := r.DB.Exec(`DELETE FROM post_reactions WHERE post_id = ? AND user_id = ?`, postID, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReactionModel) RemoveCommentReaction(userID, commentID int) error {
	stmt := `SELECT like, dislike FROM comment_reactions WHERE comment_id = ? AND user_id = ?`
	var like, dislike int
	err := r.DB.QueryRow(stmt, commentID, userID).Scan(&like, &dislike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = r.DB.Exec(`INSERT INTO comment_reactions (comment_id, user_id, like, dislike) VALUES (?, ?, ?, ?)`, commentID, userID, 1, 0)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if like == 1 || dislike == 1 {
		_, err := r.DB.Exec(`DELETE FROM comment_reactions WHERE comment_id = ? AND user_id = ?`, commentID, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

// func (r *ReactionModel) GetPostId(username string, commentID int) (int, error) {
// 	var postId int
// 	stmt := `SELECT PostID FROM comments WHERE Id = ? AND Author = ?`
// 	err := r.DB.QueryRow(stmt, commentID, username).Scan(&postId)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return postId, nil
// }

func (r *ReactionModel) GetCommentPostID(commentID int) (int, error) {
	var postID int
	stmt := `SELECT PostID FROM comments WHERE Id = ?`
	err := r.DB.QueryRow(stmt, commentID).Scan(&postID)
	if err != nil {
		return postID, err
	}
	return postID, nil
}

func (r *ReactionModel) ActivityGetID(username, typeOf string, postID int) (int, error) {
	stmt := `SELECT id FROM Activity WHERE username = ? AND type = ?`
	var id int
	err := r.DB.QueryRow(stmt, username, typeOf).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ReactionModel) GetUserNameById(id int) (string, error) {
	stmt := `SELECT name FROM Users WHERE ID = ?`

	var name string
	err := r.DB.QueryRow(stmt, id).Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (r *ReactionModel) DeleteActivity(id int) error {
	_, err := r.DB.Exec(`DELETE FROM Activity WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return nil
}

// func (r *ReactionModel) GetPost(id int) (*Post, error) {
// 	stmt := `SELECT id, title, content, created, category, user_name, image FROM posts WHERE id = ?`
// 	row := r.DB.QueryRow(stmt, id)
// 	post := &Post{}
// 	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.Created, &post.Category, &post.UserName, &post.Image)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, errorhandler.ErrNoRecord
// 		}
// 		return nil, err
// 	}
// 	profileImage, err := r.GetUserProfileImage(post.UserName)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, errorhandler.ErrNoRecord
// 		}
// 		return nil, err
// 	}
// 	post.ProfilePhoto = profileImage

// 	stmt = `SELECT COALESCE(SUM(like), 0), COALESCE(SUM(dislike), 0) FROM post_reactions WHERE post_id = ?`
// 	row = r.DB.QueryRow(stmt, id)
// 	err = row.Scan(&post.Likes, &post.Dislikes)
// 	if err != nil {
// 		return nil, err
// 	}

// 	comments, err := r.GetComments(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	post.Comments = comments
// 	return post, nil
// }

// func (r *ReactionModel) GetComments(postID int) ([]Comment, error) {
// 	commentsQuery := `
// 		SELECT c.Id, c.CContent, c.Author, c.PostID, COALESCE(SUM(r.like), 0) AS Likes, COALESCE(SUM(r.dislike), 0) AS Dislikes
// 		FROM comments AS c
// 		LEFT JOIN comment_reactions AS r ON c.Id = r.comment_id
// 		WHERE c.PostID = ?
// 		GROUP BY c.Id, c.CContent, c.Author, c.PostID
// 	`

// 	stmt, err := r.DB.Prepare(commentsQuery)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer stmt.Close()

// 	rows, err := stmt.Query(postID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	comments := []Comment{}

// 	for rows.Next() {
// 		comment := Comment{}
// 		err = rows.Scan(&comment.Id, &comment.CContent, &comment.Author, &comment.PostID, &comment.Likes, &comment.Dislikes)
// 		if err != nil {
// 			return nil, err
// 		}
// 		profileImage, err := r.GetUserProfileImage(comment.Author)
// 		if err != nil {
// 			if errors.Is(err, sql.ErrNoRows) {
// 				return nil, errorhandler.ErrNoRecord
// 			}
// 			return nil, err
// 		}
// 		comment.AuthorProfilePhoto = profileImage
// 		comments = append(comments, comment)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return comments, nil
// }

// func (r *ReactionModel) GetUserProfileImage(username string) (string, error) {
// 	stmt := `SELECT profile_photo FROM Users WHERE name = ?`

// 	var profile_photo string
// 	err := r.DB.QueryRow(stmt, username).Scan(&profile_photo)
// 	if err != nil {
// 		return "", err
// 	}

// 	return profile_photo, nil
// }
