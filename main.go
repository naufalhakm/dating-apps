package main

import (
	"go-dating-test/app/controllers"
	"go-dating-test/app/repositories"
	"go-dating-test/app/services"
	"go-dating-test/database"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db, err := database.NewPostgresSQLClient()
	if err != nil {
		log.Fatal("Could not connect to PostgresSQL:", err)
	}

	userRepo := repositories.NewUserRepository()
	userSvc := services.NewUserService(db, userRepo)
	userCtrl := controllers.NewUserController(userSvc)

	router := mux.NewRouter()
	router.HandleFunc("/api/match/recommendations", userCtrl.GetUserRecomendation).Methods("GET")

	log.Printf("Server running on :8000\n")
	log.Fatal(http.ListenAndServe(":8000", router))
}
