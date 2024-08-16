package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

func(app *application) Auth(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	//how to get env variable
	clientID := os.Getenv("CLIENT_ID")
	secret := os.Getenv("SECRET")
	fmt.Println(secret)
	payload := map[string]string{"clientID": clientID, "secret": secret}
	data, err := json.Marshal(payload)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}