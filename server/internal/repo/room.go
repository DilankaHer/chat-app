package repo

import "database/sql"

type JoinRoomRepository interface {
	JoinRoom(room *UserRoom) error
	GetRoomIds() ([]string, error)
	DeleteRoomUsersByUserId(userId string) error
	GetRooms() ([]Room, error)
}

type JoinRoomRepo struct {
	db *sql.DB
}

type UserRoom struct {
	Id     string `json:"id"`
	RoomId string `json:"roomId" validate:"required"`
	UserId string `json:"userId" validate:"required"`
}

type Room struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewJoinRoomRepo(db *sql.DB) JoinRoomRepository {
	return &JoinRoomRepo{db: db}
}

func (rr *JoinRoomRepo) JoinRoom(room *UserRoom) error {
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

func (rr *JoinRoomRepo) GetRoomIds() ([]string, error) {
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

func (rr *JoinRoomRepo) DeleteRoomUsersByUserId(userId string) error {
	query := `DELETE FROM rooms_users WHERE user_id = $1`

	_, err := rr.db.Exec(query, userId)
	if err != nil {
		return err
	}

	return nil
}

func (rr *JoinRoomRepo) GetRooms() ([]Room, error) {
	query := `SELECT id, name FROM rooms`

	var rooms []Room
	rows, err := rr.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var room Room
		err := rows.Scan(&room.Id, &room.Name)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}
