package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"groupie-tracker/internal/routes"
)

func TestGetAllData(t *testing.T) {
	done := make(chan bool) // Create a channel to signal test completion

	go func() {
		defer close(done) // Ensure the channel is closed when the goroutine finishes

		mux := routes.SetupRoutes()
		req, _ := http.NewRequest("GET", "/api/data", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		done <- true // Send a signal when the test is done
	}()

	<-done // Wait for the signal that the test is completed
}

func TestSearchArtists(t *testing.T) {
	done := make(chan bool) // Create a channel to signal test completion

	go func() {
		defer close(done) // Ensure the channel is closed when the goroutine finishes

		mux := routes.SetupRoutes()
		req, _ := http.NewRequest("GET", "/api/search?q=test", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		done <- true // Send a signal when the test is done
	}()

	<-done // Wait for the signal that the test is completed
}
