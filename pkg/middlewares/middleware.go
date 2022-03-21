package middlewares

import (
	"encoding/json"
	"github.com/spf13/cast"
	user_manager "landing-page/manager/user-manager"
	"landing-page/models"
	"landing-page/pkg/constant"
	"landing-page/pkg/utils"
	"net/http"
	"strings"
)

type Middleware struct {
	UserManager user_manager.UserManager
}

func (m *Middleware)CommonMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, Content-Length, "+
			"Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, "+
			"Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, "+
			"Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware)IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = utils.ParseJwt(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware)IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, _ := utils.ParseJwt(cookie.Value)
		user, err := m.UserManager.GetUser(&models.UserRequest{
			Id: cast.ToInt64(id),
		})
		if err != nil{
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if len(user) < 1 || strings.ToUpper(user[0].Role) != constant.RoleAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "access denied",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}
