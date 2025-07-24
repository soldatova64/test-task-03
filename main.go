package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"test-task-03/controllers"
	"test-task-03/db"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Main: Ошибка загрузки .env файла.")
	}

	db := db.ConnectDB()
	app := controllers.NewApp(db)
	router := mux.NewRouter()

	router.HandleFunc("/v1/subscription", app.GetListSubscriptionV1).Methods("GET")
	router.HandleFunc("/v1/subscription", app.CreateSubscriptionV1).Methods("POST")
	router.HandleFunc("/v1/subscription/{id}", app.GetSubscriptionV1).Methods("GET")
	router.HandleFunc("/v1/subscription/{id}", app.UpdateSubscriptionV1).Methods("PUT")
	router.HandleFunc("/v1/subscription/{id}", app.DeleteSubscriptionV1).Methods("DELETE")
	router.HandleFunc("/v1/subscription/sum", app.GetSubscriptionSumV1).Methods("GET")

	err = http.ListenAndServe(":80", router)
	if err != nil {
		log.Fatal("Main: Ошибка сервера: ", err)
	} else {
		log.Println("Main: Сервер запущен на 80-м порту.")
	}
}
