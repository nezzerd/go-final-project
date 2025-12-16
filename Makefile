.PHONY: test coverage docker-up docker-down migrate-hotel migrate-booking seed-hotel

test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

coverage:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate-hotel:
	migrate -path ./migrations/hotel -database "postgresql://hotel_user:hotel_pass@localhost:5432/hotel_db?sslmode=disable" up

migrate-booking:
	migrate -path ./migrations/booking -database "postgresql://booking_user:booking_pass@localhost:5433/booking_db?sslmode=disable" up

seed-hotel:
	go run cmd/seed/hotel/main.go

build-all:
	go build -o bin/hotel-service cmd/hotel-service/main.go
	go build -o bin/booking-service cmd/booking-service/main.go
	go build -o bin/notification-service cmd/notification-service/main.go
	go build -o bin/delivery-service cmd/delivery-service/main.go
	go build -o bin/payment-service cmd/payment-service/main.go
