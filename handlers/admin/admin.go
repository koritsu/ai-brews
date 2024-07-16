package admin

import (
	"ai-brews/models"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
)

type CreateAdminRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateAdmin(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateAdminRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		admin := models.Admin{
			Username: req.Username,
			Password: req.Password, // 실제로는 해싱이 필요합니다.
		}

		if err := db.Create(&admin).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(admin)
	}
}

func ListAdmins(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var admins []models.Admin
		if err := db.Find(&admins).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(admins)
	}
}
