package repository

import (
	"context"
	"database/sql"

	"hotel-booking-system/internal/booking/domain"
)

type PostgresBookingRepository struct {
	db *sql.DB
}

func NewPostgresBookingRepository(db *sql.DB) *PostgresBookingRepository {
	return &PostgresBookingRepository{db: db}
}

func (r *PostgresBookingRepository) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	query := `INSERT INTO bookings (id, user_id, hotel_id, room_id, check_in_date, check_out_date, 
			  total_price, status, payment_status) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
			  RETURNING created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		booking.ID, booking.UserID, booking.HotelID, booking.RoomID,
		booking.CheckInDate, booking.CheckOutDate, booking.TotalPrice,
		booking.Status, booking.PaymentStatus,
	).Scan(&booking.CreatedAt, &booking.UpdatedAt)
}

func (r *PostgresBookingRepository) GetBookingByID(ctx context.Context, id string) (*domain.Booking, error) {
	booking := &domain.Booking{}
	query := `SELECT id, user_id, hotel_id, room_id, check_in_date, check_out_date, 
			  total_price, status, payment_status, created_at, updated_at 
			  FROM bookings WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&booking.ID, &booking.UserID, &booking.HotelID, &booking.RoomID,
		&booking.CheckInDate, &booking.CheckOutDate, &booking.TotalPrice,
		&booking.Status, &booking.PaymentStatus, &booking.CreatedAt, &booking.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return booking, nil
}

func (r *PostgresBookingRepository) GetBookingsByUser(ctx context.Context, userID string) ([]domain.Booking, error) {
	query := `SELECT id, user_id, hotel_id, room_id, check_in_date, check_out_date, 
			  total_price, status, payment_status, created_at, updated_at 
			  FROM bookings WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		var booking domain.Booking
		if err := rows.Scan(
			&booking.ID, &booking.UserID, &booking.HotelID, &booking.RoomID,
			&booking.CheckInDate, &booking.CheckOutDate, &booking.TotalPrice,
			&booking.Status, &booking.PaymentStatus, &booking.CreatedAt, &booking.UpdatedAt,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	return bookings, rows.Err()
}

func (r *PostgresBookingRepository) GetBookingsByHotel(ctx context.Context, hotelID string) ([]domain.Booking, error) {
	query := `SELECT id, user_id, hotel_id, room_id, check_in_date, check_out_date, 
			  total_price, status, payment_status, created_at, updated_at 
			  FROM bookings WHERE hotel_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, hotelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		var booking domain.Booking
		if err := rows.Scan(
			&booking.ID, &booking.UserID, &booking.HotelID, &booking.RoomID,
			&booking.CheckInDate, &booking.CheckOutDate, &booking.TotalPrice,
			&booking.Status, &booking.PaymentStatus, &booking.CreatedAt, &booking.UpdatedAt,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	return bookings, rows.Err()
}

func (r *PostgresBookingRepository) UpdateBookingStatus(ctx context.Context, id, status string) error {
	query := `UPDATE bookings SET status = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, status)
	return err
}

func (r *PostgresBookingRepository) UpdatePaymentStatus(ctx context.Context, id, paymentStatus string) error {
	query := `UPDATE bookings SET payment_status = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, paymentStatus)
	return err
}
