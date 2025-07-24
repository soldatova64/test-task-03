package controllers

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"test-task-03/entity"
	"test-task-03/models"
	"test-task-03/models/responses"
	subscriptionRepo "test-task-03/repositories/subscription"
)

func (app *App) GetSubscriptionV1(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "Неверный идентификатор подписки", http.StatusBadRequest)
		return
	}

	subscription, err := subscriptionRepo.NewRepository(app.db).Get(id)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Подписка не найдена", http.StatusNotFound)
		} else {
			http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)
}

func (app *App) GetListSubscriptionV1(writer http.ResponseWriter, request *http.Request) {
	subscription, err := subscriptionRepo.NewRepository(app.db).GetList()

	if err != nil {
		http.Error(writer, "Ошибка базы данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := responses.SubscriptionResponse{
		Meta: models.Meta{Action: "home"},
		Data: subscription,
	}

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, "Ошибка формирования JSON.", http.StatusInternalServerError)
	}
}

func (app *App) CreateSubscriptionV1(w http.ResponseWriter, r *http.Request) {
	var subscription entity.Subscription
	err := json.NewDecoder(r.Body).Decode(&subscription)

	if subscription.EndDate == nil {
		newEndDate := subscription.StartDate.AddDate(100, 0, 0)
		subscription.EndDate = &newEndDate
	}

	if err != nil {
		http.Error(w, "Недопустимый JSON", http.StatusBadRequest)
		return
	}

	if subscription.ServiceName == "" {
		http.Error(w, "Поле service_name обязательно", http.StatusBadRequest)
		return
	}

	if subscription.Price <= 0 {
		http.Error(w, "Поле price должно быть положительным числом", http.StatusBadRequest)
		return
	}

	if subscription.UserID == "" {
		http.Error(w, "Поле user_id обязательно", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(subscription.UserID); err != nil {
		http.Error(w, "Неверный формат UserID (ожидается UUID)", http.StatusBadRequest)
		return
	}

	if subscription.StartDate.IsZero() {
		http.Error(w, "Поле start_date обязательно", http.StatusBadRequest)
		return
	}

	if subscription.EndDate != nil && subscription.EndDate.Before(subscription.StartDate) {
		http.Error(w, "Дата окончания не может быть раньше даты начала", http.StatusBadRequest)
		return
	}

	createdSub, err := subscriptionRepo.NewRepository(app.db).Post(subscription)

	if err != nil {
		http.Error(w, "Ошибка при создании подписки", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdSub)
}

func (app *App) UpdateSubscriptionV1(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Println(err)
		http.Error(w, "Неверный идентификатор подписки", http.StatusBadRequest)
		return
	}

	subscription, err := subscriptionRepo.NewRepository(app.db).Get(id)

	if err != nil {
		http.Error(w, "Подписка не найдена", http.StatusNotFound)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&subscription)

	if err != nil {
		http.Error(w, "Недопустимый текст запроса", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(subscription.UserID); err != nil {
		http.Error(w, "Неверный формат UserID (ожидается UUID)", http.StatusBadRequest)
		return
	}

	if subscription.EndDate != nil && subscription.EndDate.Before(subscription.StartDate) {
		http.Error(w, "Дата окончания не может быть раньше даты начала", http.StatusBadRequest)
		return
	}

	_, err = subscriptionRepo.NewRepository(app.db).Put(subscription)
	if err != nil {
		http.Error(w, "Ошибка при обновлении подписки", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *App) DeleteSubscriptionV1(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный идентификатор подписки", http.StatusBadRequest)
		return
	}

	err = subscriptionRepo.NewRepository(app.db).Delete(id)
	if err != nil {
		http.Error(w, "Ошибка при удалении подписки", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *App) GetSubscriptionSumV1(w http.ResponseWriter, r *http.Request) {
	sum, err := subscriptionRepo.NewRepository(app.db).GetSubscriptionSum(
		r.URL.Query().Get("user_id"),
		r.URL.Query().Get("service_name"),
		r.URL.Query().Get("date_from"),
		r.URL.Query().Get("date_to"),
	)

	if err != nil {
		http.Error(w, "Ошибка при вычислении суммы", http.StatusInternalServerError)
		return
	}

	response := struct {
		Sum int `json:"sum"`
	}{
		Sum: sum,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
