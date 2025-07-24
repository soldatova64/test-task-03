#Сервис управления онлайн подписками
##Описание
Сервис предоставляет REST API для управления онлайн подписками, включая создание, чтение, обновление и удаление подписок, 
а также расчет их суммарной стоимости с возможностью фильтрации.

##Технологии
Язык программирования: Go
База данных: PostgreSQL
Маршрутизация: Gorilla Mux
Миграции: golang-migrate
Контейнеризация: Docker

Для установки и запуска должны быть установлены Docker и Docker Compose.

Создайте файл .env в корне проекта с необходимыми переменными окружения:
env
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_PORT=5432

Запустите сервис с помощью Docker Compose: docker compose up -d
Сервис будет доступен по адресу: http://localhost

##API Endpoints
Подписки
GET /v1/subscription - Получить список всех активных подписок
POST /v1/subscription - Создать новую подписку
GET /v1/subscription/{id} - Получить подписку по ID
PUT /v1/subscription/{id} - Обновить подписку
DELETE /v1/subscription/{id} - Удалить подписку (мягкое удаление)

Расчет стоимости
GET /v1/subscription/sum - Получить суммарную стоимость подписок с возможностью фильтрации