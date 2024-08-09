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

type websocketConnection struct {	
	*websocket.Conn
}

var wsChan = make(chan WsPayload, 100)

var clients = make(map[websocketConnection]string)

type WsJsonResponse struct {
	Message string `json:"message"`
	Profile int `json:"profile"`
	User string `json:"user"`
}
type WsPayload struct {
	Message string `json:"message"`
	User string `json:"user"`
	Profile int `json:"profile"`
	Conn websocketConnection `json:"-"`
}
func (app *application) ChatHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := UpgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error in upgrading connection")
		app.errorLog.Println(err)
		return
	}
	// when connection is successful, send a message to the client to inform them that they are connected
	var response WsJsonResponse
	response.Profile = 0
	response.Message = "Connected to server"
	response.User = "Server"
	err = ws.WriteJSON(response)
	if err != nil {
		app.errorLog.Println(err)
	}

	conn := websocketConnection{Conn: ws}
	clients[conn] = ""

	go app.ListenForWs(&conn)
}

func(app *application) ListenForWs(conn *websocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			app.errorLog.Printf("Websocket connection error: %v", r)
		}
	}()
	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			app.errorLog.Printf("app can read payload because %v", err)
			return
		}
		app.infoLog.Printf("%v send : %v ", payload.User, payload.Message )
		payload.Conn = *conn
		wsChan <- payload
	}
}

func ListenForWsChan() {
	var response WsJsonResponse
	fmt.Println("Listening to websocket channel")
	for {
		e := <-wsChan
		response.Message = e.Message
		response.User = e.User
		response.Profile = e.Profile
BroadcastMessage(response)
	}
}

func BroadcastMessage(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			fmt.Println("Error in writing message")
			client.Close()
			delete(clients, client)
		}
	}
}
