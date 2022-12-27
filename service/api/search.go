package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.francescofazzari.it/wasa_photo/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the search query from the request
	query_search := r.URL.Query().Get("search")
	if query_search == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	from_follow := false
	query_from_follow := r.URL.Query().Get("from_follow")
	if query_from_follow == "1" {
		from_follow = true
	}

	limit, offset, err := getLimitAndOffset(r.URL.Query())
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	var userID uint32
	auth := r.Header["Authorization"][0]
	if auth == "" {
		http.Error(w, "Missing Authorization header", http.StatusBadRequest)
		return
	} else {
		_userID, err := strconv.Atoi(auth)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		userID = uint32(_userID)
	}

	dbUsers, err := rt.db.SearchUsers(userID, query_search, from_follow, offset, limit)
	if err != nil {
		ctx.Logger.Error("Error searching users", err)
		http.Error(w, "Error searching users", http.StatusInternalServerError)
		return
	}

	var users []User
	for _, u := range dbUsers {
		var user User
		user.FromDatabase(u)
		users = append(users, user)
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		ctx.Logger.Error("Error encoding users", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

}