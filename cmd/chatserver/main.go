package main

import (
	"chatappserver/database"
	"chatappserver/server"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
)

const port = ":8000"

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err.Error())
	}
	storage := database.NewStorage()
	defer storage.CloseDBConnection()
	chatServer := server.NewServer(port, storage)
	fmt.Println(storage)
	fmt.Printf("Server running on port %s\n", port)
	err := http.ListenAndServe(chatServer.Port, chatServer.Router)
	if err != nil {
		panic(err.Error())
	}
}
