package template

import (
	"yinyang/internal/models"
)

type TemplateData struct {
	CurrentYear        int
	Post               *models.Post
	Posts              []*models.Post
	Form               any
	IsAuthenticated    bool
	Category           string
	Categories         []string
	Comments           []models.Comment
	ErrorStruct        *models.ErrorStruct
	CommentError       bool
	UserName           string
	PostCategoriesForm []models.PostCategoriesForm
	Notifications     *models.Notifications
}
