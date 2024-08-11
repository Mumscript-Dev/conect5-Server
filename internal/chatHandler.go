package main

import (
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/gorilla/websocket"
)

var UpgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all connections
	},
}

type websocketConnection struct {
	*websocket.Conn
}

var wsChan = make(chan WsPayload, 100)
var clients = make(map[websocketConnection]string)
var clientsMutex = sync.Mutex{}

type WsJsonResponse struct {
	Message  string   `json:"message"`
	Profile  int      `json:"profile"`
	User     string   `json:"user"`
	UserList []string `json:"userList"`
	Action   string   `json:"action"`
}

type WsPayload struct {
	Message string             `json:"message"`
	User    string             `json:"user"`
	Profile int                `json:"profile"`
	Action  string             `json:"action"`
	Conn    websocketConnection `json:"-"`
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
	response.Action = "chat"
	err = ws.WriteJSON(response)
	if err != nil {
		app.errorLog.Println(err)
	}

	var payload WsPayload
	err = ws.ReadJSON(&payload)
	if err != nil {
		app.errorLog.Printf("app can read payload because %v", err)
		return
	}
	conn := websocketConnection{Conn: ws}

	clientsMutex.Lock()
	clients[conn] = payload.User
	clientsMutex.Unlock()

	go app.ListenForWs(&conn)
}

func (app *application) ListenForWs(conn *websocketConnection) {
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
		app.infoLog.Printf("%v send : %v ", payload.User, payload.Message)
		payload.Conn = *conn
		wsChan <- payload
	}
}

func ListenForWsChan() {
	var response WsJsonResponse
	fmt.Println("Listening to websocket channel")
	for {
		e := <-wsChan
		switch e.Action {
		case "join":
			clientsMutex.Lock()
			clients[e.Conn] = e.User
			clientsMutex.Unlock()
			response.User = "Server"
			response.Message = e.User + " has joined the chat"
			response.Profile = 0
			response.Action = "join"
			BroadcastMessage(response)
		case "chat":
			response.Message = e.Message
			response.User = e.User
			response.Profile = e.Profile
			response.Action = "chat"
			BroadcastMessage(response)
		case "userList":
			userList := getUserList()
			response.Message = fmt.Sprintf("Users in chat: %v", userList)
			response.User = e.User
			response.UserList = userList
			response.Profile = 0
			response.Action = "userList"
			BroadToUser(response)
		}
	}
}

func BroadToUser(response WsJsonResponse) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	for client := range clients {
		if clients[client] == response.User {
			err := client.WriteJSON(response)
			if err != nil {
				fmt.Printf("Error in writing message to user %v", response.User)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func BroadcastMessage(response WsJsonResponse) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			fmt.Println("Error in writing message")
			client.Close()
			delete(clients, client)
		}
	}
}

func getUserList() []string {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	var userList []string
	for _, x := range clients {
		userList = append(userList, x)
	}
	sort.Strings(userList)
	return userList
}
