package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"hotel-booking-system/internal/booking/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	return db, mock
}

func TestNewPostgresBookingRepository(t *testing.T) {
	db, _ := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestCreateBooking_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	booking := &domain.Booking{
		ID:            "booking-123",
		UserID:        "user-123",
		HotelID:       "hotel-123",
		RoomID:        "room-123",
		CheckInDate:   time.Now(),
		CheckOutDate:  time.Now().Add(24 * time.Hour),
		TotalPrice:    5000.0,
		Status:        "pending",
		PaymentStatus: "pending",
	}

	createdAt := time.Now()
	updatedAt := time.Now()

	mock.ExpectQuery(`INSERT INTO bookings`).
		WithArgs(
			booking.ID, booking.UserID, booking.HotelID, booking.RoomID,
			booking.CheckInDate, booking.CheckOutDate, booking.TotalPrice,
			booking.Status, booking.PaymentStatus,
		).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
			AddRow(createdAt, updatedAt))

	err := repo.CreateBooking(context.Background(), booking)
	assert.NoError(t, err)
	assert.Equal(t, createdAt, booking.CreatedAt)
	assert.Equal(t, updatedAt, booking.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateBooking_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	booking := &domain.Booking{
		ID:            "booking-123",
		UserID:        "user-123",
		HotelID:       "hotel-123",
		RoomID:        "room-123",
		CheckInDate:   time.Now(),
		CheckOutDate:  time.Now().Add(24 * time.Hour),
		TotalPrice:    5000.0,
		Status:        "pending",
		PaymentStatus: "pending",
	}

	mock.ExpectQuery(`INSERT INTO bookings`).
		WithArgs(
			booking.ID, booking.UserID, booking.HotelID, booking.RoomID,
			booking.CheckInDate, booking.CheckOutDate, booking.TotalPrice,
			booking.Status, booking.PaymentStatus,
		).
		WillReturnError(errors.New("database error"))

	err := repo.CreateBooking(context.Background(), booking)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookingByID_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	bookingID := "booking-123"
	createdAt := time.Now()
	updatedAt := time.Now()
	checkIn := time.Now()
	checkOut := time.Now().Add(24 * time.Hour)

	mock.ExpectQuery(`SELECT.*FROM bookings WHERE id`).
		WithArgs(bookingID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "hotel_id", "room_id", "check_in_date", "check_out_date",
			"total_price", "status", "payment_status", "created_at", "updated_at",
		}).AddRow(
			bookingID, "user-123", "hotel-123", "room-123",
			checkIn, checkOut, 5000.0, "pending", "pending",
			createdAt, updatedAt,
		))

	booking, err := repo.GetBookingByID(context.Background(), bookingID)
	assert.NoError(t, err)
	assert.NotNil(t, booking)
	assert.Equal(t, bookingID, booking.ID)
	assert.Equal(t, "user-123", booking.UserID)
	assert.Equal(t, 5000.0, booking.TotalPrice)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookingByID_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	bookingID := "non-existent"

	mock.ExpectQuery(`SELECT.*FROM bookings WHERE id`).
		WithArgs(bookingID).
		WillReturnError(sql.ErrNoRows)

	booking, err := repo.GetBookingByID(context.Background(), bookingID)
	assert.Error(t, err)
	assert.Nil(t, booking)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookingByID_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	bookingID := "booking-123"

	mock.ExpectQuery(`SELECT.*FROM bookings WHERE id`).
		WithArgs(bookingID).
		WillReturnError(errors.New("connection error"))

	booking, err := repo.GetBookingByID(context.Background(), bookingID)
	assert.Error(t, err)
	assert.Nil(t, booking)
	assert.Contains(t, err.Error(), "connection error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookingsByUser_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	userID := "user-123"
	createdAt := time.Now()
	updatedAt := time.Now()
	checkIn := time.Now()
	checkOut := time.Now().Add(24 * time.Hour)

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "hotel_id", "room_id", "check_in_date", "check_out_date",
		"total_price", "status", "payment_status", "created_at", "updated_at",
	}).
		AddRow("booking-1", userID, "hotel-1", "room-1", checkIn, checkOut, 5000.0, "pending", "pending", createdAt, updatedAt).
		AddRow("booking-2", userID, "hotel-2", "room-2", checkIn, checkOut, 6000.0, "confirmed", "paid", createdAt, updatedAt)

	mock.ExpectQuery(`SELECT.*FROM bookings WHERE user_id`).
		WithArgs(userID).
		WillReturnRows(rows)

	bookings, err := repo.GetBookingsByUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, bookings, 2)
	assert.Equal(t, "booking-1", bookings[0].ID)
	assert.Equal(t, "booking-2", bookings[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookingsByUser_EmptyResult(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	userID := "user-no-bookings"

	mock.ExpectQuery(`SELECT.*FROM bookings WHERE user_id`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "hotel_id", "room_id", "check_in_date", "check_out_date",
			"total_price", "status", "payment_status", "created_at", "updated_at",
		}))

	bookings, err := repo.GetBookingsByUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Empty(t, bookings)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookingsByUser_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	userID := "user-123"

	mock.ExpectQuery(`SELECT.*FROM bookings WHERE user_id`).
		WithArgs(userID).
		WillReturnError(errors.New("query error"))

	bookings, err := repo.GetBookingsByUser(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, bookings)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookingsByUser_ScanError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	userID := "user-123"

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "hotel_id", "room_id", "check_in_date", "check_out_date",
		"total_price", "status", "payment_status", "created_at", "updated_at",
	}).AddRow("invalid", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	mock.ExpectQuery(`SELECT.*FROM bookings WHERE user_id`).
		WithArgs(userID).
		WillReturnRows(rows)

	bookings, err := repo.GetBookingsByUser(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, bookings)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookingsByHotel_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	hotelID := "hotel-123"
	createdAt := time.Now()
	updatedAt := time.Now()
	checkIn := time.Now()
	checkOut := time.Now().Add(24 * time.Hour)

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "hotel_id", "room_id", "check_in_date", "check_out_date",
		"total_price", "status", "payment_status", "created_at", "updated_at",
	}).
		AddRow("booking-1", "user-1", hotelID, "room-1", checkIn, checkOut, 5000.0, "pending", "pending", createdAt, updatedAt)

	mock.ExpectQuery(`SELECT.*FROM bookings WHERE hotel_id`).
		WithArgs(hotelID).
		WillReturnRows(rows)

	bookings, err := repo.GetBookingsByHotel(context.Background(), hotelID)
	assert.NoError(t, err)
	assert.Len(t, bookings, 1)
	assert.Equal(t, hotelID, bookings[0].HotelID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookingsByHotel_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	hotelID := "hotel-123"

	mock.ExpectQuery(`SELECT.*FROM bookings WHERE hotel_id`).
		WithArgs(hotelID).
		WillReturnError(errors.New("query error"))

	bookings, err := repo.GetBookingsByHotel(context.Background(), hotelID)
	assert.Error(t, err)
	assert.Nil(t, bookings)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookingsByHotel_RowsError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	hotelID := "hotel-123"
	createdAt := time.Now()
	updatedAt := time.Now()
	checkIn := time.Now()
	checkOut := time.Now().Add(24 * time.Hour)

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "hotel_id", "room_id", "check_in_date", "check_out_date",
		"total_price", "status", "payment_status", "created_at", "updated_at",
	}).AddRow("booking-1", "user-1", hotelID, "room-1", checkIn, checkOut, 5000.0, "pending", "pending", createdAt, updatedAt).
		RowError(0, errors.New("row error"))

	mock.ExpectQuery(`SELECT.*FROM bookings WHERE hotel_id`).
		WithArgs(hotelID).
		WillReturnRows(rows)

	bookings, err := repo.GetBookingsByHotel(context.Background(), hotelID)
	assert.Error(t, err)
	assert.Nil(t, bookings)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBookingStatus_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	bookingID := "booking-123"
	status := "confirmed"

	mock.ExpectExec(`UPDATE bookings SET status`).
		WithArgs(bookingID, status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdateBookingStatus(context.Background(), bookingID, status)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBookingStatus_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	bookingID := "non-existent"
	status := "confirmed"

	mock.ExpectExec(`UPDATE bookings SET status`).
		WithArgs(bookingID, status).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.UpdateBookingStatus(context.Background(), bookingID, status)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBookingStatus_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	bookingID := "booking-123"
	status := "confirmed"

	mock.ExpectExec(`UPDATE bookings SET status`).
		WithArgs(bookingID, status).
		WillReturnError(errors.New("update error"))

	err := repo.UpdateBookingStatus(context.Background(), bookingID, status)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePaymentStatus_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	bookingID := "booking-123"
	paymentStatus := "paid"

	mock.ExpectExec(`UPDATE bookings SET payment_status`).
		WithArgs(bookingID, paymentStatus).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdatePaymentStatus(context.Background(), bookingID, paymentStatus)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePaymentStatus_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	bookingID := "booking-123"
	paymentStatus := "paid"

	mock.ExpectExec(`UPDATE bookings SET payment_status`).
		WithArgs(bookingID, paymentStatus).
		WillReturnError(errors.New("update error"))

	err := repo.UpdatePaymentStatus(context.Background(), bookingID, paymentStatus)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePaymentStatus_AllStatuses(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresBookingRepository(db)
	bookingID := "booking-123"
	statuses := []string{"pending", "paid", "failed", "refunded"}

	for _, status := range statuses {
		mock.ExpectExec(`UPDATE bookings SET payment_status`).
			WithArgs(bookingID, status).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdatePaymentStatus(context.Background(), bookingID, status)
		assert.NoError(t, err)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}
