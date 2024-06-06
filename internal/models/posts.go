package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"
	errorhandler "yinyang/internal/errors"
)

type PostCreateForm struct {
	Title      string
	Content    string
	Image      string
	Category   string
	Categories []string
	Error      PostCreateFormError
}

type PostCreateFormError struct {
	TitleError   string
	ContentError string
	ImageError   string
}

type PostCategoriesForm struct {
	Id             int
	CategoriesName string
}

type CategoriesForm struct {
	Game        string
	Films       string
	Programming string
	Anime       string
	Sport       string
}

type Post struct {
	ID              int
	Title           string
	Content         string
	Created         time.Time
	Category        string
	UserName        string
	Image           string
	ProfilePhoto    string
	Comments        []Comment
	Likes           int
	Dislikes        int
	IsAuthenticated bool
}

type Comment struct {
	Id                 int
	Author             string
	AuthorProfilePhoto string
	CContent           string
	PostID             int
	Likes              int
	Dislikes           int
	IsAuthenticated    bool
}

type Activity struct {
	ID       int
	Username string
	Author   string
	Type     string
	Post     *Post
	Comment  *Comment
	Seen     int
	Date     string
}

type Model struct {
	DB *sql.DB
}

func (m *Model) Insert(title, content, category, userName, image string) (int, error) {
	stmt := `INSERT INTO posts (title, content, created, category, user_name, image)
    VALUES (?, ?, datetime('now'), ?, ?, ?)`
	result, err := m.DB.Exec(stmt, title, content, category, userName, image)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *Model) DeletePost(postID int) error {
	_, err := m.DB.Exec(`DELETE FROM posts WHERE id = ?`, postID)
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) DeleteActivity(id int) error {
	_, err := m.DB.Exec(`DELETE FROM Activity WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) ActivityInsert(username, typeOf string, postId, commentId, seen int) error {
	stmt := `INSERT INTO Activity (username, author, type, post, comment, seen, date)
    VALUES (?, ?, ?, ?, ?, ?, ?)`
	date := CurrentTimeDateTime()
	var comment *Comment
	var author string
	post, err := m.Get(postId)
	if err != nil {
		return err
	}
	if commentId != 0 {
		comment, err = m.GetComment(commentId)
		if err != nil {
			return err
		}
		if typeOf == "createcomment" {
			author = post.UserName
		} else {
			author = comment.Author
		}
	} else {
		author = post.UserName
	}
	postByte, err := SavePost(post)
	if err != nil {
		return err
	}
	commentByte, err := SaveComment(comment)
	if err != nil {
		return err
	}
	_, err = m.DB.Exec(stmt, username, author, typeOf, string(postByte), string(commentByte), seen, date)
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) ActivitiesGet(username string) ([]*Activity, error) {
	stmt := `SELECT id, username, author, type, post, comment, seen, date FROM Activity WHERE username = ? OR author = ? ORDER BY id DESC`
	rows, err := m.DB.Query(stmt, username, username)
	if err != nil {
		return nil, err
	}
	activities := []*Activity{}
	for rows.Next() {
		activity := &Activity{}
		var postStr string
		var commentStr string
		err = rows.Scan(&activity.ID, &activity.Username, &activity.Author, &activity.Type, &postStr, &commentStr, &activity.Seen, &activity.Date)
		if err != nil {
			return nil, err
		}
		post, err := GetPost(postStr)
		if err != nil {
			return nil, err
		}
		comment, err := GetComment(commentStr)
		if err != nil {
			return nil, err
		}
		activity.Post = post
		activity.Comment = comment
		activities = append(activities, activity)
	}
	return activities, nil
}

// func (m *Model) GetCommentTitleByID(commentID int) (*Comment, error) {
// 	stmt := `SELECT CContent FROM comments WHERE id = ?`
// 	err := m.DB.QueryRow(stmt, id)
// }

func (m *Model) ActivityGet(id int) (*Activity, error) {
	stmt := `SELECT id, username, type, post, comment, seen, date FROM Activity WHERE id = ?`
	row := m.DB.QueryRow(stmt, id)
	activity := &Activity{}
	var postStr string
	var commentStr string
	err := row.Scan(&activity.ID, &activity.Username, &activity.Type, &postStr, &commentStr, &activity.Seen, &activity.Date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorhandler.ErrNoRecord
		}
		return nil, err
	}
	post, err := GetPost(postStr)
	if err != nil {
		return nil, err
	}
	comment, err := GetComment(commentStr)
	if err != nil {
		return nil, err
	}
	activity.Post = post
	activity.Comment = comment
	return activity, nil
}

func (m *Model) AllActivitiesGet() ([]Activity, error) {
	rows, err := m.DB.Query(`SELECT id, username, type, post, comment, seen, date FROM Activity`)
	if err != nil {
		return nil, err
	}
	activityList := []Activity{}
	for rows.Next() {
		activity := Activity{}
		var postStr string
		var commentStr string
		if err = rows.Scan(&activity.ID, &activity.Username, &activity.Type, &postStr, &commentStr, &activity.Seen, &activity.Date); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errorhandler.ErrNoRecord
			}
			return nil, err
		}
		post, err := GetPost(postStr)
		if err != nil {
			return nil, err
		}
		comment, err := GetComment(commentStr)
		if err != nil {
			return nil, err
		}
		activity.Post = post
		activity.Comment = comment
		activityList = append(activityList, activity)
	}
	return activityList, nil
}

func (m *Model) ActivityGetID(username, typeOf string, postID int) (int, error) {
	stmt := `SELECT id FROM Activity WHERE username = ? AND type = ? AND postId = ?`
	var id int
	err := m.DB.QueryRow(stmt, username, typeOf, postID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *Model) Get(id int) (*Post, error) {
	stmt := `SELECT id, title, content, created, category, user_name, image FROM posts WHERE id = ?`
	row := m.DB.QueryRow(stmt, id)
	post := &Post{}
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.Created, &post.Category, &post.UserName, &post.Image)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorhandler.ErrNoRecord
		}
		return nil, err
	}
	profileImage, err := m.GetUserProfileImage(post.UserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorhandler.ErrNoRecord
		}
		return nil, err
	}
	post.ProfilePhoto = profileImage

	stmt = `SELECT COALESCE(SUM(like), 0), COALESCE(SUM(dislike), 0) FROM post_reactions WHERE post_id = ?`
	row = m.DB.QueryRow(stmt, id)
	err = row.Scan(&post.Likes, &post.Dislikes)
	if err != nil {
		return nil, err
	}

	comments, err := m.GetComments(id)
	if err != nil {
		return nil, err
	}
	post.Comments = comments
	return post, nil
}

func (m *Model) Latest() ([]*Post, error) {
	stmt := `SELECT id, title, content, created, category, user_name, image FROM posts ORDER BY id DESC LIMIT 100`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*Post{}
	for rows.Next() {
		s := &Post{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Category, &s.UserName, &s.Image)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *Model) ByCategory(category string) ([]*Post, error) {
	stmt := `SELECT id, title, content, created, category, user_name, image FROM posts WHERE category = ? ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(stmt, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*Post{}
	for rows.Next() {
		s := &Post{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Category, &s.UserName, &s.Image)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *Model) GetPostsByUser(name string) ([]*Post, error) {
	stmt := `SELECT id, title, content, created, category, user_name, image FROM posts WHERE user_name = ? ORDER BY id DESC`
	rows, err := m.DB.Query(stmt, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*Post{}
	for rows.Next() {
		p := &Post{}
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.Created, &p.Category, &p.UserName, &p.Image)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *Model) GetComments(postID int) ([]Comment, error) {
	commentsQuery := `
		SELECT c.Id, c.CContent, c.Author, c.PostID, COALESCE(SUM(r.like), 0) AS Likes, COALESCE(SUM(r.dislike), 0) AS Dislikes
		FROM comments AS c
		LEFT JOIN comment_reactions AS r ON c.Id = r.comment_id
		WHERE c.PostID = ?
		GROUP BY c.Id, c.CContent, c.Author, c.PostID
	`

	stmt, err := m.DB.Prepare(commentsQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}

	for rows.Next() {
		comment := Comment{}
		err = rows.Scan(&comment.Id, &comment.CContent, &comment.Author, &comment.PostID, &comment.Likes, &comment.Dislikes)
		if err != nil {
			return nil, err
		}
		profileImage, err := m.GetUserProfileImage(comment.Author)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errorhandler.ErrNoRecord
			}
			return nil, err
		}
		comment.AuthorProfilePhoto = profileImage
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (m *Model) PostComment(CommentInput Comment) (int64, error) {
	result, err := m.DB.Exec("INSERT INTO comments (CContent, Author, PostID) VALUES ($1,$2,$3)", CommentInput.CContent, CommentInput.Author, CommentInput.PostID)
	if err != nil {
		return 0, err
	}
	commentId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return commentId, nil
}

func (m *Model) GetComment(commentID int) (*Comment, error) {
	stmt := "SELECT ID, CContent, Author, PostID from comments where Id = ?"
	comment := &Comment{}
	err := m.DB.QueryRow(stmt, commentID).Scan(&comment.Id, &comment.CContent, &comment.Author, &comment.PostID)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (m *Model) RemoveComment(commentID int) error {
	_, err := m.DB.Exec("DELETE FROM comments WHERE Id = ?", commentID)
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) GetPostsByUserReaction(userID int) ([]*Post, error) {
	stmt := `SELECT p.id, p.title, p.content, p.created, p.category, p.user_name, p.image
		FROM posts p
		INNER JOIN post_reactions pr ON p.id = pr.post_id
		WHERE pr.user_id = ? AND pr.like = 1
		ORDER BY p.id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*Post{}
	for rows.Next() {
		p := &Post{}
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.Created, &p.Category, &p.UserName, &p.Image)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *Model) GetUserProfileImage(username string) (string, error) {
	stmt := `SELECT profile_photo FROM Users WHERE name = ?`

	var profile_photo string
	err := m.DB.QueryRow(stmt, username).Scan(&profile_photo)
	if err != nil {
		return "", err
	}

	return profile_photo, nil
}

func CurrentTimeDateTime() string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("02.01.2006 15:04")
	return formattedTime
}

func SavePost(post *Post) ([]byte, error) {
	postData, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}
	return postData, nil
}

func GetPost(messagesData string) (*Post, error) {
	var post *Post
	err := json.Unmarshal([]byte(messagesData), &post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func SaveComment(comment *Comment) ([]byte, error) {
	commentData, err := json.Marshal(comment)
	if err != nil {
		return nil, err
	}
	return commentData, nil
}

func GetComment(messagesData string) (*Comment, error) {
	var comment *Comment
	err := json.Unmarshal([]byte(messagesData), &comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (m *Model) GetUnseenActivityCount(username string) (int, error) {
	stmt := `SELECT COUNT(*) FROM Activity WHERE seen = 0 AND username != ? AND author = ?`
	var count int
	row := m.DB.QueryRow(stmt, username, username)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *Model) MarkAllAsSeen(username string) error {
	stmt := `UPDATE Activity SET seen = 1 WHERE seen = 0 AND username != ? AND author = ?`

	_, err := m.DB.Exec(stmt, username, username)
	if err != nil {
		return err
	}

	stmt2 := `UPDATE Activity SET seen = 1 WHERE seen = 0 AND username = ? AND author = ?`

	_, err = m.DB.Exec(stmt2, username, username)
	if err != nil {
		return err
	}

	return nil
}
