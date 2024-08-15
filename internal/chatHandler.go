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
			app.errorLog.Printf("app can not read payload because %v", err)
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
			clients[e.Conn] = fmt.Sprintf("%v-%v", e.User, e.Profile)
			response.User = e.User
			userList := getUserList()
			response.Profile = e.Profile
			response.Message = "Server has fetched the user list"
			response.Action = "userList"
			response.UserList = userList
			BroadcastToUser(response)
		case "leave":
			fmt.Sprintln("someone has left the chat")
			clientsMutex.Lock()
			delete(clients, e.Conn)
			clientsMutex.Unlock()
			response.Action = "leave"
			response.Message = e.User + " has left the chat"
			response.User = "Server"
			response.Profile = 0
			BroadcastMessage(response)
		}
	}
}

func BroadcastToUser(response WsJsonResponse) {
	fmt.Println(128, response.User)
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	user := fmt.Sprintf("%v-%v", response.User, response.Profile)
	for client := range clients {
		if clients[client] == user {
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
	fmt.Println(164, userList)
	return userList
}
