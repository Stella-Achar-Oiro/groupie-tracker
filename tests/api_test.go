package tests

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "groupie-tracker/internal/routes"
)

func TestGetAllData(t *testing.T) {
    mux := routes.SetupRoutes()
    req, _ := http.NewRequest("GET", "/api/data", nil)
    rr := httptest.NewRecorder()

    mux.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }
}

func TestSearchArtists(t *testing.T) {
    mux := routes.SetupRoutes()
    req, _ := http.NewRequest("GET", "/api/search?q=test", nil)
    rr := httptest.NewRecorder()

    mux.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }
}