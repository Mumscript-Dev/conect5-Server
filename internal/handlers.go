package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var UpgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true //allow all connections
	},
}

type WsJsonResponse struct {
	Message string `json:"message"`
	Action string `json:"action"`
	MessageType string `json:"message_type"`
}

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

func (app *application) UserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User"))
}

func (app *application) ChatHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := UpgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error in upgrading connection")
		app.errorLog.Println(err)
		return
	}

	var response WsJsonResponse
	response.Message = `<em><small>Connected to server</small></em>`
	err = ws.WriteJSON(response)
	if err != nil {
		app.errorLog.Println(err)
	}
	defer ws.Close()
}
