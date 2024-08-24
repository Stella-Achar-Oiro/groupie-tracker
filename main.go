package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "context"
    "time"

    "groupie-tracker/handlers"
)

func withLogging(h http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        h.ServeHTTP(w, r)
    }
}

func main() {
    // Serve static files
    fileServer := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fileServer))

    // Routes with middleware
    http.HandleFunc("/", withLogging(handlers.IndexHandler))
    http.HandleFunc("/home", withLogging(handlers.HomeHandler))
    http.HandleFunc("/artists/", withLogging(handlers.ArtistHandler))
    http.HandleFunc("/filter", withLogging(handlers.FilterHandler))

    // Start server
    srv := &http.Server{
        Addr: ":8010",
        Handler: nil,  // Default ServeMux is used
    }

    // Start the server in a goroutine
    go func() {
        fmt.Println("Server started at http://localhost:8010/")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    // Wait for interrupt signal to gracefully shutdown the server
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    <-c

    // Gracefully shutdown the server
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    fmt.Println("Server exiting")
}
