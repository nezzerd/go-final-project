package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"hotel-booking-system/internal/hotel/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupMockDBForRoom(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	return db, mock
}

func TestNewPostgresRoomRepository(t *testing.T) {
	db, _ := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestCreateRoom_Success(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	room := &domain.Room{
		ID:            "room-123",
		HotelID:       "hotel-123",
		RoomNumber:    "101",
		RoomType:      "Standard",
		PricePerNight: 5000.0,
		Capacity:      2,
		Description:   "Comfortable room",
		IsAvailable:   true,
	}

	createdAt := time.Now()
	updatedAt := time.Now()

	mock.ExpectQuery(`INSERT INTO rooms`).
		WithArgs(
			room.ID, room.HotelID, room.RoomNumber, room.RoomType,
			room.PricePerNight, room.Capacity, room.Description, room.IsAvailable,
		).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
			AddRow(createdAt, updatedAt))

	err := repo.CreateRoom(context.Background(), room)
	assert.NoError(t, err)
	assert.Equal(t, createdAt, room.CreatedAt)
	assert.Equal(t, updatedAt, room.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateRoom_DatabaseError(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	room := &domain.Room{
		ID:            "room-123",
		HotelID:       "hotel-123",
		RoomNumber:    "101",
		RoomType:      "Standard",
		PricePerNight: 5000.0,
		Capacity:      2,
		Description:   "Comfortable room",
		IsAvailable:   true,
	}

	mock.ExpectQuery(`INSERT INTO rooms`).
		WithArgs(
			room.ID, room.HotelID, room.RoomNumber, room.RoomType,
			room.PricePerNight, room.Capacity, room.Description, room.IsAvailable,
		).
		WillReturnError(errors.New("duplicate key"))

	err := repo.CreateRoom(context.Background(), room)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRoomByID_Success(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	roomID := "room-123"
	createdAt := time.Now()
	updatedAt := time.Now()

	mock.ExpectQuery(`SELECT.*FROM rooms WHERE id`).
		WithArgs(roomID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "hotel_id", "room_number", "room_type", "price_per_night",
			"capacity", "description", "is_available", "created_at", "updated_at",
		}).AddRow(
			roomID, "hotel-123", "101", "Standard", 5000.0,
			2, "Comfortable room", true, createdAt, updatedAt,
		))

	room, err := repo.GetRoomByID(context.Background(), roomID)
	assert.NoError(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, roomID, room.ID)
	assert.Equal(t, "101", room.RoomNumber)
	assert.Equal(t, 5000.0, room.PricePerNight)
	assert.True(t, room.IsAvailable)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRoomByID_NotFound(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	roomID := "non-existent"

	mock.ExpectQuery(`SELECT.*FROM rooms WHERE id`).
		WithArgs(roomID).
		WillReturnError(sql.ErrNoRows)

	room, err := repo.GetRoomByID(context.Background(), roomID)
	assert.Error(t, err)
	assert.Nil(t, room)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRoomByID_DatabaseError(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	roomID := "room-123"

	mock.ExpectQuery(`SELECT.*FROM rooms WHERE id`).
		WithArgs(roomID).
		WillReturnError(errors.New("connection error"))

	room, err := repo.GetRoomByID(context.Background(), roomID)
	assert.Error(t, err)
	assert.Nil(t, room)
	assert.Contains(t, err.Error(), "connection error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRoomsByHotel_Success(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	hotelID := "hotel-123"
	createdAt := time.Now()
	updatedAt := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "hotel_id", "room_number", "room_type", "price_per_night",
		"capacity", "description", "is_available", "created_at", "updated_at",
	}).
		AddRow("room-1", hotelID, "101", "Standard", 5000.0, 2, "Room 1", true, createdAt, updatedAt).
		AddRow("room-2", hotelID, "102", "Deluxe", 8000.0, 3, "Room 2", true, createdAt, updatedAt)

	mock.ExpectQuery(`SELECT.*FROM rooms WHERE hotel_id`).
		WithArgs(hotelID).
		WillReturnRows(rows)

	rooms, err := repo.GetRoomsByHotel(context.Background(), hotelID)
	assert.NoError(t, err)
	assert.Len(t, rooms, 2)
	assert.Equal(t, hotelID, rooms[0].HotelID)
	assert.Equal(t, hotelID, rooms[1].HotelID)
	assert.Equal(t, "101", rooms[0].RoomNumber)
	assert.Equal(t, "102", rooms[1].RoomNumber)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRoomsByHotel_EmptyResult(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	hotelID := "hotel-no-rooms"

	mock.ExpectQuery(`SELECT.*FROM rooms WHERE hotel_id`).
		WithArgs(hotelID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "hotel_id", "room_number", "room_type", "price_per_night",
			"capacity", "description", "is_available", "created_at", "updated_at",
		}))

	rooms, err := repo.GetRoomsByHotel(context.Background(), hotelID)
	assert.NoError(t, err)
	assert.Empty(t, rooms)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRoomsByHotel_DatabaseError(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	hotelID := "hotel-123"

	mock.ExpectQuery(`SELECT.*FROM rooms WHERE hotel_id`).
		WithArgs(hotelID).
		WillReturnError(errors.New("query error"))

	rooms, err := repo.GetRoomsByHotel(context.Background(), hotelID)
	assert.Error(t, err)
	assert.Nil(t, rooms)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRoomsByHotel_ScanError(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	hotelID := "hotel-123"

	rows := sqlmock.NewRows([]string{
		"id", "hotel_id", "room_number", "room_type", "price_per_night",
		"capacity", "description", "is_available", "created_at", "updated_at",
	}).AddRow(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	mock.ExpectQuery(`SELECT.*FROM rooms WHERE hotel_id`).
		WithArgs(hotelID).
		WillReturnRows(rows)

	rooms, err := repo.GetRoomsByHotel(context.Background(), hotelID)
	assert.Error(t, err)
	assert.Nil(t, rooms)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRoomsByHotel_RowsError(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	hotelID := "hotel-123"
	createdAt := time.Now()
	updatedAt := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "hotel_id", "room_number", "room_type", "price_per_night",
		"capacity", "description", "is_available", "created_at", "updated_at",
	}).AddRow("room-1", hotelID, "101", "Standard", 5000.0, 2, "Room 1", true, createdAt, updatedAt).
		RowError(0, errors.New("row error"))

	mock.ExpectQuery(`SELECT.*FROM rooms WHERE hotel_id`).
		WithArgs(hotelID).
		WillReturnRows(rows)

	rooms, err := repo.GetRoomsByHotel(context.Background(), hotelID)
	assert.Error(t, err)
	assert.Nil(t, rooms)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateRoom_Success(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	room := &domain.Room{
		ID:            "room-123",
		RoomNumber:    "101",
		RoomType:      "Deluxe",
		PricePerNight: 8000.0,
		Capacity:      3,
		Description:   "Updated room",
		IsAvailable:   false,
	}
	updatedAt := time.Now()

	mock.ExpectQuery(`UPDATE rooms SET`).
		WithArgs(
			room.ID, room.RoomNumber, room.RoomType, room.PricePerNight,
			room.Capacity, room.Description, room.IsAvailable,
		).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).AddRow(updatedAt))

	err := repo.UpdateRoom(context.Background(), room)
	assert.NoError(t, err)
	assert.Equal(t, updatedAt, room.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateRoom_NotFound(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	room := &domain.Room{
		ID:            "non-existent",
		RoomNumber:    "101",
		RoomType:      "Standard",
		PricePerNight: 5000.0,
		Capacity:      2,
		Description:   "Room",
		IsAvailable:   true,
	}

	mock.ExpectQuery(`UPDATE rooms SET`).
		WithArgs(
			room.ID, room.RoomNumber, room.RoomType, room.PricePerNight,
			room.Capacity, room.Description, room.IsAvailable,
		).
		WillReturnError(sql.ErrNoRows)

	err := repo.UpdateRoom(context.Background(), room)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateRoom_DatabaseError(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	room := &domain.Room{
		ID:            "room-123",
		RoomNumber:    "101",
		RoomType:      "Standard",
		PricePerNight: 5000.0,
		Capacity:      2,
		Description:   "Room",
		IsAvailable:   true,
	}

	mock.ExpectQuery(`UPDATE rooms SET`).
		WithArgs(
			room.ID, room.RoomNumber, room.RoomType, room.PricePerNight,
			room.Capacity, room.Description, room.IsAvailable,
		).
		WillReturnError(errors.New("update error"))

	err := repo.UpdateRoom(context.Background(), room)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteRoom_Success(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	roomID := "room-123"

	mock.ExpectExec(`DELETE FROM rooms WHERE id`).
		WithArgs(roomID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteRoom(context.Background(), roomID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteRoom_NotFound(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	roomID := "non-existent"

	mock.ExpectExec(`DELETE FROM rooms WHERE id`).
		WithArgs(roomID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeleteRoom(context.Background(), roomID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteRoom_DatabaseError(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	roomID := "room-123"

	mock.ExpectExec(`DELETE FROM rooms WHERE id`).
		WithArgs(roomID).
		WillReturnError(errors.New("delete error"))

	err := repo.DeleteRoom(context.Background(), roomID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delete error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRoomPrice_Success(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	hotelID := "hotel-123"
	roomID := "room-123"
	expectedPrice := 5000.0

	mock.ExpectQuery(`SELECT price_per_night FROM rooms WHERE id`).
		WithArgs(roomID, hotelID).
		WillReturnRows(sqlmock.NewRows([]string{"price_per_night"}).AddRow(expectedPrice))

	price, err := repo.GetRoomPrice(context.Background(), hotelID, roomID)
	assert.NoError(t, err)
	assert.Equal(t, expectedPrice, price)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRoomPrice_NotFound(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	hotelID := "hotel-123"
	roomID := "non-existent"

	mock.ExpectQuery(`SELECT price_per_night FROM rooms WHERE id`).
		WithArgs(roomID, hotelID).
		WillReturnError(sql.ErrNoRows)

	price, err := repo.GetRoomPrice(context.Background(), hotelID, roomID)
	assert.Error(t, err)
	assert.Equal(t, 0.0, price)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRoomPrice_DatabaseError(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	hotelID := "hotel-123"
	roomID := "room-123"

	mock.ExpectQuery(`SELECT price_per_night FROM rooms WHERE id`).
		WithArgs(roomID, hotelID).
		WillReturnError(errors.New("query error"))

	price, err := repo.GetRoomPrice(context.Background(), hotelID, roomID)
	assert.Error(t, err)
	assert.Equal(t, 0.0, price)
	assert.Contains(t, err.Error(), "query error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRoomPrice_WrongHotel(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	hotelID := "wrong-hotel"
	roomID := "room-123"

	mock.ExpectQuery(`SELECT price_per_night FROM rooms WHERE id`).
		WithArgs(roomID, hotelID).
		WillReturnError(sql.ErrNoRows)

	price, err := repo.GetRoomPrice(context.Background(), hotelID, roomID)
	assert.Error(t, err)
	assert.Equal(t, 0.0, price)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateRoom_WithZeroPrice(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	room := &domain.Room{
		ID:            "room-123",
		HotelID:       "hotel-123",
		RoomNumber:    "101",
		RoomType:      "Standard",
		PricePerNight: 0.0,
		Capacity:      2,
		Description:   "Free room",
		IsAvailable:   true,
	}

	createdAt := time.Now()
	updatedAt := time.Now()

	mock.ExpectQuery(`INSERT INTO rooms`).
		WithArgs(
			room.ID, room.HotelID, room.RoomNumber, room.RoomType,
			room.PricePerNight, room.Capacity, room.Description, room.IsAvailable,
		).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
			AddRow(createdAt, updatedAt))

	err := repo.CreateRoom(context.Background(), room)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRoomsByHotel_OrderedByRoomNumber(t *testing.T) {
	db, mock := setupMockDBForRoom(t)
	defer db.Close()

	repo := NewPostgresRoomRepository(db)
	hotelID := "hotel-123"
	createdAt := time.Now()
	updatedAt := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "hotel_id", "room_number", "room_type", "price_per_night",
		"capacity", "description", "is_available", "created_at", "updated_at",
	}).
		AddRow("room-3", hotelID, "103", "Standard", 5000.0, 2, "Room 3", true, createdAt, updatedAt).
		AddRow("room-1", hotelID, "101", "Standard", 5000.0, 2, "Room 1", true, createdAt, updatedAt).
		AddRow("room-2", hotelID, "102", "Standard", 5000.0, 2, "Room 2", true, createdAt, updatedAt)

	mock.ExpectQuery(`SELECT.*FROM rooms WHERE hotel_id.*ORDER BY room_number`).
		WithArgs(hotelID).
		WillReturnRows(rows)

	rooms, err := repo.GetRoomsByHotel(context.Background(), hotelID)
	assert.NoError(t, err)
	assert.Len(t, rooms, 3)
	assert.NoError(t, mock.ExpectationsWereMet())
}
