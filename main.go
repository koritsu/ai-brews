package main

import (
	"ai-brews/handlers/admin"
	"ai-brews/handlers/api"
	"errors"
	"github.com/pquerna/otp/totp"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"ai-brews/config"
	"ai-brews/migrations"
	"ai-brews/models"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config file" + err.Error())
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = cfg.DBPath
	}

	// Ensure the database directory exists
	dbDir := filepath.Dir(dbPath)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
			log.Fatal("Failed to create database directory")
		}
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	// Run migrations
	migrations.Migrate()

	// Create default admin if not exists
	var defaultAdmin models.Admin
	if err := db.Where("username = ?", cfg.AdminUsername).First(&defaultAdmin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			key, err := totp.Generate(totp.GenerateOpts{
				Issuer:      "AiBrews",
				AccountName: "admin",
			})
			if err != nil {
				log.Fatal(err)
			}

			defaultAdmin := models.Admin{
				Username: cfg.AdminUsername,
				Password: key.Secret(),
			}
			if err := db.Create(&defaultAdmin).Error; err != nil {
				log.Fatal("Failed to create default admin user")
			}
		} else {
			log.Fatal("Failed to check for default admin user")
		}
	}

	router := mux.NewRouter()

	fileServer := http.FileServer(http.Dir("./static"))
	// Serve static files
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", addCacheControl(fileServer)))

	// Serve HTML files
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	}).Methods("GET")

	router.HandleFunc("/about.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/about.html")
	}).Methods("GET")

	router.HandleFunc("/admin/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/admin_login.html")
	}).Methods("GET")

	// 어드민 권한 필요
	router.Handle("/admin/create", admin.AuthMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "templates/admin_create.html")
		}))).Methods("GET")

	router.Handle("/admin/list", admin.AuthMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "templates/admin_list.html")
		}))).Methods("GET")

	router.Handle("/admin/recipe/create", admin.AuthMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "templates/admin_recipe_create.html")
		}))).Methods("GET")

	router.Handle("/admin/recipe/update/{id}", admin.AuthMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "templates/admin_recipe_update.html")
		}))).Methods("GET")

	// API Endpoints
	router.HandleFunc("/recipes", api.Index(db)).Methods("GET")
	router.HandleFunc("/recipe/{id}", api.Detail(db)).Methods("GET")

	// 어드민 권한 필요 (나중에 회원 오픈 개별 생성 오픈 예정)
	router.Handle("/recipe", admin.AuthMiddleware(admin.CreateRecipe(db, cfg))).Methods("POST")
	router.Handle("/recipe/{id}", admin.AuthMiddleware(admin.UpdateRecipe(db))).Methods("PUT")

	// Admin Endpoints
	router.HandleFunc("/admin/login", admin.AdminLogin(db)).Methods("POST")
	router.HandleFunc("/admin/logout", admin.AdminLogout()).Methods("POST")

	// 어드민 권한 필요
	router.Handle("/admin/list", admin.AuthMiddleware(admin.ListAdmins(db))).Methods("GET")
	router.Handle("/admin", admin.AuthMiddleware(admin.CreateAdmin(db))).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.ServerPort
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
}

// 캐시 제어 헤더를 추가하는 미들웨어
func addCacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}
		next.ServeHTTP(w, r)
	})
}
