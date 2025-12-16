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

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	return db, mock
}

func TestNewPostgresHotelRepository(t *testing.T) {
	db, _ := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestCreateHotel_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	hotel := &domain.Hotel{
		ID:          "hotel-123",
		Name:        "Grand Hotel",
		Description: "Luxury hotel",
		Address:     "123 Main St",
		OwnerID:     "owner-123",
	}

	createdAt := time.Now()
	updatedAt := time.Now()

	mock.ExpectQuery(`INSERT INTO hotels`).
		WithArgs(hotel.ID, hotel.Name, hotel.Description, hotel.Address, hotel.OwnerID).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
			AddRow(createdAt, updatedAt))

	err := repo.CreateHotel(context.Background(), hotel)
	assert.NoError(t, err)
	assert.Equal(t, createdAt, hotel.CreatedAt)
	assert.Equal(t, updatedAt, hotel.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateHotel_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	hotel := &domain.Hotel{
		ID:          "hotel-123",
		Name:        "Grand Hotel",
		Description: "Luxury hotel",
		Address:     "123 Main St",
		OwnerID:     "owner-123",
	}

	mock.ExpectQuery(`INSERT INTO hotels`).
		WithArgs(hotel.ID, hotel.Name, hotel.Description, hotel.Address, hotel.OwnerID).
		WillReturnError(errors.New("duplicate key"))

	err := repo.CreateHotel(context.Background(), hotel)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetHotelByID_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	hotelID := "hotel-123"
	createdAt := time.Now()
	updatedAt := time.Now()

	mock.ExpectQuery(`SELECT.*FROM hotels WHERE id`).
		WithArgs(hotelID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "address", "owner_id", "created_at", "updated_at",
		}).AddRow(
			hotelID, "Grand Hotel", "Luxury hotel", "123 Main St", "owner-123",
			createdAt, updatedAt,
		))

	hotel, err := repo.GetHotelByID(context.Background(), hotelID)
	assert.NoError(t, err)
	assert.NotNil(t, hotel)
	assert.Equal(t, hotelID, hotel.ID)
	assert.Equal(t, "Grand Hotel", hotel.Name)
	assert.Equal(t, "owner-123", hotel.OwnerID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetHotelByID_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	hotelID := "non-existent"

	mock.ExpectQuery(`SELECT.*FROM hotels WHERE id`).
		WithArgs(hotelID).
		WillReturnError(sql.ErrNoRows)

	hotel, err := repo.GetHotelByID(context.Background(), hotelID)
	assert.Error(t, err)
	assert.Nil(t, hotel)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetHotelByID_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	hotelID := "hotel-123"

	mock.ExpectQuery(`SELECT.*FROM hotels WHERE id`).
		WithArgs(hotelID).
		WillReturnError(errors.New("connection error"))

	hotel, err := repo.GetHotelByID(context.Background(), hotelID)
	assert.Error(t, err)
	assert.Nil(t, hotel)
	assert.Contains(t, err.Error(), "connection error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetHotels_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	limit := 10
	offset := 0
	createdAt := time.Now()
	updatedAt := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "address", "owner_id", "created_at", "updated_at",
	}).
		AddRow("hotel-1", "Hotel 1", "Desc 1", "Addr 1", "owner-1", createdAt, updatedAt).
		AddRow("hotel-2", "Hotel 2", "Desc 2", "Addr 2", "owner-2", createdAt, updatedAt)

	mock.ExpectQuery(`SELECT.*FROM hotels ORDER BY created_at DESC LIMIT`).
		WithArgs(limit, offset).
		WillReturnRows(rows)

	hotels, err := repo.GetHotels(context.Background(), limit, offset)
	assert.NoError(t, err)
	assert.Len(t, hotels, 2)
	assert.Equal(t, "hotel-1", hotels[0].ID)
	assert.Equal(t, "hotel-2", hotels[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetHotels_EmptyResult(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	limit := 10
	offset := 0

	mock.ExpectQuery(`SELECT.*FROM hotels ORDER BY created_at DESC LIMIT`).
		WithArgs(limit, offset).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "address", "owner_id", "created_at", "updated_at",
		}))

	hotels, err := repo.GetHotels(context.Background(), limit, offset)
	assert.NoError(t, err)
	assert.Empty(t, hotels)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetHotels_WithPagination(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	limit := 5
	offset := 10
	createdAt := time.Now()
	updatedAt := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "address", "owner_id", "created_at", "updated_at",
	}).AddRow("hotel-11", "Hotel 11", "Desc 11", "Addr 11", "owner-11", createdAt, updatedAt)

	mock.ExpectQuery(`SELECT.*FROM hotels ORDER BY created_at DESC LIMIT`).
		WithArgs(limit, offset).
		WillReturnRows(rows)

	hotels, err := repo.GetHotels(context.Background(), limit, offset)
	assert.NoError(t, err)
	assert.Len(t, hotels, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetHotels_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	limit := 10
	offset := 0

	mock.ExpectQuery(`SELECT.*FROM hotels ORDER BY created_at DESC LIMIT`).
		WithArgs(limit, offset).
		WillReturnError(errors.New("query error"))

	hotels, err := repo.GetHotels(context.Background(), limit, offset)
	assert.Error(t, err)
	assert.Nil(t, hotels)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetHotels_ScanError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	limit := 10
	offset := 0

	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "address", "owner_id", "created_at", "updated_at",
	}).AddRow(nil, nil, nil, nil, nil, nil, nil)

	mock.ExpectQuery(`SELECT.*FROM hotels ORDER BY created_at DESC LIMIT`).
		WithArgs(limit, offset).
		WillReturnRows(rows)

	hotels, err := repo.GetHotels(context.Background(), limit, offset)
	assert.Error(t, err)
	assert.Nil(t, hotels)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetHotelsByOwner_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	ownerID := "owner-123"
	createdAt := time.Now()
	updatedAt := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "address", "owner_id", "created_at", "updated_at",
	}).
		AddRow("hotel-1", "Hotel 1", "Desc 1", "Addr 1", ownerID, createdAt, updatedAt).
		AddRow("hotel-2", "Hotel 2", "Desc 2", "Addr 2", ownerID, createdAt, updatedAt)

	mock.ExpectQuery(`SELECT.*FROM hotels WHERE owner_id`).
		WithArgs(ownerID).
		WillReturnRows(rows)

	hotels, err := repo.GetHotelsByOwner(context.Background(), ownerID)
	assert.NoError(t, err)
	assert.Len(t, hotels, 2)
	assert.Equal(t, ownerID, hotels[0].OwnerID)
	assert.Equal(t, ownerID, hotels[1].OwnerID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetHotelsByOwner_EmptyResult(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	ownerID := "owner-no-hotels"

	mock.ExpectQuery(`SELECT.*FROM hotels WHERE owner_id`).
		WithArgs(ownerID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "address", "owner_id", "created_at", "updated_at",
		}))

	hotels, err := repo.GetHotelsByOwner(context.Background(), ownerID)
	assert.NoError(t, err)
	assert.Empty(t, hotels)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetHotelsByOwner_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	ownerID := "owner-123"

	mock.ExpectQuery(`SELECT.*FROM hotels WHERE owner_id`).
		WithArgs(ownerID).
		WillReturnError(errors.New("query error"))

	hotels, err := repo.GetHotelsByOwner(context.Background(), ownerID)
	assert.Error(t, err)
	assert.Nil(t, hotels)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateHotel_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	hotel := &domain.Hotel{
		ID:          "hotel-123",
		Name:        "Updated Hotel",
		Description: "Updated description",
		Address:     "Updated address",
	}
	updatedAt := time.Now()

	mock.ExpectQuery(`UPDATE hotels SET`).
		WithArgs(hotel.ID, hotel.Name, hotel.Description, hotel.Address).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).AddRow(updatedAt))

	err := repo.UpdateHotel(context.Background(), hotel)
	assert.NoError(t, err)
	assert.Equal(t, updatedAt, hotel.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateHotel_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	hotel := &domain.Hotel{
		ID:          "non-existent",
		Name:        "Updated Hotel",
		Description: "Updated description",
		Address:     "Updated address",
	}

	mock.ExpectQuery(`UPDATE hotels SET`).
		WithArgs(hotel.ID, hotel.Name, hotel.Description, hotel.Address).
		WillReturnError(sql.ErrNoRows)

	err := repo.UpdateHotel(context.Background(), hotel)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateHotel_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	hotel := &domain.Hotel{
		ID:          "hotel-123",
		Name:        "Updated Hotel",
		Description: "Updated description",
		Address:     "Updated address",
	}

	mock.ExpectQuery(`UPDATE hotels SET`).
		WithArgs(hotel.ID, hotel.Name, hotel.Description, hotel.Address).
		WillReturnError(errors.New("update error"))

	err := repo.UpdateHotel(context.Background(), hotel)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteHotel_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	hotelID := "hotel-123"

	mock.ExpectExec(`DELETE FROM hotels WHERE id`).
		WithArgs(hotelID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteHotel(context.Background(), hotelID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteHotel_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	hotelID := "non-existent"

	mock.ExpectExec(`DELETE FROM hotels WHERE id`).
		WithArgs(hotelID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeleteHotel(context.Background(), hotelID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteHotel_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPostgresHotelRepository(db)
	hotelID := "hotel-123"

	mock.ExpectExec(`DELETE FROM hotels WHERE id`).
		WithArgs(hotelID).
		WillReturnError(errors.New("delete error"))

	err := repo.DeleteHotel(context.Background(), hotelID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delete error")
	assert.NoError(t, mock.ExpectationsWereMet())
}
