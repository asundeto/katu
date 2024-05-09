package models

import (
	"database/sql"
	"encoding/json"
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
	UnseenMessagesCount int
}

type MessageSolo struct {
	ID      int
	Time    string
	Message string
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

func (m *UserModel) InsertChat(firstPerson, secondPerson string, message MessageSolo) error {
	check, err := m.ChatExists(firstPerson, secondPerson)
	if err != nil {
		return err
	}
	var chatHistory []MessageSolo
	var historyBytes []byte
	if check {
		chatId, err := m.GetChatID(firstPerson, secondPerson)
		if err != nil {
			return err
		}
		chatHistory, err = m.GetChatHistory(chatId)
		if err != nil {
			return err
		}
		chatHistory = append(chatHistory, message)
		updateStmt := "UPDATE Chats SET history = ? WHERE id = ?"
		historyBytes, err = SaveMessages(chatHistory)
		if err != nil {
			return err
		}
		_, err = m.DB.Exec(updateStmt, string(historyBytes), chatId)
		if err != nil {
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

func (m *UserModel) UpdateChatHistory(firstName, secondName string, oldHistory []MessageSolo) error {
	stmt, err := m.DB.Prepare("UPDATE Chats SET history = ? WHERE (first_person_name = ? AND second_person_name = ?) OR (first_person_name = ? AND second_person_name = ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var historyBytes []byte
	newHistory := MarkMessagesAsSeen(oldHistory, firstName)
	historyBytes, err = SaveMessages(newHistory)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(string(historyBytes), firstName, secondName, secondName, firstName)
	if err != nil {
		return err
	}

	return nil
}

func MarkMessagesAsSeen(messages []MessageSolo, author string) []MessageSolo {
	for i, msg := range messages {
		if msg.Author != author {
			// If the Author does not match, update Seen to true
			messages[i].Seen = true
		}
	}
	return messages
}

func (m *UserModel) GetChatHistory(id int) ([]MessageSolo, error) {
	stmt := `SELECT history FROM Chats WHERE id = ?`

	var history []MessageSolo
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

func (m *UserModel) GetHistoryOfChat(firstPerson, secondPerson string) ([]MessageSolo, error) {
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

func SaveMessages(messages []MessageSolo) ([]byte, error) {
	messagesData, err := json.Marshal(messages)
	if err != nil {
		return nil, err
	}
	return messagesData, nil
}

func GetMessages(messagesData string) ([]MessageSolo, error) {
	var messages []MessageSolo
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
		return nil, err
	}
	defer rows.Close()

	chats, err := u.processRows(rows)
	if err != nil {
		return nil, err
	}

	query2 := "SELECT first_person_name, history FROM Chats WHERE second_person_name = ?"
	rows2, err := u.DB.Query(query2, username)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	chats2, err := u.processRows(rows2)
	if err != nil {
		return nil, err
	}

	chats = append(chats, chats2...)

	return chats, nil
}

func (u *UserModel) processRows(rows *sql.Rows) ([]StartedChat, error) {
	var chats []StartedChat
	for rows.Next() {
		var chat StartedChat
		var byteMsg string
		var userName string
		if err := rows.Scan(&userName, &byteMsg); err != nil {
			return nil, err
		}
		messagesArr, err := GetMessages(byteMsg)
		if err != nil {
			return nil, err
		}
		lastMessageStruct := GetLastMessage(messagesArr)
		userPhoto, err := u.GetPhotoByUserName(userName)
		if err != nil {
			return nil, err
		}
		chat.With = userName
		chat.WithPhoto = userPhoto
		chat.WithStatus = u.GetUserStatus(chat.With) == 1
		chat.LastMessage = ShortLastMessage(lastMessageStruct.Message)
		chat.LastMessageTime = lastMessageStruct.Time
		chat.UnseenMessagesCount = GetUnseenMessagesCount(userName, messagesArr)
		chats = append(chats, chat)
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
		var history []MessageSolo
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

func GetUnseenMessagesCount(username string, messages []MessageSolo) int {
	var cnt int
	for _, msg := range messages {
		if msg.Author == username && !msg.Seen {
			cnt++
		}
	}
	return cnt
}

func GetUnseenMessagesCountReverse(username string, messages []MessageSolo) int {
	var cnt int
	for _, msg := range messages {
		if msg.Author != username && !msg.Seen {
			cnt++
		}
	}
	return cnt
}

func GetLastMessage(messages []MessageSolo) MessageSolo {
	if len(messages) == 0 {
		return MessageSolo{}
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
