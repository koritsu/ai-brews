package models

import (
	"gorm.io/gorm"
)

type Recipe struct {
	gorm.Model
	ID        int    `json:"id"`
	Title     string `json:"title"`
	ImageURL  string `json:"image_url"`
	RecipeURL string `json:"recipe_url"`
	ShareURL  string `json:"share_url"`
	Email     string `json:"email"`
}
