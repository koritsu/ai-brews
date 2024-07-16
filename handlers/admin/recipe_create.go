package admin

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"ai-brews/config"
	"ai-brews/models"
	"gorm.io/gorm"
)

type CreateRecipeRequest struct {
	Title     string `json:"title"`
	ImageURL  string `json:"image_url"`
	RecipeURL string `json:"recipe_url"`
	ShareURL  string `json:"share_url"`
	Email     string `json:"email"`
}

func CreateRecipe(db *gorm.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRecipeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Use default image if ImageURL is empty
		if req.ImageURL == "" {
			rand.Seed(time.Now().UnixNano())
			defaultImages := []string{
				"/static/default_image1.jpg",
				"/static/default_image2.jpg",
				"/static/default_image3.jpg",
				"/static/default_image4.jpg",
				"/static/default_image5.jpg",
				"/static/default_image6.jpg",
				"/static/default_image7.jpg",
				"/static/default_image8.jpg",
				"/static/default_image9.jpg",
				"/static/default_image10.jpg",
			}
			req.ImageURL = defaultImages[rand.Intn(len(defaultImages))]
		}

		recipe := models.Recipe{
			Title:     req.Title,
			ImageURL:  req.ImageURL,
			RecipeURL: req.RecipeURL,
			ShareURL:  req.ShareURL,
			Email:     req.Email,
		}

		if err := db.Create(&recipe).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(recipe)
		if err != nil {
			return
		}
	}
}
