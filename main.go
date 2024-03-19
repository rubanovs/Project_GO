package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/gorilla/mux"
)

// Command представляет bash-скрипт команды
type Command struct {
	ID      int
	Content string
	Output  string // Добавляем поле для вывода команды
}

var db *sql.DB

func main() {
	initDB()

	r := mux.NewRouter()
	r.HandleFunc("/commands", createCommand).Methods("POST")
	r.HandleFunc("/commands", getCommands).Methods("GET")
	r.HandleFunc("/commands/{id}", getCommand).Methods("GET")
	r.HandleFunc("/commands/{id}/stop", stopCommand).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Сервер запущен на порту %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func initDB() {
	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=commands sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	createTable()
}

func createTable() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS commands (
			id SERIAL PRIMARY KEY,
			content TEXT,
			output TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func insertCommand(content string) error {
	_, err := db.Exec("INSERT INTO commands (content) VALUES ($1)", content)
	return err
}

func getCommands(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, content, output FROM commands")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var cmd Command
		err := rows.Scan(&cmd.ID, &cmd.Content, &cmd.Output)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		commands = append(commands, cmd)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commands)
}

func getCommand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var cmd Command
	err = db.QueryRow("SELECT id, content, output FROM commands WHERE id = $1", id).Scan(&cmd.ID, &cmd.Content, &cmd.Output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cmd)
}

func stopCommand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем команду по идентификатору
	var cmd Command
	err = db.QueryRow("SELECT id, content FROM commands WHERE id = $1", id).Scan(&cmd.ID, &cmd.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Запускаем команду и получаем ее вывод
	output, err := executeCommand(cmd.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Обновляем запись в базе данных с выводом команды
	_, err = db.Exec("UPDATE commands SET output = $1 WHERE id = $2", output, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный статус
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Команда успешно остановлена")
}

// Функция для выполнения команды и сохранения ее вывода
func executeCommand(content string) (string, error) {
	cmd := exec.Command("bash", "-c", content)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func createCommand(w http.ResponseWriter, r *http.Request) {
	// Читаем содержимое команды из тела запроса
	content := r.FormValue("content")

	// Вставляем команду в базу данных
	_, err := db.Exec("INSERT INTO commands (content) VALUES ($1)", content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Команда успешно создана")
}
