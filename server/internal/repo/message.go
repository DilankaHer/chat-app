package repo

import "database/sql"

type MessageRepository interface {
	AddMessasge(message *Message) error
	GetMessagesByRoomId(roomId string) ([]Message, error)
}

type MessageRepo struct {
	db *sql.DB
}
type Message struct {
	MessageId string `json:"messageId"`
	RoomId    string `json:"roomId" validate:"required"`
	UserId    string `json:"userId" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Content   string `json:"content" validate:"required"`
	IsError   bool   `json:"isError"`
}

func NewMessageRepo(db *sql.DB) MessageRepository {
	return &MessageRepo{db: db}
}

func (mr *MessageRepo) AddMessasge(message *Message) error {
	tx, err := mr.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO messages (room_id, user_id, content) VALUES ($1, $2, $3) RETURNING id`

	err = tx.QueryRow(query, message.RoomId, message.UserId, message.Content).Scan(&message.MessageId)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (mr *MessageRepo) GetMessagesByRoomId(roomId string) ([]Message, error) {
	query := `SELECT m.id, m.room_id, m.user_id, m.content, u.username FROM messages AS m
	          INNER JOIN users AS u ON u.id = m.user_id WHERE m.room_id = $1`

	rows, err := mr.db.Query(query, roomId)
	if err != nil {
		return nil, err
	}

	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.MessageId, &message.RoomId, &message.UserId, &message.Content, &message.Username)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
