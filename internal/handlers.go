package main

import (
	"fmt"
	"net/http"
	"sync"
	"github.com/gorilla/websocket"
)

var UpgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true //allow all connections
	},
}

type Game struct {
	Player1 *websocket.Conn
	Player2 *websocket.Conn
	Mux     sync.Mutex
}
var currentGame = &Game{}
type WsJsonResponse struct {
	Message string `json:"message"`
	Profile int `json:"profile"`
	User string `json:"user"`
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
	response.Message = "Connected to server"
	response.Action = "2,3"
	err = ws.WriteJSON(response)
	if err != nil {
		app.errorLog.Println(err)
	}

	for {
		var msg WsJsonResponse
		err := ws.ReadJSON(&msg)
		if err != nil {
				app.errorLog.Println("Error reading json.", err)
				break
		}
		app.infoLog.Printf("Received message: %s", msg.Message)

		// Optionally, send a response back to the client
		response.Message = fmt.Sprintf("Message received: %s", msg.Message)
		err = ws.WriteJSON(response)
		if err != nil {
				app.errorLog.Println(err)
				break
		}
}
	defer ws.Close()
}

func (app *application) GameHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := UpgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		app.errorLog.Println("Error upgrading connection:", err)
		return
	}

	currentGame.Mux.Lock()
	defer currentGame.Mux.Unlock()

	if currentGame.Player1 == nil {
		currentGame.Player1 = ws
		app.infoLog.Println("Player 1 connected")
	} else if currentGame.Player2 == nil {
		currentGame.Player2 = ws
		app.infoLog.Println("Player 2 connected")
	} else {
		app.errorLog.Println("Game already has two players")
		ws.Close()
		return
	}

	defer func() {
		if currentGame.Player1 == ws {
			currentGame.Player1 = nil
		} else if currentGame.Player2 == ws {
			currentGame.Player2 = nil
		}
		ws.Close()
	}()

	for {
		var msg WsJsonResponse
		err := ws.ReadJSON(&msg)
		if err != nil {
			app.errorLog.Println("Error reading json:", err)
			break
		}

		// Forward the message to the other player
		var otherPlayer *websocket.Conn
		if ws == currentGame.Player1 {
			otherPlayer = currentGame.Player2
		} else {
			otherPlayer = currentGame.Player1
		}

		if otherPlayer != nil {
			err = otherPlayer.WriteJSON(msg)
			if err != nil {
				app.errorLog.Println("Error sending message to other player:", err)
				break
			}
		} else {
			app.infoLog.Println("Waiting for the other player to connect")
		}
	}


}