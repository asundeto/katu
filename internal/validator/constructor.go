package validator

import (
	"strings"
	"yinyang/internal/models"
	"yinyang/internal/template"
)

type Constructor struct {
	Users *models.UserModel
}

func GetCats(form *models.CategoriesForm) string {
	var cats []string
	if form.Game != "" {
		cats = append(cats, form.Game)
	}
	if form.Films != "" {
		cats = append(cats, form.Films)
	}
	if form.Programming != "" {
		cats = append(cats, form.Programming)
	}
	if form.Anime != "" {
		cats = append(cats, form.Anime)
	}
	if form.Sport != "" {
		cats = append(cats, form.Sport)
	}
	return strings.Join(cats, " ")
}

func AddCategoriesToData(data *template.TemplateData) *template.TemplateData {
	categForm := []models.PostCategoriesForm{}
	categStr := StringToArray(data.Post.Category)
	for i := 0; i < len(categStr); i++ {
		oneCat := models.PostCategoriesForm{
			Id:             i + 1,
			CategoriesName: categStr[i],
		}
		categForm = append(categForm, oneCat)
	}
	data.PostCategoriesForm = categForm
	return data
}

func (c *Constructor) PostCreateCheck(form *models.PostCreateForm) *models.PostCreateForm {
	if ValidPostTitle(form.Title) != "" {
		form.Error.TitleError = ValidPostTitle(form.Title)
	}
	if ValidPostContent(form.Content) != "" {
		form.Error.ContentError = ValidPostContent(form.Content)
	}
	return form
}

func CommentForm(comment, username string, id int) models.Comment {
	comment = strings.Replace(comment, "\n", "<br>", -1)

	commentInput := models.Comment{
		Author:   username,
		CContent: comment,
		PostID:   id,
	}
	return commentInput
}
