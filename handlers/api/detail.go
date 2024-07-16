package api

import (
	"ai-brews/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

func Detail(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var recipe models.Recipe
		if err := db.First(&recipe, id).Error; err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(recipe)
	}
}
