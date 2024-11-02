package main

import (
	"net/http"
)


func (app *application) UserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}
// type userParams struct {
// 	Email string `json:"email"`
// 	ProfileIndex string `json:"profileIndex"`
// 	Username string `json:"username"` 
// }

// type userResponse struct {
// 	Email string `json:"email"`
// 	ProfileIndex string `json:"profileIndex"`
// 	Username string `json:"username"`
// 	ID string `json:"id"`
// }

// func (app *application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
// 	decoder := json.NewDecoder(r.Body)
// 	params := userParams{}
// 	err := decoder.Decode(&params)
// 	if err != nil {
// 		app.errorLog.Print(err)
// 		w.WriteHeader(400)
// 		w.Write([]byte("Invalid request"))
// 	}
// 	// Generate a new UUID for the user ID
// 	id := uuid.New()

// 	// Call the CreateUser query generated by sqlc
// 	ctx := context.Background()
// 	user, err := app.queries.CreateUser(ctx, database.CreateUserParams{
// 		ID:           id,
// 		Email:        params.Email,
// 		Username:     sql.NullString{String: params.Username, Valid: true},
// 		Profileindex: sql.NullString{String: params.ProfileIndex, Valid: true},
// 	})
// 	if err != nil {
// 		app.errorLog.Print(err)
// 		w.WriteHeader(500)
// 		w.Write([]byte("Error creating user"))
// 		return
// 	}
// 	userCreated := userResponse{
// 		Email:        user.Email,
// 		ProfileIndex: user.Profileindex.String,
// 		Username:     user.Username.String,
// 		ID:           user.ID.(string),
// 	}
// 	data, err := json.Marshal(userCreated)
// 	if err != nil {
// 		app.errorLog.Print(err)
// 		w.WriteHeader(500)
// 		w.Write([]byte("Error creating user"))
// 		return
// 	}
// 	w.Header().Add("Content-Type", "application/json")
// 	w.WriteHeader(200)
// 	w.Write(data)
// }