package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"books/core"
	"books/core/storage/repositories"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() (*gin.Engine, *core.Core) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	repo := repositories.NewBookStorageInMemoryRepository()
	appCore := core.NewCore(repo)

	controllers := NewControllers(appCore)
	controllers.RegisterRoutes(router)

	return router, appCore
}

func TestAddBook(t *testing.T) {
	router, _ := setupTestRouter()

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful add",
			requestBody: map[string]interface{}{
				"isbn":   "9783161484100",
				"title":  "Test Book",
				"author": "Test Author",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "missing title",
			requestBody: map[string]interface{}{
				"isbn":   "9783161484100",
				"author": "Test Author",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing author",
			requestBody: map[string]interface{}{
				"isbn":  "9783161484100",
				"title": "Test Book",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing isbn",
			requestBody: map[string]interface{}{
				"title":  "Test Book",
				"author": "Test Author",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid isbn format",
			requestBody: map[string]interface{}{
				"isbn":   "invalid",
				"title":  "Test Book",
				"author": "Test Author",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty body",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d. Body: %s", tc.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestGetAllBooks(t *testing.T) {
	router, appCore := setupTestRouter()

	appCore.AddBook(nil, "Test Book", "Test Author", "9783161484100")

	req, _ := http.NewRequest(http.MethodGet, "/books", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	books, ok := response["books"].([]interface{})
	if !ok {
		t.Errorf("expected books array in response")
		return
	}

	if len(books) != 1 {
		t.Errorf("expected 1 book, got %d", len(books))
	}
}

func TestGetBookByISBN(t *testing.T) {
	router, appCore := setupTestRouter()

	validISBN := "9783161484100"

	appCore.AddBook(nil, "Test Book", "Test Author", validISBN)

	tests := []struct {
		name           string
		isbn           string
		expectedStatus int
	}{
		{
			name:           "existing book",
			isbn:           validISBN,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "non-existing book",
			isbn:           "9780306406157",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/books/isbn/"+tc.isbn, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d. Body: %s", tc.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	router, appCore := setupTestRouter()

	validISBN := "9783161484100"

	appCore.AddBook(nil, "Test Book", "Test Author", validISBN)

	tests := []struct {
		name           string
		isbn           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "successful update",
			isbn: validISBN,
			requestBody: map[string]interface{}{
				"title":  "Updated Title",
				"author": "Updated Author",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "partial update - title only",
			isbn: validISBN,
			requestBody: map[string]interface{}{
				"title": "Another Title",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "non-existing book",
			isbn: "9780306406157",
			requestBody: map[string]interface{}{
				"title": "New Title",
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(http.MethodPut, "/books/"+tc.isbn, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d. Body: %s", tc.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestDeleteBook(t *testing.T) {
	router, appCore := setupTestRouter()

	validISBN := "9783161484100"

	appCore.AddBook(nil, "Test Book", "Test Author", validISBN)

	tests := []struct {
		name           string
		isbn           string
		expectedStatus int
	}{
		{
			name:           "successful delete",
			isbn:           validISBN,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "non-existing book",
			isbn:           "9780306406157",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, "/books/"+tc.isbn, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d. Body: %s", tc.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestHealthCheck(t *testing.T) {
	router, _ := setupTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "ok" {
		t.Errorf("expected status 'ok', got '%v'", response["status"])
	}
}
