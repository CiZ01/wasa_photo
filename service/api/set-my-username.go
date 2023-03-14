package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.francescofazzari.it/wasa_photo/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

/*
setMyUsername is the handler for the API endpoint PUT /profiles/:profileUserID/username.
It change the username of the user with the given userID to the new username.
The new username must be in the body of the request.
The request body must be a JSON object with the following fields:
- username: string
If the new username is already taken, the request will fail.
If the user is not authorized, the request will fail.
If the username is not valid, the request will fail.
*/
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the userID from the URL
	profileUserID, err := strconv.Atoi(ps.ByName("profileUserID"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	userID := ctx.UserID

	// Check if the user is authorized
	if profileUserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Read the request body, save the new username in the variable newUsername
	type NewUsernameBody struct {
		NewUsername string `json:"username"`
	}

	var newUsernameBody NewUsernameBody

	if err := json.NewDecoder(r.Body).Decode(&newUsernameBody); err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the new username is valid
	newUsername := newUsernameBody.NewUsername

	if !IsValid(newUsername) {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	// Change the username, if the new username is already taken, the request will fail
	if err := rt.db.ChangeUsername(userID, newUsername); err != nil {
		http.Error(w, "Username already taken. Username must be unique", http.StatusBadRequest)
		return
	}

	// Return a success message
	w.WriteHeader(http.StatusOK)
}
