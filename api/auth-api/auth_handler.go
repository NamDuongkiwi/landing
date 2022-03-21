package auth_api

import (
	"database/sql"
	"encoding/json"
	"github.com/spf13/cast"
	user_manager "landing-page/manager/user-manager"
	"landing-page/models"
	"landing-page/pkg/utils"
	"net/http"
	"time"
)

type AuthHandler struct {
	UserManager user_manager.UserManager
}

func (h *AuthHandler)Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "request is invalid", http.StatusBadRequest)
		return
	}
	users, err := h.UserManager.GetUser(&models.UserRequest{
		Email: req.Email,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "email does not exist", http.StatusBadRequest)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	user := users[0]
	if err := user.ComparePassword(req.Password); err != nil {
		http.Error(w, "password incorrect", http.StatusBadRequest)
		return
	}

	token, err := utils.GenerateJwt(cast.ToString(user.Id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	})

	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler)Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "logout successfully",
	})
}
