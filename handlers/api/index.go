package api

import (
	"ai-brews/models"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
)

func Index(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var recipes []models.Recipe
		if err := db.Find(&recipes).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(recipes)
	}
}
