package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

var db *sql.DB

func main() {
	// Підключення до SQLite
	var err error
	db, err = sql.Open("sqlite", "./database/app.db")
	if err != nil {
		log.Fatal("Помилка підключення до бази даних:", err)
	}
	defer db.Close()

	// Створення таблиці
	createTable()

	// Статичні файли
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Маршрути
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/dashboard", dashboardHandler)

	// Запуск сервера
	fmt.Println("Сервер запущено на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println("Помилка при завантаженні шаблону:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			log.Println("Помилка при завантаженні шаблону:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			sendError(w, "Невірний формат даних")
			return
		}

		// Перевірка даних
		if user.Username == "" || user.Password == "" {
			sendError(w, "Усі поля обов'язкові")
			return
		}

		// Пошук користувача в БД
		var storedPassword string
		err = db.QueryRow("SELECT password FROM users WHERE username = ?", user.Username).Scan(&storedPassword)
		if err != nil {
			sendError(w, "Користувач не знайдений")
			return
		}

		// Перевірка пароля
		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
		if err != nil {
			sendError(w, "Невірний пароль")
			return
		}

		// Встановлення сесії (спрощено)
		http.SetCookie(w, &http.Cookie{
			Name:  "username",
			Value: user.Username,
		})

		json.NewEncoder(w).Encode(Response{Success: true})
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Видалення сесії (спрощено)
	http.SetCookie(w, &http.Cookie{
		Name:   "username",
		Value:  "",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не підтримується", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		sendError(w, "Невірний формат даних")
		return
	}

	// Перевірка даних
	if user.Username == "" || user.Email == "" || user.Password == "" {
		sendError(w, "Усі поля обов'язкові")
		return
	}

	// Хешування пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		sendError(w, "Помилка при хешуванні пароля")
		return
	}

	// Збереження в БД
	stmt, err := db.Prepare("INSERT INTO users(username, email, password) VALUES(?, ?, ?)")
	if err != nil {
		sendError(w, "Помилка бази даних")
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Email, string(hashedPassword))
	if err != nil {
		sendError(w, "Користувач вже існує")
		return
	}

	json.NewEncoder(w).Encode(Response{Success: true})
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Перевірка сесії (спрощено)
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Отримання даних користувача з БД
	var user User
	err = db.QueryRow("SELECT username, email FROM users WHERE username = ?", cookie.Value).Scan(&user.Username, &user.Email)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Відображення сторінки
	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		log.Println("Помилка при завантаженні шаблону:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, user)
}

func createTable() {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		email TEXT UNIQUE,
		password TEXT
	);`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("Помилка при створенні таблиці:", err)
	}
}

func sendError(w http.ResponseWriter, message string) {
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error:   message,
	})
}
