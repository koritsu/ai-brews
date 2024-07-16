package admin

import (
	"ai-brews/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"time"
)

func UpdateRecipe(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var updatedRecipe models.Recipe
		if err := json.NewDecoder(r.Body).Decode(&updatedRecipe); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Use default image if ImageURL is empty
		if updatedRecipe.ImageURL == "" {
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
			updatedRecipe.ImageURL = defaultImages[rand.Intn(len(defaultImages))]
		}

		var recipe models.Recipe
		if err := db.First(&recipe, id).Error; err != nil {
			http.NotFound(w, r)
			return
		}

		recipe.Title = updatedRecipe.Title
		recipe.ImageURL = updatedRecipe.ImageURL
		recipe.RecipeURL = updatedRecipe.RecipeURL
		recipe.ShareURL = updatedRecipe.ShareURL
		recipe.Email = updatedRecipe.Email

		db.Save(&recipe)

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(recipe)
		if err != nil {
			return
		}
	}
}
