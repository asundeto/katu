package handlers

import (
	"regexp"
	"strings"
)

func IsEmailValid(s string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(s)
}

func isPasswordValid(s string) bool {
	if len(s) < 7 {
		return false
	}
	passRegex := regexp.MustCompile(`[a-z]`)
	passRegex2 := regexp.MustCompile(`[A-Z]`)
	passRegex4 := regexp.MustCompile(`[0-9]`)
	passRegex3 := regexp.MustCompile(`[$#%!?.*]`)
	if passRegex.MatchString(s) && passRegex2.MatchString(s) && passRegex3.MatchString(s) && passRegex4.MatchString(s) {
		return true
	}
	return false
}

func stringToArray(s string) []string {
	result := strings.Fields(s)
	return result
}

func getCats(form *CategoriesForm) string {
	var cats string
	if form.Game != "" {
		cats += form.Game + " "
	}
	if form.Films != "" {
		cats += form.Films + " "
	}
	if form.Programming != "" {
		cats += form.Programming + " "
	}
	if form.Anime != "" {
		cats += form.Anime + " "
	}
	if form.Sport != "" {
		cats += form.Sport + " "
	}

	return cats
}

func countLines(comment string) int {
	lines := strings.Split(comment, "\n")
	return len(lines)
}

func addCategoriesToData(data *TemplateData) *TemplateData {
	categForm := []PostCategoriesForm{}
	categStr := stringToArray(data.Post.Category)
	for i := 0; i < len(categStr); i++ {
		oneCat := PostCategoriesForm{
			Id:             i + 1,
			CategoriesName: categStr[i],
		}
		categForm = append(categForm, oneCat)
	}
	data.PostCategoriesForm = categForm
	return data
}

func userNameCheck(userName string) string {
	if len(userName) < 4 {return "Username is too short!"}
	if len(userName) > 10 {return "Username is too long!"}
	if userName[0] >= '0' && userName[0] <= '9' {return "Username can`t start with number!"}
	return ""
}
