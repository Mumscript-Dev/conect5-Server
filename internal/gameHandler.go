package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// BoardSize and Connect are the size of the board and the number of pieces in a row needed to win.
const (
	BoardSize = 12 // Size of the board (15x15 for Connect 5)
	Connect   = 5  // Number of pieces in a row needed to win
)

// GameState represents the state of a game.
type GameState struct {
	Board  [BoardSize][BoardSize]string
	Turn   string // player1 or player2
	Winner string
}

// gameStates is a map of game IDs to game states. the string will be the game ID
var gameStates = make(map[string]*GameState)
var gameStatesMutex = sync.Mutex{}

// GameMove represents a move in a game.
// The player is either "X" or "O". 
// gameID will be player1-player2
type GameMove struct {
	GameID string `json:"gameId"`
	Row    int    `json:"row"`
	Col    int    `json:"col"`
	Player string `json:"player"`
}

func (app *application) GameHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := UpgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error in upgrading connection")
		app.errorLog.Println(err)
		return
	}

	go app.ListenForGameMoves(ws)
}

func (app *application) ListenForGameMoves(conn *websocket.Conn) {
	defer func() {
		if r := recover(); r != nil {
			app.errorLog.Printf("WebSocket connection error: %v", r)
		}
	}()
	for {
		var move GameMove
		err := conn.ReadJSON(&move)
		if err != nil {
			app.errorLog.Printf("Error reading game move: %v", err)
			return
		}

		gameStatesMutex.Lock()
		gameState, exists := gameStates[move.GameID]
		if !exists {
			// If the game does not exist, create a new game state
			gameState = &GameState{
				Board:  [BoardSize][BoardSize]string{},
				Turn:   "X",
				Winner: "",
			}
			gameStates[move.GameID] = gameState
		}

		if gameState.Winner != "" {
			// If there's already a winner, no more moves should be accepted.
			conn.WriteJSON(gameState)
			gameStatesMutex.Unlock()
			continue
		}

		if gameState.Board[move.Row][move.Col] == "" && move.Player == gameState.Turn {
			// Valid move
			gameState.Board[move.Row][move.Col] = move.Player

			if checkWinner(gameState, move.Row, move.Col, move.Player) {
				gameState.Winner = move.Player
			} else {
				// Change turn
				if gameState.Turn == "X" {
					gameState.Turn = "O"
				} else {
					gameState.Turn = "X"
				}
			}
		}

		conn.WriteJSON(gameState)
		gameStatesMutex.Unlock()
	}
}

func checkWinner(gameState *GameState, row, col int, player string) bool {
	directions := [][2]int{
		{0, 1},  // Horizontal
		{1, 0},  // Vertical
		{1, 1},  // Diagonal \
		{1, -1}, // Diagonal /
	}

	for _, dir := range directions {
		count := 1 // The piece just placed counts as 1
		// Check in the positive direction
		count += countPieces(gameState, row, col, dir[0], dir[1], player)
		// Check in the negative direction
		count += countPieces(gameState, row, col, -dir[0], -dir[1], player)

		if count >= Connect {
			return true
		}
	}

	return false
}

func countPieces(gameState *GameState, row, col, rowDir, colDir int, player string) int {
	count := 0
	for i := 1; i < Connect; i++ {
		r := row + i*rowDir
		c := col + i*colDir
		if r >= 0 && r < BoardSize && c >= 0 && c < BoardSize && gameState.Board[r][c] == player {
			count++
		} else {
			break
		}
	}
	return count
}
