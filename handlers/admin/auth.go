package admin

import (
	"ai-brews/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	_ "github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AdminLogin(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var admin models.Admin
		if err := db.Where("username = ?", req.Username).First(&admin).Error; err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// 시간 기반 일회용 비밀번호(TOTP) 검증
		valid := totp.Validate(req.Password, admin.Password)
		if valid {
			fmt.Println("Code is valid!")
		} else {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Create session
		session, _ := store.Get(r, "admin-session")
		session.Values["authenticated"] = true
		session.Save(r, w)

		w.WriteHeader(http.StatusOK)
	}
}

func AdminLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "admin-session")
		session.Values["authenticated"] = false
		session.Save(r, w)

		w.WriteHeader(http.StatusOK)
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "admin-session")

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
