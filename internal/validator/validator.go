package validator

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	errorhandler "yinyang/internal/errors"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*\\.(?:com|ru|su|net|org|gov|edu|co\\.uk)$")

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func StringToArray(s string) []string {
	result := strings.Fields(s)
	return result
}

func CountLines(comment string) int {
	lines := strings.Split(comment, "\n")
	return len(lines)
}

func IsEmailValid(s string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(s)
}

func IsPasswordValid(s string) bool {
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

func UserNameCheck(userName string) error {
	if len(userName) < 4 {
		return errorhandler.ErrShortUsername
	}
	if len(userName) > 10 {
		return errorhandler.ErrLongUsername
	}
	if userName[0] >= '0' && userName[0] <= '9' {
		return errorhandler.ErrUsernameStart
	}
	return nil
}

func RandomStr(s string) (string, error) {
	fileExt := getFileExtension(s)
	err := fileExtensionCheck(fileExt)
	if err != nil {
		return "", err
	}
	// Initialize random number generator
	rand.Seed(time.Now().UnixNano())
	// Define symbols to choose from
	symbols := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	// Generate random string of 5 symbols
	result := make([]byte, 5)
	for i := range result {
		result[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(result) + fileExt, nil
}

func getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) == 1 {
		return ""
	}
	return "." + parts[len(parts)-1]
}

func fileExtensionCheck(s string) error {
	allowedExts := []string{".png", ".jpg", ".jpeg", ".gif", ".PNG", ".JPG", ".JPEG", ".GIF"}
	valid := false
	for _, ext := range allowedExts {
		if s == ext {
			valid = true
			break
		}
	}
	if !valid {
		return errorhandler.ErrPostImageExtension
	}
	return nil
}

func ValidPostId(r *http.Request) int {
	idStr := strings.TrimPrefix(r.URL.Path, "/post/view/")
	if idStr[0] == '0' {
		return 0
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		return 0
	}
	return id
}

func ValidUserName(r *http.Request) string {
	username := strings.TrimPrefix(r.URL.Path, "/user/profile/")
	return username
}

func ValidChatUserName(r *http.Request) string {
	username := strings.TrimPrefix(r.URL.Path, "/chat/")
	if username == "" {
		return ""
	}
	if username[0] == '0' {
		return ""
	}
	return username
}

func ValidActivityID(r *http.Request) int {
	idStr := strings.TrimPrefix(r.URL.Path, "/remove/activity/")
	if idStr[0] == '0' {
		return 0
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		return 0
	}
	return id
}

func ValidChatUserName2(r *http.Request) string {
	username := strings.TrimPrefix(r.URL.Path, "/sendmessage/")
	if username == "" {
		return ""
	}
	if username[0] == '0' {
		return ""
	}
	return username
}

func ValidPostTitle(s string) string {
	if len(s) < 3 || len(s) > 20 {
		return errorhandler.ErrPostTitle.Error()
	}
	return ""
}

func ValidPostContent(s string) string {
	if len(s) < 5 || len(s) > 400 {
		return errorhandler.ErrPostContent.Error()
	}
	return ""
}

func ToUpperFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	first := strings.ToUpper(string(s[0]))
	return first + s[1:]
}

func MatchesCategory(postCategory, formCategory string) bool {
	return formCategory != "" && postCategory != "" && strings.Contains(postCategory, formCategory)
}

func PathExists(path string) error {
	uploadsPath := filepath.Join("ui", "static", path)

	// Check if the folder already exists
	if _, err := os.Stat(uploadsPath); os.IsNotExist(err) {
		// Create the folder if it doesn't exist
		err := os.MkdirAll(uploadsPath, 0755)
		if err != nil {
			fmt.Println("Error creating folder:", err)
			return err
		}
	}
	return nil
}

func GetCurrentTime() string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("15:04")
	return formattedTime
}

func GetCurrentDate() string {
	currentTime := time.Now()
	formattedDate := currentTime.Format("02.01.2006 15:04")
	return formattedDate
}

func ChatMessageCorrector(input string) string {
	var result strings.Builder

	for i, char := range input {
		result.WriteRune(char)
		if (i+1)%40 == 0 {
			result.WriteString("<br>")
		}
	}

	return result.String()
}
