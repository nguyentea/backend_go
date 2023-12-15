package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/m/handlers"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func main() {
	handlers.ConnectDB()
	defer handlers.CloseDB()

	http.HandleFunc("/createuser", handlers.CreateUserEndpoint)
	http.HandleFunc("/getuser", handlers.GetUserEndpoint)
	http.HandleFunc("/login", handlers.LoginHandler)

	fmt.Println("Server is running on :8888")
	log.Fatal(http.ListenAndServe("0.0.0.0:8888", nil))
}
