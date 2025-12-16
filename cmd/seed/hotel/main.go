package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"hotel-booking-system/internal/hotel/domain"
	"hotel-booking-system/internal/hotel/repository"
	"hotel-booking-system/pkg/database"
	"hotel-booking-system/pkg/logger"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	logger.Init("info")
	log := logger.GetLogger()

	log.Info("waiting for database to be ready...")
	time.Sleep(5 * time.Second)

	dbCfg := database.Config{
		Host:     os.Getenv("HOTEL_DB_HOST"),
		Port:     os.Getenv("HOTEL_DB_PORT"),
		User:     os.Getenv("HOTEL_DB_USER"),
		Password: os.Getenv("HOTEL_DB_PASSWORD"),
		DBName:   os.Getenv("HOTEL_DB_NAME"),
	}

	log.Infof("connecting to database: %s:%s", dbCfg.Host, dbCfg.Port)

	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = database.NewPostgresConnection(dbCfg)
		if err == nil {
			break
		}
		log.WithError(err).Warnf("failed to connect to database, retry %d/10", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.WithError(err).Fatal("failed to connect to database after retries")
	}
	defer db.Close()

	hotelRepo := repository.NewPostgresHotelRepository(db)
	roomRepo := repository.NewPostgresRoomRepository(db)

	ctx := context.Background()

	hotels := []domain.Hotel{
		{ID: uuid.New().String(), Name: "Гранд Отель Москва", Description: "Роскошный отель в центре Москвы", Address: "Москва, Тверская улица, 1", OwnerID: uuid.New().String()},
		{ID: uuid.New().String(), Name: "Отель Санкт-Петербург", Description: "Комфортабельный отель у Невского проспекта", Address: "Санкт-Петербург, Невский проспект, 50", OwnerID: uuid.New().String()},
		{ID: uuid.New().String(), Name: "Казань Плаза", Description: "Современный отель в историческом центре", Address: "Казань, улица Баумана, 25", OwnerID: uuid.New().String()},
		{ID: uuid.New().String(), Name: "Сочи Парк Отель", Description: "Отель у моря с видом на горы", Address: "Сочи, Курортный проспект, 75", OwnerID: uuid.New().String()},
		{ID: uuid.New().String(), Name: "Екатеринбург Центр", Description: "Бизнес-отель в центре города", Address: "Екатеринбург, проспект Ленина, 40", OwnerID: uuid.New().String()},
		{ID: uuid.New().String(), Name: "Новосибирск Сити", Description: "Современный отель для деловых путешественников", Address: "Новосибирск, Красный проспект, 15", OwnerID: uuid.New().String()},
		{ID: uuid.New().String(), Name: "Владивосток Океан", Description: "Отель с видом на бухту", Address: "Владивосток, Океанский проспект, 10", OwnerID: uuid.New().String()},
		{ID: uuid.New().String(), Name: "Нижний Новгород Волга", Description: "Отель на берегу Волги", Address: "Нижний Новгород, Верхне-Волжская набережная, 3", OwnerID: uuid.New().String()},
		{ID: uuid.New().String(), Name: "Калининград Европа", Description: "Уютный отель в европейском стиле", Address: "Калининград, проспект Мира, 20", OwnerID: uuid.New().String()},
		{ID: uuid.New().String(), Name: "Красноярск Сибирь", Description: "Комфортный отель в сибирском городе", Address: "Красноярск, проспект Мира, 35", OwnerID: uuid.New().String()},
	}

	roomTypes := []string{"Стандарт", "Улучшенный", "Люкс", "Делюкс", "Президентский люкс"}
	basePrices := []float64{3000, 5000, 8000, 12000, 25000}

	for _, hotel := range hotels {
		if err := hotelRepo.CreateHotel(ctx, &hotel); err != nil {
			log.WithError(err).Errorf("failed to create hotel %s", hotel.Name)
			continue
		}
		log.Infof("created hotel: %s", hotel.Name)

		for i := 1; i <= 5; i++ {
			for j, roomType := range roomTypes {
				room := domain.Room{
					ID:            uuid.New().String(),
					HotelID:       hotel.ID,
					RoomNumber:    fmt.Sprintf("%d%02d", i, j+1),
					RoomType:      roomType,
					PricePerNight: basePrices[j],
					Capacity:      (j/2 + 1) * 2,
					Description:   fmt.Sprintf("Номер типа %s на %d этаже", roomType, i),
					IsAvailable:   true,
				}

				if err := roomRepo.CreateRoom(ctx, &room); err != nil {
					log.WithError(err).Errorf("failed to create room %s", room.RoomNumber)
					continue
				}
			}
		}
		log.Infof("created 25 rooms for hotel: %s", hotel.Name)
	}

	log.Info("seed completed successfully")
}
