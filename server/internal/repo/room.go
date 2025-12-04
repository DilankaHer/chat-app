package repo

import (
	"database/sql"
	"duhchat/internal/api/model"
)

type RoomRepository interface {
	JoinRoom(room *model.UserRoom) error
	CreateRoom(room *Room) error
	GetRoomIds() ([]string, error)
	DeleteRoomUsersByUserId(userId string) error
	GetRooms() ([]Room, error)
}
type RoomRepo struct {
	db *sql.DB
}

type Room struct {
	RoomId string `json:"roomId"`
	Name   string `json:"name"`
}

func NewRoomRepo(db *sql.DB) RoomRepository {
	return &RoomRepo{db: db}
}

func (rr *RoomRepo) JoinRoom(room *model.UserRoom) error {
	tx, err := rr.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO rooms_users (room_id, user_id) VALUES ($1, $2) RETURNING id`

	err = tx.QueryRow(query, room.RoomId, room.UserId).Scan(&room.Id)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (rr *RoomRepo) CreateRoom(room *Room) error {
	tx, err := rr.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO rooms (name) VALUES ($1) RETURNING id`

	err = tx.QueryRow(query, room.Name).Scan(&room.RoomId)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (rr *RoomRepo) GetRoomIds() ([]string, error) {
	query := `SELECT id FROM rooms`

	var ids []string
	rows, err := rr.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (rr *RoomRepo) DeleteRoomUsersByUserId(userId string) error {
	query := `DELETE FROM rooms_users WHERE user_id = $1`

	_, err := rr.db.Exec(query, userId)
	if err != nil {
		return err
	}

	return nil
}

func (rr *RoomRepo) GetRooms() ([]Room, error) {
	query := `SELECT id, name FROM rooms`

	var rooms []Room
	rows, err := rr.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var room Room
		err := rows.Scan(&room.RoomId, &room.Name)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}
