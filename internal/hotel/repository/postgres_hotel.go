package repository

import (
	"context"
	"database/sql"

	"hotel-booking-system/internal/hotel/domain"
)

type PostgresHotelRepository struct {
	db *sql.DB
}

func NewPostgresHotelRepository(db *sql.DB) *PostgresHotelRepository {
	return &PostgresHotelRepository{db: db}
}

func (r *PostgresHotelRepository) CreateHotel(ctx context.Context, hotel *domain.Hotel) error {
	query := `INSERT INTO hotels (id, name, description, address, owner_id) 
			  VALUES ($1, $2, $3, $4, $5) 
			  RETURNING created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		hotel.ID, hotel.Name, hotel.Description, hotel.Address, hotel.OwnerID,
	).Scan(&hotel.CreatedAt, &hotel.UpdatedAt)
}

func (r *PostgresHotelRepository) GetHotelByID(ctx context.Context, id string) (*domain.Hotel, error) {
	hotel := &domain.Hotel{}
	query := `SELECT id, name, description, address, owner_id, created_at, updated_at 
			  FROM hotels WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&hotel.ID, &hotel.Name, &hotel.Description, &hotel.Address,
		&hotel.OwnerID, &hotel.CreatedAt, &hotel.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return hotel, nil
}

func (r *PostgresHotelRepository) GetHotels(ctx context.Context, limit, offset int) ([]domain.Hotel, error) {
	query := `SELECT id, name, description, address, owner_id, created_at, updated_at 
			  FROM hotels ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hotels []domain.Hotel
	for rows.Next() {
		var hotel domain.Hotel
		if err := rows.Scan(
			&hotel.ID, &hotel.Name, &hotel.Description, &hotel.Address,
			&hotel.OwnerID, &hotel.CreatedAt, &hotel.UpdatedAt,
		); err != nil {
			return nil, err
		}
		hotels = append(hotels, hotel)
	}
	return hotels, rows.Err()
}

func (r *PostgresHotelRepository) GetHotelsByOwner(ctx context.Context, ownerID string) ([]domain.Hotel, error) {
	query := `SELECT id, name, description, address, owner_id, created_at, updated_at 
			  FROM hotels WHERE owner_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hotels []domain.Hotel
	for rows.Next() {
		var hotel domain.Hotel
		if err := rows.Scan(
			&hotel.ID, &hotel.Name, &hotel.Description, &hotel.Address,
			&hotel.OwnerID, &hotel.CreatedAt, &hotel.UpdatedAt,
		); err != nil {
			return nil, err
		}
		hotels = append(hotels, hotel)
	}
	return hotels, rows.Err()
}

func (r *PostgresHotelRepository) UpdateHotel(ctx context.Context, hotel *domain.Hotel) error {
	query := `UPDATE hotels SET name = $2, description = $3, address = $4, 
			  updated_at = CURRENT_TIMESTAMP WHERE id = $1
			  RETURNING updated_at`
	return r.db.QueryRowContext(ctx, query,
		hotel.ID, hotel.Name, hotel.Description, hotel.Address,
	).Scan(&hotel.UpdatedAt)
}

func (r *PostgresHotelRepository) DeleteHotel(ctx context.Context, id string) error {
	query := `DELETE FROM hotels WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
