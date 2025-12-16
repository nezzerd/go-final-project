package repository

import (
	"context"
	"database/sql"

	"hotel-booking-system/internal/hotel/domain"
)

type PostgresRoomRepository struct {
	db *sql.DB
}

func NewPostgresRoomRepository(db *sql.DB) *PostgresRoomRepository {
	return &PostgresRoomRepository{db: db}
}

func (r *PostgresRoomRepository) CreateRoom(ctx context.Context, room *domain.Room) error {
	query := `INSERT INTO rooms (id, hotel_id, room_number, room_type, price_per_night, capacity, description, is_available) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
			  RETURNING created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		room.ID, room.HotelID, room.RoomNumber, room.RoomType,
		room.PricePerNight, room.Capacity, room.Description, room.IsAvailable,
	).Scan(&room.CreatedAt, &room.UpdatedAt)
}

func (r *PostgresRoomRepository) GetRoomByID(ctx context.Context, id string) (*domain.Room, error) {
	room := &domain.Room{}
	query := `SELECT id, hotel_id, room_number, room_type, price_per_night, capacity, 
			  description, is_available, created_at, updated_at 
			  FROM rooms WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&room.ID, &room.HotelID, &room.RoomNumber, &room.RoomType,
		&room.PricePerNight, &room.Capacity, &room.Description,
		&room.IsAvailable, &room.CreatedAt, &room.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (r *PostgresRoomRepository) GetRoomsByHotel(ctx context.Context, hotelID string) ([]domain.Room, error) {
	query := `SELECT id, hotel_id, room_number, room_type, price_per_night, capacity, 
			  description, is_available, created_at, updated_at 
			  FROM rooms WHERE hotel_id = $1 ORDER BY room_number`
	rows, err := r.db.QueryContext(ctx, query, hotelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []domain.Room
	for rows.Next() {
		var room domain.Room
		if err := rows.Scan(
			&room.ID, &room.HotelID, &room.RoomNumber, &room.RoomType,
			&room.PricePerNight, &room.Capacity, &room.Description,
			&room.IsAvailable, &room.CreatedAt, &room.UpdatedAt,
		); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, rows.Err()
}

func (r *PostgresRoomRepository) UpdateRoom(ctx context.Context, room *domain.Room) error {
	query := `UPDATE rooms SET room_number = $2, room_type = $3, price_per_night = $4, 
			  capacity = $5, description = $6, is_available = $7, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = $1 RETURNING updated_at`
	return r.db.QueryRowContext(ctx, query,
		room.ID, room.RoomNumber, room.RoomType, room.PricePerNight,
		room.Capacity, room.Description, room.IsAvailable,
	).Scan(&room.UpdatedAt)
}

func (r *PostgresRoomRepository) DeleteRoom(ctx context.Context, id string) error {
	query := `DELETE FROM rooms WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresRoomRepository) GetRoomPrice(ctx context.Context, hotelID, roomID string) (float64, error) {
	var price float64
	query := `SELECT price_per_night FROM rooms WHERE id = $1 AND hotel_id = $2`
	err := r.db.QueryRowContext(ctx, query, roomID, hotelID).Scan(&price)
	return price, err
}
