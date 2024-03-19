package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCreateCommand(t *testing.T) {
	// Создаем временную базу данных для тестов
	db, err := sql.Open("postgres", "user=postgres dbname=commands_test sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Инициализируем роутер
	r := mux.NewRouter()
	r.HandleFunc("/commands", createCommand).Methods("POST")

	// Создаем тестовый запрос
	content := "echo 'Hello, World!'"
	reqBody := bytes.NewBufferString("content=" + content)
	req, err := http.NewRequest("POST", "/commands", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем тестовый ResponseRecorder
	rr := httptest.NewRecorder()

	// Выполняем запрос
	r.ServeHTTP(rr, req)

	// Проверяем код состояния
	if rr.Code != http.StatusCreated {
		t.Errorf("Ожидался код состояния %d, получен %d", http.StatusCreated, rr.Code)
	}

	// Проверяем, что команда была добавлена в базу данных
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM commands").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("Ожидалось 1 команда в базе данных, получено %d", count)
	}
}

func TestGetCommands(t *testing.T) {
	// Создаем временную базу данных для тестов
	db, err := sql.Open("postgres", "user=postgres dbname=commands_test sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Вставляем тестовые команды в базу данных
	_, err = db.Exec("INSERT INTO commands (content) VALUES ($1), ($2)", "echo 'Hello'", "ls -l")
	if err != nil {
		t.Fatal(err)
	}

	// Инициализируем роутер
	r := mux.NewRouter()
	r.HandleFunc("/commands", getCommands).Methods("GET")

	// Создаем тестовый запрос
	req, err := http.NewRequest("GET", "/commands", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем тестовый ResponseRecorder
	rr := httptest.NewRecorder()

	// Выполняем запрос
	r.ServeHTTP(rr, req)

	// Проверяем код состояния
	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался код состояния %d, получен %d", http.StatusOK, rr.Code)
	}

	// Проверяем, что полученные данные соответствуют ожидаемым
	var commands []Command
	err = json.NewDecoder(rr.Body).Decode(&commands)
	if err != nil {
		t.Fatal(err)
	}

	if len(commands) != 2 {
		t.Errorf("Ожидалось 2 команды, получено %d", len(commands))
	}
}

func TestGetCommand(t *testing.T) {
	// Создаем временную базу данных для тестов
	db, err := sql.Open("postgres", "user=postgres dbname=commands_test sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Вставляем тестовую команду в базу данных
	_, err = db.Exec("INSERT INTO commands (content) VALUES ($1)", "echo 'Hello'")
	if err != nil {
		t.Fatal(err)
	}

	// Инициализируем роутер
	r := mux.NewRouter()
	r.HandleFunc("/commands/{id}", getCommand).Methods("GET")

	// Создаем тестовый запрос
	req, err := http.NewRequest("GET", "/commands/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем тестовый ResponseRecorder
	rr := httptest.NewRecorder()

	// Выполняем запрос
	r.ServeHTTP(rr, req)

	// Проверяем код состояния
	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался код состояния %d, получен %d", http.StatusOK, rr.Code)
	}

	// Проверяем, что полученная команда соответствует ожидаемой
	var cmd Command
	err = json.NewDecoder(rr.Body).Decode(&cmd)
	if err != nil {
		t.Fatal(err)
	}

	if cmd.ID != 1 || cmd.Content != "echo 'Hello'" {
		t.Errorf("Полученная команда не соответствует ожидаемой")
	}
}

func TestStopCommand(t *testing.T) {
	// Создаем временную базу данных для тестов
	db, err := sql.Open("postgres", "user=postgres dbname=commands_test sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Вставляем тестовую команду в базу данных
	_, err = db.Exec("INSERT INTO commands (content) VALUES ($1)", "echo 'Hello'")
	if err != nil {
		t.Fatal(err)
	}

	// Инициализируем роутер
	r := mux.NewRouter()
	r.HandleFunc("/commands/{id}/stop", stopCommand).Methods("POST")

	// Создаем тестовый запрос
	req, err := http.NewRequest("POST", "/commands/1/stop", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем тестовый ResponseRecorder
	rr := httptest.NewRecorder()

	// Выполняем запрос
	r.ServeHTTP(rr, req)

	// Проверяем код состояния
	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался код состояния %d, получен %d", http.StatusOK, rr.Code)
	}

	// Здесь вы можете добавить дополнительные проверки, чтобы убедиться, что команда была остановлена

	// Пример: Проверяем, что в ответе есть сообщение об успешной остановке
	expectedResponse := "Команда успешно остановлена"
	if rr.Body.String() != expectedResponse {
		t.Errorf("Ожидалось сообщение '%s', получено '%s'", expectedResponse, rr.Body.String())
	}
}
