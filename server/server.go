package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type Visitor struct {
	LastSeen time.Time
	Tokens   int
}

var (
	TotalRequests     int
	ActiveConnections int
	StatsMutex        sync.Mutex
	Visitors          = make(map[string]*Visitor)
	VisitorsMutex     sync.Mutex
	RateLimit         = 5
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		VisitorsMutex.Lock()
		v, exists := Visitors[ip]
		if !exists || time.Since(v.LastSeen) > time.Minute {
			Visitors[ip] = &Visitor{LastSeen: time.Now(), Tokens: RateLimit - 1}
			VisitorsMutex.Unlock()
			next.ServeHTTP(w, r)
			return
		}
		if v.Tokens > 0 {
			v.Tokens--
			v.LastSeen = time.Now()
			VisitorsMutex.Unlock()
			next.ServeHTTP(w, r)
			return
		}
		VisitorsMutex.Unlock()
		http.Error(w, "Rate limit exceeded. Try again later.", http.StatusTooManyRequests)
	})
}

func SlowHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Started slow handler")
	time.Sleep(5 * time.Second)
	fmt.Fprintln(w, "this is a slow response")
	fmt.Println("Finished slow handler")
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		StatsMutex.Lock()
		TotalRequests++
		ActiveConnections++
		StatsMutex.Unlock()

		fmt.Printf("Recieved request : %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)

		StatsMutex.Lock()
		ActiveConnections--
		StatsMutex.Unlock()
	})
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	StatsMutex.Lock()
	defer StatsMutex.Unlock()
	fmt.Fprintf(w, "total Requests : %d\n Active connections : %d\n", TotalRequests, ActiveConnections)
}

// NotFoundHandler serves a custom 404 page
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	http.ServeFile(w, r, "static/404.html")
}

// FileUploadHandler handles file uploads via POST
func FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
			<form enctype="multipart/form-data" action="/upload" method="post">
				<input type="file" name="myfile" />
				<input type="submit" value="Upload" />
			</form>
		`)
		return
	}

	file, header, err := r.FormFile("myfile")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "static/500.html")
		return
	}
	defer file.Close()

	out, err := os.Create("uploads/" + header.Filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "static/500.html")
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "static/500.html")
		return
	}

	fmt.Fprintf(w, "File %s uploaded successfully!", header.Filename)
}
