package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.francescofazzari.it/wasa_photo/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) GetLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the profileUserID and photoID from the URL
	_profileUserID, err := strconv.Atoi(ps.ByName("profileUserID"))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	profileUserID := uint32(_profileUserID)

	_postID, err := strconv.Atoi(ps.ByName("postID"))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	photoID := uint32(_postID)

	_userID := r.Header.Get("Authorization")

	userID, err := strconv.Atoi(_userID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if _userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	isBanned, err := rt.db.IsBanned(profileUserID, uint32(userID))
	if err != nil {
		ctx.Logger.WithError(err).Error("Error checking if user is banned")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if isBanned {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	offset, limit := uint32(0), uint32(10)
	if ps.ByName("offset") != "" {
		_offset, err := strconv.Atoi(ps.ByName("offset"))
		if err != nil {
			http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
			return
		}
		offset = uint32(_offset)
	}

	if ps.ByName("limit") != "" {
		_limit, err := strconv.Atoi(ps.ByName("limit"))
		if err != nil {
			http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
			return
		}
		limit = uint32(_limit)
	}

	likes, err := rt.db.GetLikes(photoID, profileUserID, offset, limit)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error getting likes")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(likes)
}
