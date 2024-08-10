package main

import (
	"FORUM-GO/databaseAPI"
	"FORUM-GO/webAPI"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	Id         int
	Username   string
	Title      string
	Categories []string
	Content    string
	CreatedAt  string
	UpVotes    int
	DownVotes  int
	Comments   []Comment
}

type Comment struct {
	Id        int
	PostId    int
	Username  string
	Content   string
	CreatedAt string
	UpVotes   int
	DownVotes int
}

// Database
var database *sql.DB

func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles(fmt.Sprintf("public/html/%s.html", tmpl))
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func badRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	renderTemplate(w, "400")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	renderTemplate(w, "404")
}

func internalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	renderTemplate(w, "500")
}
func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	renderTemplate(w, "405")
}

func handleUpvote(w http.ResponseWriter, r *http.Request) {
	// Handle upvote logic
	response := map[string]string{"status": "success"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Middleware to recover from panics and handle errors
func errorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				internalServerErrorHandler(w, r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Middleware to handle not found and bad request errors
func statusCodeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a ResponseRecorder to capture the status code
		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rr, r)

		switch rr.statusCode {
		case http.StatusMethodNotAllowed:
			methodNotAllowedHandler(w, r)
		}
	})
}

// ResponseRecorder to capture the status code
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func main() {
	// check if DB exists
	var _, err = os.Stat("database.db")

	// create DB if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create("database.db")
		if err != nil {
			return
		}
		defer file.Close()
	}

	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		fmt.Println("Timeout error:", err)
	}

	database, _ = sql.Open("sqlite3", "./database.db")

	databaseAPI.CreateUsersTable(database)
	databaseAPI.CreatePostTable(database)
	databaseAPI.CreateCommentTable(database)
	databaseAPI.CreateVoteTable(database)
	databaseAPI.CreateCategoriesTable(database)
	databaseAPI.CreateCategories(database)
	databaseAPI.CreateCategoriesIcons(database)

	webAPI.SetDatabase(database)

	fs := http.FileServer(http.Dir("public"))
	router := http.NewServeMux()
	fmt.Println("Starting server on port 8080")

	// Apply the middleware
	handler := statusCodeMiddleware(errorHandlingMiddleware(router))
	router.HandleFunc("/", webAPI.Index)
	router.HandleFunc("/register", webAPI.Register)
	router.HandleFunc("/login", webAPI.Login)
	router.HandleFunc("/post", webAPI.DisplayPost)
	router.HandleFunc("/filter", webAPI.GetPostsByApi)
	router.HandleFunc("/newpost", webAPI.NewPost)
	router.HandleFunc("/api/register", webAPI.RegisterApi)
	router.HandleFunc("/api/login", webAPI.LoginApi)
	router.HandleFunc("/api/logout", webAPI.LogoutAPI)
	router.HandleFunc("/api/createpost", webAPI.CreatePostApi)
	router.HandleFunc("/api/comments", webAPI.CommentsApi)
	router.HandleFunc("/api/vote", webAPI.VoteApi)
	router.HandleFunc("/api/comments/votes", webAPI.CommentVoteHandler)

	http.HandleFunc("/api/comments/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleUpvote(w, r)
		} else {
			methodNotAllowedHandler(w, r)
		}
	})

	router.Handle("/public/", http.StripPrefix("/public/", fs))

	// Custom error handlers
	router.HandleFunc("/400", badRequestHandler)
	router.HandleFunc("/404", notFoundHandler)
	router.HandleFunc("/500", internalServerErrorHandler)
	router.HandleFunc("/405", methodNotAllowedHandler)

	// Start the HTTP server
	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
