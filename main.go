package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

var staticFiles embed.FS

type ServerInfo struct {
	StartTime time.Time
	BuildInfo *debug.BuildInfo
	Version   string
}

var serverInfo ServerInfo

func main() {
	serverInfo = ServerInfo{
		StartTime: time.Now(),
		Version:   "1.0.0",
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		serverInfo.BuildInfo = info
	}

	logFile := setupLogging()
	if logFile != nil {
		defer logFile.Close()
	}

	staticFS, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		log.Fatalf("Failed to get static files: %v", err)
	}

	fileServer := http.FileServer(http.FS(staticFS))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		log.Printf("%s - %s %s %s", r.RemoteAddr, r.Method, r.URL.Path, r.UserAgent())

		fileServer.ServeHTTP(w, r)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		info := map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now(),
			"uptime":    time.Since(serverInfo.StartTime).String(),
			"version":   serverInfo.Version,
		}
		json.NewEncoder(w).Encode(info)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("Server starting on http://localhost:%s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func setupLogging() *os.File {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("Failed to create log directory: %v", err)
		return nil
	}

	timestamp := time.Now().Format("2006-01-02")
	logPath := filepath.Join(logDir, fmt.Sprintf("server-%s.log", timestamp))

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return nil
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("Server logging initialized. Log file: %s", logPath)

	return logFile
}
