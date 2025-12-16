# Руководство по использованию сервисов Hotel Booking

## Содержание
1. [Обзор сервисов](#обзор-сервисов)
2. [Развертывание инфраструктуры](#развертывание-инфраструктуры)
3. [Детальное описание сервисов](#детальное-описание-сервисов)
4. [Архитектура](#архитектура)

---

## Обзор сервисов

**Перечень сервисов:**
- **Hotel Service** — HTTP `http://localhost:8081` — управление отелями и номерами
- **Booking Service** — HTTP `http://localhost:8082` — управление бронированиями
- **Notification Service** — фоновый сервис — обработка событий и отправка уведомлений
- **Delivery Service** — HTTP `http://localhost:8084` — доставка уведомлений по различным каналам
- **Payment Service** — HTTP `http://localhost:8085` — обработка платежей

**Порты (настраиваются через `.env`):**
- `HOTEL_SERVICE_PORT=8081`
- `BOOKING_SERVICE_PORT=8082`
- `NOTIFICATION_SERVICE_PORT=8083` (используется только для метрик)
- `DELIVERY_SERVICE_PORT=8084`
- `PAYMENT_SERVICE_PORT=8085`
- `PROMETHEUS_PORT=2112` (базовый порт для метрик)

**Метрики Prometheus:**
- Hotel Service: `http://localhost:2112/metrics`
- Booking Service: `http://localhost:2113/metrics`
- Notification Service: `http://localhost:2114/metrics`
- Delivery Service: `http://localhost:2115/metrics`
- Payment Service: `http://localhost:2116/metrics`

**Веб-интерфейсы:**
- Prometheus UI: `http://localhost:9090`
- Jaeger UI: `http://localhost:16686`

---

## Развертывание инфраструктуры

### Шаг 1: Клонирование репозитория

Создайте файл `.env` на основе `env.example`:

```bash
cp env.example .env
```

Откройте `.env` и при необходимости настройте переменные окружения.

### Шаг 2: Запуск через Docker Compose

Запустите все сервисы одной командой:

```bash
docker-compose up -d
```

Эта команда запустит:
- PostgreSQL (2 базы данных: hotel_db и booking_db)
- Kafka и Zookeeper
- Jaeger
- Prometheus
- Seed Service (заполнение базы данных тестовыми данными)
- Все микросервисы

### Шаг 3: Проверка работоспособности

Проверьте доступность сервисов:

```bash
# Hotel Service
curl http://localhost:8081/api/hotels?limit=5

# Booking Service
curl http://localhost:8082/api/bookings/user/test-user

# Delivery Service
curl http://localhost:8084/api/notifications/send \
  -H "Content-Type: application/json" \
  -d '{"channel":"email","recipient":"test@example.com","message":"Test"}'

# Payment Service
curl http://localhost:8085/api/payments \
  -H "Content-Type: application/json" \
  -d '{"booking_id":"test-123","amount":1000.0}'
```

### Шаг 4: Просмотр логов

```bash
# Все сервисы
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f hotel-service
docker-compose logs -f booking-service
docker-compose logs -f notification-service
docker-compose logs -f delivery-service
docker-compose logs -f payment-service
```

## Детальное описание сервисов

### Hotel Service — API (`http://localhost:8081`)

Управление отелями и номерами.

#### Endpoints

**GET** `/api/hotels` — получить список отелей
- Ответ: массив объектов `Hotel`
- Пример:
  ```bash
  curl http://localhost:8081/api/hotels
  ```

**GET** `/api/hotels/{id}` — получить детали отеля
- Ответ: объект `Hotel`
- Пример:
  ```bash
  curl http://localhost:8081/api/hotels/{hotel-id}
  ```

**POST** `/api/hotels` — создать отель
- Body JSON:
  ```json
  {
    "name": "Grand Hotel",
    "description": "Роскошный отель в центре города",
    "address": "ул. Ленина, д. 1, Москва",
    "owner_id": "550e8400-e29b-41d4-a716-446655440000"
  }
  ```
- **Важно:** `owner_id` должен быть валидным UUID (используйте `uuidgen` для генерации)
- Ответ: созданный объект `Hotel` (HTTP 201)
- Пример:
  ```bash
  OWNER_ID=$(uuidgen)
  curl -X POST http://localhost:8081/api/hotels \
    -H "Content-Type: application/json" \
    -d "{
      \"name\": \"Grand Hotel\",
      \"description\": \"Роскошный отель\",
      \"address\": \"ул. Ленина, д. 1\",
      \"owner_id\": \"$OWNER_ID\"
    }"
  ```

**PUT** `/api/hotels/{id}` — обновить отель
- Body JSON:
  ```json
  {
    "name": "Updated Hotel Name",
    "description": "Обновленное описание",
    "address": "ул. Новая, д. 2",
    "owner_id": "550e8400-e29b-41d4-a716-446655440000"
  }
  ```
- Ответ: обновленный объект `Hotel` (HTTP 200)

**GET** `/api/hotels/{id}/rooms` — получить отель со всеми номерами
- Ответ: объект `HotelWithRooms`
  ```json
  {
    "hotel": {
      "id": "uuid",
      "name": "Grand Hotel",
      "description": "...",
      "address": "...",
      "owner_id": "uuid",
      "created_at": "timestamp",
      "updated_at": "timestamp"
    },
    "rooms": [
      {
        "id": "uuid",
        "hotel_id": "uuid",
        "room_number": "101",
        "room_type": "Standard",
        "price_per_night": 5000.0,
        "capacity": 2,
        "description": "Стандартный номер",
        "is_available": true,
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ]
  }
  ```

**POST** `/api/rooms` — создать номер
- Body JSON:
  ```json
  {
    "hotel_id": "550e8400-e29b-41d4-a716-446655440000",
    "room_number": "101",
    "room_type": "Standard",
    "price_per_night": 5000.0,
    "capacity": 2,
    "description": "Стандартный номер с видом на город",
    "is_available": true
  }
  ```
- Ответ: созданный объект `Room` (HTTP 201)

#### JSON схемы

**Hotel:**
```json
{
  "id": "uuid",
  "name": "string",
  "description": "string",
  "address": "string",
  "owner_id": "uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

**Room:**
```json
{
  "id": "uuid",
  "hotel_id": "uuid",
  "room_number": "string",
  "room_type": "string",
  "price_per_night": "float64",
  "capacity": "int",
  "description": "string",
  "is_available": "bool",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

---

### Booking Service — API (`http://localhost:8082`)

Управление бронированиями.

#### Endpoints

**POST** `/api/bookings` — создать бронирование
- Body JSON:
  ```json
  {
    "user_id": "user-123",
    "hotel_id": "550e8400-e29b-41d4-a716-446655440000",
    "room_id": "550e8400-e29b-41d4-a716-446655440000",
    "check_in_date": "2024-12-20T14:00:00Z",
    "check_out_date": "2024-12-25T12:00:00Z"
  }
  ```
- **Формат дат:** RFC3339 (ISO 8601), например: `2024-12-20T14:00:00Z`
- **Важно:** `user_id` может быть любой строкой (VARCHAR(255) в БД)
- Ответ: объект `Booking` (HTTP 201)
- Сервис автоматически:
    1. Проверяет доступность комнаты через Hotel Service (HTTP запрос)
    2. Получает цену за ночь
    3. Рассчитывает `total_price` на основе количества ночей
    4. Устанавливает `status` = `"confirmed"` и `payment_status` = `"pending"`
    5. Публикует Kafka-событие в топик `booking.created`
    6. Создает платеж через Payment Service
- Пример:
  ```bash
  curl -X POST http://localhost:8082/api/bookings \
    -H "Content-Type: application/json" \
    -d '{
      "user_id": "user-123",
      "hotel_id": "<hotel-uuid>",
      "room_id": "<room-uuid>",
      "check_in_date": "2024-12-20T14:00:00Z",
      "check_out_date": "2024-12-25T12:00:00Z"
    }'
  ```

**GET** `/api/bookings/{id}` — получить бронирование по ID
- Ответ: объект `Booking`

**GET** `/api/bookings/user/{userId}` — получить все бронирования пользователя
- Ответ: массив объектов `Booking`

**GET** `/api/bookings/hotel/{hotelId}` — получить все бронирования отеля
- Ответ: массив объектов `Booking`

**POST** `/api/webhooks/payment` — webhook для обновления статуса оплаты
- Body JSON:
  ```json
  {
    "payment_id": "payment-uuid",
    "booking_id": "booking-uuid",
    "status": "paid",
    "amount": 1000.0,
    "processed_at": "2024-12-15T10:00:00Z"
  }
  ```
- **Возможные статусы:** `pending`, `paid`, `failed`, `refunded`
- Ответ: HTTP 200 OK (пустое тело)
- Используется Payment Service для уведомления о статусе платежа

#### JSON схема

**Booking:**
```json
{
  "id": "uuid",
  "user_id": "string",
  "hotel_id": "uuid",
  "room_id": "uuid",
  "check_in_date": "timestamp (RFC3339)",
  "check_out_date": "timestamp (RFC3339)",
  "total_price": 25000.0,
  "status": "confirmed|pending|cancelled",
  "payment_status": "pending|paid|failed|refunded",
  "created_at": "timestamp (RFC3339)",
  "updated_at": "timestamp (RFC3339)"
}
```

---

### Delivery Service — API (`http://localhost:8084`)

Доставка уведомлений по различным каналам.

#### Endpoints

**POST** `/api/notifications/send` — отправить уведомление
- Body JSON:
  ```json
  {
    "channel": "email",
    "recipient": "user@example.com",
    "subject": "Бронирование подтверждено",
    "message": "Ваше бронирование успешно создано!"
  }
  ```
- **Каналы:** `email`, `sms`, `telegram`
- **Поля:**
    - `channel` (обязательно) — канал доставки
    - `recipient` (обязательно) — получатель (email, телефон, telegram chat_id)
    - `subject` (опционально) — тема (для email)
    - `message` (обязательно) — текст сообщения
- Ответ:
  ```json
  {
    "success": true,
    "message": "notification sent successfully"
  }
  ```
- Пример:
  ```bash
  curl -X POST http://localhost:8084/api/notifications/send \
    -H "Content-Type: application/json" \
    -d '{
      "channel": "email",
      "recipient": "user@example.com",
      "subject": "Тест",
      "message": "Тестовое сообщение"
    }'
  ```

---

### Payment Service — API (`http://localhost:8085`)

Обработка платежей.

#### Endpoints

**POST** `/api/payments` — создать платеж
- Body JSON:
  ```json
  {
    "booking_id": "550e8400-e29b-41d4-a716-446655440000",
    "amount": 1000.0,
    "currency": "RUB"
  }
  ```
- **Поля:**
    - `booking_id` (обязательно) — ID бронирования
    - `amount` (обязательно) — сумма платежа
    - `currency` (опционально) — валюта (по умолчанию `RUB`)
- Ответ: HTTP 202 Accepted
  ```json
  {
    "payment_id": "payment-uuid",
    "status": "processing",
    "message": "payment is being processed"
  }
  ```
- Сервис асинхронно обрабатывает платеж и отправляет webhook в Booking Service
- Пример:
  ```bash
  curl -X POST http://localhost:8085/api/payments \
    -H "Content-Type: application/json" \
    -d '{
      "booking_id": "booking-uuid",
      "amount": 1000.0,
      "currency": "RUB"
    }'
  ```

---

### Notification Service

Фоновый сервис без HTTP API. Работает как Kafka consumer и отправляет уведомления через Delivery Service.

#### Функционал

- Подписывается на топик `booking.created` в Kafka
- При получении события о создании бронирования:
    1. Получает `owner_id` отеля через Hotel Service
    2. Отправляет уведомление клиенту через Delivery Service
    3. Отправляет уведомление владельцу отеля через Delivery Service

---
## Архитектура

### Структура проекта

```
/
├── cmd/                    # Точки входа приложений
│   ├── booking-service/
│   ├── hotel-service/
│   ├── notification-service/
│   ├── delivery-service/
│   └── payment-service/
│
├── internal/              # Приватный код приложения
│   └── {service}/
│       ├── domain/        # Доменные модели и интерфейсы
│       ├── repository/    # Реализация репозиториев
│       ├── usecase/       # Бизнес-логика
│       └── delivery/      # HTTP handlers и routes
│
├── pkg/                   # Публичные библиотеки
│   ├── database/          # Подключение к PostgreSQL
│   ├── hotelclient/       # HTTP клиент для Hotel Service
│   ├── httpclient/        # HTTP клиенты (Delivery, Payment, Hotel)
│   ├── kafka/             # Producer и Consumer для Kafka
│   ├── logger/            # Структурированное логирование
│   ├── metrics/           # Prometheus метрики
│   └── tracing/           # Jaeger трейсинг
│
├── migrations/            # SQL миграции БД
│   ├── hotel/
│   └── booking/
│
└── deployments/           # Docker конфигурации
