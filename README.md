<p align="center">
  <img src="https://raw.githubusercontent.com/golang-samples/gopher-vector/master/gopher.svg" alt="Go Gopher" width="120"/>
</p>

<h1 align="center">Go Multithreaded Web Server 🚀</h1>

<p align="center">
  <b>A modern, secure, and concurrent web server built with Go</b><br>
  <i>Featuring HTTPS, real-time dashboard, file uploads, rate limiting, and more!</i>
</p>

---

## 🌟 Features

- ⚡ **Concurrent/Multithreaded**: Handles many clients at once using Go's goroutines
- 🔒 **HTTPS/TLS**: Secure connections out of the box
- 📊 **Real-Time Dashboard**: Modern UI to monitor stats and interact with the server
- 📁 **File Uploads**: Upload files directly from the dashboard
- 🚦 **Rate Limiting**: Prevents abuse by limiting requests per IP
- 🛑 **Custom Error Pages**: Friendly 404 and 500 error pages
- 🏁 **Graceful Shutdown**: Safely finishes requests before stopping
- ⚙️ **Configurable**: Set ports and limits via command-line flags
- 📝 **Easy to Extend**: Modular codebase for adding new features

---

## 🚀 Quick Start

1. **Clone the repo:**
   ```sh
   git clone https://github.com/yourusername/go-multithreaded-webserver.git
   cd go-multithreaded-webserver
   ```
2. **Initialize Go modules:**
   ```sh
   go mod tidy
   ```
3. **Generate self-signed certificates (for HTTPS):**
   ```sh
   openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/CN=localhost"
   ```
4. **Create uploads directory:**
   ```sh
   mkdir uploads
   ```
5. **Run the server:**
   ```sh
   go run . --https-port=8443 --http-port=8080 --rate-limit=100
   ```
6. **Open the dashboard:**
   - Visit [https://localhost:8443/dashboard.html](https://localhost:8443/dashboard.html) in your browser
   - Accept the self-signed certificate warning

---


## ⚙️ Configuration

You can configure the server using command-line flags:

| Flag           | Default | Description                        |
|----------------|---------|------------------------------------|
| --https-port   | 8443    | HTTPS port to listen on            |
| --http-port    | 8080    | HTTP port for redirect             |
| --rate-limit   | 5       | Requests per minute per IP         |

Example:
```sh
go run . --https-port=9443 --http-port=9080 --rate-limit=20
```

---

## 📦 Project Structure

```
multithreaded-webserver/
├── main.go                # Entry point
├── server/                # Handlers, middleware, logic
│   └── server.go
├── static/                # Static files (dashboard, error pages)
│   ├── dashboard.html
│   ├── 404.html
│   ├── 500.html
│   └── ...
├── uploads/               # Uploaded files
├── cert.pem, key.pem      # TLS certificates
├── go.mod
└── README.md
```

---

## 💡 How It Works

- **Concurrency:** Go's HTTP server handles each request in a separate goroutine.
- **Rate Limiting:** Middleware tracks requests per IP and blocks excessive requests.
- **Dashboard:** Uses AJAX to fetch real-time stats and interact with endpoints.
- **Graceful Shutdown:** Listens for OS signals and finishes in-progress requests before exiting.
- **HTTPS:** All traffic is encrypted for security.


---

## 🛠️ Extending the Project

- Add authentication middleware
- Connect to a database
- Add REST API endpoints
- Deploy with Docker

---




