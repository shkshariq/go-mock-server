package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type ApiResponse struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       json.RawMessage   `json:"body"`
}

var mockAPIs = struct {
	sync.RWMutex
	m map[string]ApiResponse
}{m: make(map[string]ApiResponse)}

func initDB() (*sql.DB, error) {
	// Replace 'username', 'password', and 'database' with your MySQL credentials
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/database")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func saveAPI(db *sql.DB, apiPath string, apiResponse ApiResponse) error {
	headersJSON, err := json.Marshal(apiResponse.Headers)
	if err != nil {
		fmt.Println("ERROR: saveAPI:headersJSON", err)
		return err
	}

	bodyJSON, err := json.Marshal(apiResponse.Body)
	if err != nil {
		fmt.Println("ERROR: saveAPI:bodyJSON", err)
		return err
	}

	_, err = db.Exec(`
	INSERT INTO mock_apis (api_path, status_code, headers, body)
	VALUES (?, ?, ?, ?);
`, apiPath, apiResponse.StatusCode, string(headersJSON), string(bodyJSON))

	fmt.Println("ERROR: saveAPI:Exec", err)
	return err
}

func loadAPI(db *sql.DB, apiPath string) (ApiResponse, error) {
	var apiResponse ApiResponse
	var headersJSON, bodyJSON string

	err := db.QueryRow("SELECT status_code, headers, body FROM mock_apis WHERE api_path = ?", apiPath).Scan(&apiResponse.StatusCode, &headersJSON, &bodyJSON)
	if err != nil {
		return ApiResponse{}, err
	}

	err = json.Unmarshal([]byte(headersJSON), &apiResponse.Headers)
	if err != nil {
		return ApiResponse{}, err
	}

	err = json.Unmarshal([]byte(bodyJSON), &apiResponse.Body)
	if err != nil {
		return ApiResponse{}, err
	}

	return apiResponse, nil
}

func handleCreateUpdateAPI(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var apiResponse ApiResponse
	err := json.NewDecoder(r.Body).Decode(&apiResponse)
	if err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	apiPath := r.URL.Query().Get("api_path")
	if apiPath == "" {
		http.Error(w, "api_path parameter is required", http.StatusBadRequest)
		return
	}

	mockAPIs.Lock()
	mockAPIs.m[apiPath] = apiResponse
	mockAPIs.Unlock()

	err = saveAPI(db, apiPath, apiResponse)
	if err != nil {
		http.Error(w, "Failed to save the API data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleMockAPI(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	apiPath := r.URL.Path
	if apiPath == "" {
		http.Error(w, "api_path parameter is required", http.StatusBadRequest)
		return
	}

	mockAPIs.RLock()
	apiResponse, err := loadAPI(db, apiPath)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "API not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to load API data", http.StatusInternalServerError)
		}
		return
	}
	mockAPIs.RUnlock()

	for key, value := range apiResponse.Headers {
		w.Header().Set(key, value)
	}

	w.WriteHeader(apiResponse.StatusCode)
	w.Write(apiResponse.Body)
}

func main() {

	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to initialize the database:", err)
	}
	defer db.Close()

	http.HandleFunc("/create_update_api", func(w http.ResponseWriter, r *http.Request) {
		handleCreateUpdateAPI(db, w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleMockAPI(db, w, r)
	})

	fmt.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
