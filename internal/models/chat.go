package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

type Chat struct {
	ID               int
	FirstPersonName  string
	SecondPersonName string
	History          []string
}

type StartedChat struct {
	With                string
	WithPhoto           string
	WithStatus          bool
	LastMessage         string
	LastMessageTime     string
	LastMessageDate     string
	UnseenMessagesCount int
}

type MessageStruct struct {
	ID      int
	Time    string
	Date    string
	Message string
	Type    string //image //string //document
	Path    string
	Author  string
	Seen    bool
}

type Notifications struct {
	ActionCount  int
	MessageCount int
}

type ChatModel struct {
	DB *sql.DB
}

var kazakhMonths = []string{
	"Қаңтар", "Ақпан", "Наурыз", "Сәуір", "Мамыр", "Маусым",
	"Шілде", "Тамыз", "Қыркүйек", "Қазан", "Қараша", "Желтоқсан",
}

func (m *UserModel) InsertChat(firstPerson, secondPerson string, message MessageStruct) error {
	check, err := m.ChatExists(firstPerson, secondPerson)
	if err != nil {
		fmt.Println("----------------------- InsertChat ERROR 1 --------------------------")
		return err
	}
	var chatHistory []MessageStruct
	var historyBytes []byte
	if check {
		chatId, err := m.GetChatID(firstPerson, secondPerson)
		if err != nil {
			fmt.Println("----------------------- InsertChat ERROR 2 --------------------------")
			return err
		}
		chatHistory, err = m.GetChatHistory(chatId)
		if err != nil {
			fmt.Println("----------------------- InsertChat ERROR 3 --------------------------")
			return err
		}

		chatHistory = m.AddDateToMessage(chatHistory, message)

		updateStmt := "UPDATE Chats SET history = ? WHERE id = ?"
		historyBytes, err = SaveMessages(chatHistory)
		if err != nil {
			fmt.Println("----------------------- InsertChat ERROR 4 --------------------------")
			return err
		}
		_, err = m.DB.Exec(updateStmt, string(historyBytes), chatId)
		if err != nil {
			fmt.Println("----------------------- InsertChat ERROR 5 --------------------------")
			return err
		}
		return nil
	}
	chatHistory = append(chatHistory, message)
	historyBytes, err = SaveMessages(chatHistory)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO Chats (first_person_name, second_person_name, history)
		VALUES(?, ?, ?)`
	_, err = m.DB.Exec(stmt, firstPerson, secondPerson, string(historyBytes))
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) AddDateToMessage(chatHistory []MessageStruct, message MessageStruct) []MessageStruct {
	boly, _ := CompareDates(chatHistory[len(chatHistory)-1].Date, message.Date)
	if boly {
		kz, _ := ConvertToKazakhCalendar(message.Date)
		dateMessage := MessageStruct{
			Type: "calendar",
			Date: message.Date,
			Message: kz,
		}
		chatHistory = append(chatHistory, dateMessage)
	}
	chatHistory = append(chatHistory, message)
	return chatHistory
}

func (m *UserModel) GetChatID(firstPerson, secondPerson string) (int, error) {
	stmt := `SELECT id FROM Chats WHERE first_person_name = ? AND second_person_name = ?`

	var id int
	err := m.DB.QueryRow(stmt, firstPerson, secondPerson).Scan(&id)
	if err != nil {
		err := m.DB.QueryRow(stmt, secondPerson, firstPerson).Scan(&id)
		if err != nil {
			return 0, nil
		}
	}
	return id, nil
}

func (m *UserModel) UpdateChatHistory(firstName, secondName string, oldHistory []MessageStruct) error {
	stmt, err := m.DB.Prepare("UPDATE Chats SET history = ? WHERE (first_person_name = ? AND second_person_name = ?) OR (first_person_name = ? AND second_person_name = ?)")
	if err != nil {
		fmt.Println(" ----------------------- UpdateChatHistory 1 ------------------------")
		return err
	}
	defer stmt.Close()

	var historyBytes []byte
	newHistory := MarkMessagesAsSeen(oldHistory, firstName)
	historyBytes, err = SaveMessages(newHistory)
	if err != nil {
		fmt.Println(" ----------------------- UpdateChatHistory 2 ------------------------")
		return err
	}
	_, err = stmt.Exec(string(historyBytes), firstName, secondName, secondName, firstName)
	if err != nil {
		fmt.Println(" ----------------------- UpdateChatHistory 3 ------------------------" + err.Error())
		return err
	}

	return nil
}

func MarkMessagesAsSeen(messages []MessageStruct, author string) []MessageStruct {
	for i, msg := range messages {
		if msg.Author != author {
			// If the Author does not match, update Seen to true
			messages[i].Seen = true
		}
	}
	return messages
}

func (m *UserModel) GetChatHistory(id int) ([]MessageStruct, error) {
	stmt := `SELECT history FROM Chats WHERE id = ?`

	var history []MessageStruct
	var historyByte string
	err := m.DB.QueryRow(stmt, id).Scan(&historyByte)
	if err != nil {
		return nil, err
	}
	history, err = GetMessages(historyByte)
	if err != nil {
		return nil, err
	}

	return history, nil
}

func (m *UserModel) ChatExists(firstPerson, secondPerson string) (bool, error) {
	stmt := "SELECT COUNT(*) FROM Chats WHERE first_person_name = ? AND second_person_name = ?"
	var cnt int
	err := m.DB.QueryRow(stmt, firstPerson, secondPerson).Scan(&cnt)
	if err != nil {
		return false, err
	}
	if cnt == 0 {
		err := m.DB.QueryRow(stmt, secondPerson, firstPerson).Scan(&cnt)
		if err != nil {
			return false, err
		}
	}
	if cnt == 0 {
		return false, nil
	}
	return true, nil
}

func (m *UserModel) SaveStringArray(array []string) ([]byte, error) {
	jsonData, err := json.Marshal(array)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func GetStringArray(jsonData string) ([]string, error) {
	var stringArray []string
	err := json.Unmarshal([]byte(jsonData), &stringArray)
	if err != nil {
		return nil, err
	}

	return stringArray, nil
}

func (m *UserModel) GetHistoryOfChat(firstPerson, secondPerson string) ([]MessageStruct, error) {
	exist, err := m.ChatExists(firstPerson, secondPerson)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	chatId, err := m.GetChatID(firstPerson, secondPerson)
	if err != nil {
		return nil, err
	}
	history, err := m.GetChatHistory(chatId)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func SaveMessages(messages []MessageStruct) ([]byte, error) {
	messagesData, err := json.Marshal(messages)
	if err != nil {
		return nil, err
	}
	return messagesData, nil
}

func GetMessages(messagesData string) ([]MessageStruct, error) {
	var messages []MessageStruct
	err := json.Unmarshal([]byte(messagesData), &messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (u *UserModel) GetMyChats(username string) ([]StartedChat, error) {
	query := "SELECT second_person_name, history FROM Chats WHERE first_person_name = ?"
	rows, err := u.DB.Query(query, username)
	if err != nil {
		fmt.Println("ERROR 1 --------------------------------------------------------------")
		return nil, err
	}
	defer rows.Close()

	chats, err := u.processRows(rows)
	if err != nil {
		fmt.Println("ERROR 2 --------------------------------------------------------------")
		return nil, err
	}

	query2 := "SELECT first_person_name, history FROM Chats WHERE second_person_name = ?"
	rows2, err := u.DB.Query(query2, username)
	if err != nil {
		fmt.Println("ERROR 3 --------------------------------------------------------------")
		return nil, err
	}
	defer rows2.Close()

	chats2, err := u.processRows(rows2)
	if err != nil {
		fmt.Println("ERROR 4 --------------------------------------------------------------")
		return nil, err
	}

	chats = append(chats, chats2...)

	sorted, err := u.sortChatsByTime(chats)
	if err != nil {
		fmt.Println("ERROR 5 --------------------------------------------------------------")
		return nil, err
	}

	return sorted, nil
}

func (u *UserModel) sortChatsByTime(chats []StartedChat) ([]StartedChat, error) {
	// Custom sort function
	sort.Slice(chats, func(i, j int) bool {
		// Combine date and time strings
		dateTimeI := chats[i].LastMessageDate
		dateTimeJ := chats[j].LastMessageDate

		// Compare the combined strings
		return dateTimeI > dateTimeJ
	})

	return chats, nil
}

func (u *UserModel) processRows(rows *sql.Rows) ([]StartedChat, error) {
	var chats []StartedChat
	for rows.Next() {
		var chat StartedChat
		var byteMsg string
		var userName string
		var exist bool
		if err := rows.Scan(&userName, &byteMsg); err != nil {
			fmt.Println("processRows ERROR 1 --------------------------------------------------------------")
			// return nil, err
			exist = true
		}
		messagesArr, err := GetMessages(byteMsg)
		if err != nil {
			fmt.Println("processRows ERROR 2 --------------------------------------------------------------")
			// return nil, err
			exist = true
		}
		lastMessageStruct := GetLastMessage(messagesArr)
		userPhoto, err := u.GetPhotoByUserName(userName)
		if err != nil {
			//fmt.Println("processRows ERROR 3 --------------------------------------------------------------")
			// return
			exist = true
		}
		chat.With = userName
		chat.WithPhoto = userPhoto
		chat.WithStatus = u.GetUserStatus(chat.With) == 1
		// chat.LastMessage = ShortLastMessage(lastMessageStruct.Message)
		chat.LastMessage = lastMessageStruct.Message
		chat.LastMessageTime = lastMessageStruct.Time
		chat.LastMessageDate = lastMessageStruct.Date
		chat.UnseenMessagesCount = GetUnseenMessagesCount(userName, messagesArr)
		if !exist {
			chats = append(chats, chat)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

func (u *UserModel) GetAllUnseenMessagesCount(username string) int {
	var cnt int
	query := "SELECT history FROM Chats WHERE first_person_name = ? OR second_person_name = ?"
	rows, err := u.DB.Query(query, username, username)
	if err != nil {
		return 0
	}
	defer rows.Close()

	for rows.Next() {
		var history []MessageStruct
		var historyByte string
		if err := rows.Scan(&historyByte); err != nil {
			return 0
		}
		history, err = GetMessages(historyByte)
		if err != nil {
			return 0
		}
		localCnt := GetUnseenMessagesCountReverse(username, history)
		cnt += localCnt
	}
	return cnt
}

func GetUnseenMessagesCount(username string, messages []MessageStruct) int {
	var cnt int
	for _, msg := range messages {
		if msg.Author == username && !msg.Seen {
			cnt++
		}
	}
	return cnt
}

func GetUnseenMessagesCountReverse(username string, messages []MessageStruct) int {
	var cnt int
	for _, msg := range messages {
		if msg.Author != username && !msg.Seen {
			cnt++
		}
	}
	return cnt
}

func GetLastMessage(messages []MessageStruct) MessageStruct {
	if len(messages) == 0 {
		return MessageStruct{}
	}
	return messages[len(messages)-1]
}

func ShortLastMessage(message string) string {
	if len(message) > 30 {
		return message[:30] + "..."
	} else {
		return message
	}
}

func CompareDates(dateTimeStr1, dateTimeStr2 string) (bool, error) {
	layout := "02.01.2006 15:04"

	t1, err1 := time.Parse(layout, dateTimeStr1)
	if err1 != nil {
		return false, err1
	}

	t2, err2 := time.Parse(layout, dateTimeStr2)
	if err2 != nil {
		return false, err2
	}

	if t1.Year() != t2.Year() || t1.Month() != t2.Month() || t1.Day() != t2.Day() {
		return true, nil
	}

	return false, nil
}

func ConvertToKazakhCalendar(dateStr string) (string, error) {
	layout := "02.01.2006 15:04"
	
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return "", err
	}
	
	day := t.Day()
	month := kazakhMonths[t.Month()-1]
	
	formattedDate := fmt.Sprintf("%s %d", month, day)
	
	return formattedDate, nil
}

// SAVE FOR TIME

// func (m *UserModel) InsertChat(firstPerson, secondPerson string, message MessageStruct) error {
// 	maxRetries := 5
// 	var err error
// 	for i := 0; i < maxRetries; i++ {
// 		err = m.insertChat(firstPerson, secondPerson, message)
// 		if err != nil {
// 			if sqliteError, ok := err.(sqlite3.Error); ok && sqliteError.Code == sqlite3.ErrBusy {
// 				time.Sleep(time.Second) // Wait for 1 second before retrying
// 				continue
// 			}
// 		}
// 		break
// 	}
// 	return err
// }

// func (m *UserModel) insertChat(firstPerson, secondPerson string, message MessageStruct) error {
// 	check, err := m.ChatExists(firstPerson, secondPerson)
// 	if err != nil {
// 		fmt.Println("----------------------- InsertChat ERROR 1 --------------------------")
// 		return err
// 	}
// 	var chatHistory []MessageStruct
// 	var historyBytes []byte
// 	if check {
// 		chatId, err := m.GetChatID(firstPerson, secondPerson)
// 		if err != nil {
// 			fmt.Println("----------------------- InsertChat ERROR 2 --------------------------")
// 			return err
// 		}
// 		chatHistory, err = m.GetChatHistory(chatId)
// 		if err != nil {
// 			fmt.Println("----------------------- InsertChat ERROR 3 --------------------------")
// 			return err
// 		}
// 		chatHistory = append(chatHistory, message)
// 		updateStmt := "UPDATE Chats SET history = ? WHERE id = ?"
// 		historyBytes, err = SaveMessages(chatHistory)
// 		if err != nil {
// 			fmt.Println("----------------------- InsertChat ERROR 4 --------------------------")
// 			return err
// 		}
// 		_, err = m.DB.Exec(updateStmt, string(historyBytes), chatId)
// 		if err != nil {
// 			fmt.Println("----------------------- InsertChat ERROR 5 --------------------------")
// 			return err
// 		}
// 		return nil
// 	}
// 	chatHistory = append(chatHistory, message)
// 	historyBytes, err = SaveMessages(chatHistory)
// 	if err != nil {
// 		return err
// 	}
// 	stmt := `INSERT INTO Chats (first_person_name, second_person_name, history)
//              VALUES(?, ?, ?)`
// 	_, err = m.DB.Exec(stmt, firstPerson, secondPerson, string(historyBytes))
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// const maxRetries = 5
// const retryDelay = 100 * time.Millisecond // 100 milliseconds

// func (m *UserModel) UpdateChatHistory(firstName, secondName string, oldHistory []MessageStruct) error {
// 	var stmt *sql.Stmt
// 	var err error

// 	// Retry logic for "database is locked" error
// 	for i := 0; i < maxRetries; i++ {
// 		stmt, err = m.DB.Prepare("UPDATE Chats SET history = ? WHERE (first_person_name = ? AND second_person_name = ?) OR (first_person_name = ? AND second_person_name = ?)")
// 		if err == nil {
// 			break
// 		}
// 		if isDatabaseLockedError(err) {
// 			time.Sleep(retryDelay)
// 			continue
// 		}
// 		fmt.Println(" ----------------------- UpdateChatHistory 1 ------------------------")
// 		return err
// 	}
// 	if err != nil {
// 		return errors.New("UpdateChatHistory: failed to prepare statement after retries")
// 	}
// 	defer stmt.Close()

// 	var historyBytes []byte
// 	newHistory := MarkMessagesAsSeen(oldHistory, firstName)
// 	historyBytes, err = SaveMessages(newHistory)
// 	if err != nil {
// 		fmt.Println(" ----------------------- UpdateChatHistory 2 ------------------------")
// 		return err
// 	}

// 	for i := 0; i < maxRetries; i++ {
// 		_, err = stmt.Exec(string(historyBytes), firstName, secondName, secondName, firstName)
// 		if err == nil {
// 			break
// 		}
// 		if isDatabaseLockedError(err) {
// 			time.Sleep(retryDelay)
// 			continue
// 		}
// 		fmt.Println(" ----------------------- UpdateChatHistory 3 ------------------------" + err.Error())
// 		return err
// 	}
// 	if err != nil {
// 		return errors.New("UpdateChatHistory: failed to execute statement after retries")
// 	}

// 	return nil
// }

// Helper function to check if the error is a "database is locked" error
// func isDatabaseLockedError(err error) bool {
// 	return err.Error() == "database is locked"
// }

// func (m *UserModel) UpdateChatHistory(firstName, secondName string, oldHistory []MessageStruct) error {
// 	maxRetries := 5
// 	var err error

// 	for i := 0; i < maxRetries; i++ {
// 		mutex.Lock()
// 		err = m.updateChatHistory(firstName, secondName, oldHistory)
// 		mutex.Unlock()

// 		if err != nil {
// 			if sqliteError, ok := err.(sqlite3.Error); ok && sqliteError.Code == sqlite3.ErrBusy {
// 				time.Sleep(time.Second) // Wait for 1 second before retrying
// 				continue
// 			}
// 		}
// 		break
// 	}

// 	return err
// }

// func (m *UserModel) updateChatHistory(firstName, secondName string, oldHistory []MessageStruct) error {
// 	stmt, err := m.DB.Prepare("UPDATE Chats SET history = ? WHERE (first_person_name = ? AND second_person_name = ?) OR (first_person_name = ? AND second_person_name = ?)")
// 	if err != nil {
// 		fmt.Println(" ----------------------- UpdateChatHistory 1 ------------------------")
// 		return err
// 	}
// 	defer stmt.Close()

// 	var historyBytes []byte
// 	newHistory := MarkMessagesAsSeen(oldHistory, firstName)
// 	historyBytes, err = SaveMessages(newHistory)
// 	if err != nil {
// 		fmt.Println(" ----------------------- UpdateChatHistory 2 ------------------------")
// 		return err
// 	}

// 	_, err = stmt.Exec(string(historyBytes), firstName, secondName, secondName, firstName)
// 	if err != nil {
// 		fmt.Println(" ----------------------- UpdateChatHistory 3 ------------------------" + err.Error())
// 		return err
// 	}

// 	return nil
// }
