package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"multithreaded-webserver/server"
)

func main() {

	httpsPort := flag.Int("https-port", 8443, "HTTPS port to listen on")
	httpPort := flag.Int("http-port", 8080, "HTTP port for redirect")
	rateLimit := flag.Int("rate-limit", 5, "Requests per minute per IP")
	flag.Parse()

	server.RateLimit = *rateLimit
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Try to serve static file
		_, err := os.Stat("static" + r.URL.Path)
		if os.IsNotExist(err) && r.URL.Path != "/" {
			server.NotFoundHandler(w, r)
			return
		}
		fs.ServeHTTP(w, r)
	})
	mux.HandleFunc("/slow", server.SlowHandler)
	mux.HandleFunc("/stats", server.StatsHandler)
	mux.HandleFunc("/upload", server.FileUploadHandler)

	rateLimitedMux := server.RateLimitMiddleware(mux)
	loggedMux := server.LoggingMiddleware(rateLimitedMux)

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(*httpsPort),
		Handler: loggedMux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Server starting on https://localhost:%d ...\n", *httpsPort)
		if err := srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil && err != http.ErrServerClosed {
			fmt.Printf("ListenAndServeTLS error: %v\n", err)
		}
	}()

	go func() {
		http.ListenAndServe(":"+strconv.Itoa(*httpPort), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpsURL := "https://localhost:" + strconv.Itoa(*httpsPort) + r.RequestURI
			http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
		}))
	}()

	<-stop // Wait for signal

	fmt.Println("\nShutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	} else {
		fmt.Println("Server exited gracefully.")
	}

}
