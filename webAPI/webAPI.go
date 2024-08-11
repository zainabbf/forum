package webAPI

import (
	"FORUM-GO/databaseAPI"
	"database/sql"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	IsLoggedIn bool
	Username   string
}

type HomePage struct {
	User              User
	Categories        []string
	Icons             []string
	PostsByCategories [][]databaseAPI.Post
}

type PostsPage struct {
	User  User
	Title string
	Posts []databaseAPI.Post
	Icon  string
}

type PostPage struct {
	User User
	Post databaseAPI.Post
}

var database *sql.DB

func SetDatabase(db *sql.DB) {
	database = db
}

func renderNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	t, _ := template.ParseFiles("public/HTML/404.html")
	t.Execute(w, nil)
}

func renderBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	t, _ := template.ParseFiles("public/HTML/400.html")
	t.Execute(w, nil)
}

// Index displays the Index page
func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderNotFound(w)
		return
	}
	if isLoggedIn(r) {
		cookie, _ := r.Cookie("SESSION")
		payload := HomePage{
			User:              User{IsLoggedIn: true, Username: databaseAPI.GetUser(database, cookie.Value)},
			Categories:        databaseAPI.GetCategories(database),
			Icons:             databaseAPI.GetCategoriesIcons(database),
			PostsByCategories: databaseAPI.GetPostsByCategories(database),
		}
		t, _ := template.ParseGlob("public/HTML/*.html")
		t.ExecuteTemplate(w, "forum.html", payload)
		return
	}
	payload := HomePage{
		User:              User{IsLoggedIn: false},
		Categories:        databaseAPI.GetCategories(database),
		Icons:             databaseAPI.GetCategoriesIcons(database),
		PostsByCategories: databaseAPI.GetPostsByCategories(database),
	}
	t, _ := template.ParseGlob("public/HTML/*.html")
	t.ExecuteTemplate(w, "forum.html", payload)
}

// DisplayPost displays a post on a template
func DisplayPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		renderBadRequest(w)
		return
	}
	post, err := databaseAPI.GetPost(database, id)
	if err != nil {
		if err == sql.ErrNoRows {
			renderNotFound(w)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	payload := PostPage{
		Post: post,
	}
	if !isLoggedIn(r) {
		payload.User = User{IsLoggedIn: false}
	} else {
		cookie, _ := r.Cookie("SESSION")
		username := databaseAPI.GetUser(database, cookie.Value)
		payload.User = User{IsLoggedIn: true, Username: username}
	}
	payload.Post.Comments = databaseAPI.GetComments(database, id)
	t, err := template.ParseGlob("public/HTML/*.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, "detail.html", payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// GetPostsByApi GetPostByApi gets all post filtered by the given parameters
func GetPostsByApi(w http.ResponseWriter, r *http.Request) {
	method := r.URL.Query().Get("by")
	if method == "category" {
		category := r.URL.Query().Get("category")
		posts := databaseAPI.GetPostsByCategory(database, category)
		payload := PostsPage{
			Title: "Posts in category " + category,
			Posts: posts,
			Icon:  databaseAPI.GetCategoryIcon(database, category),
		}
		if isLoggedIn(r) {
			payload.User = User{IsLoggedIn: true}
		}
		t, _ := template.ParseGlob("public/HTML/*.html")
		t.ExecuteTemplate(w, "posts.html", payload)
		return
	}
	if method == "myposts" {
		if isLoggedIn(r) {
			cookie, _ := r.Cookie("SESSION")
			username := databaseAPI.GetUser(database, cookie.Value)
			posts := databaseAPI.GetPostsByUser(database, username)
			payload := PostsPage{
				User:  User{IsLoggedIn: true},
				Title: "My posts",
				Posts: posts,
				Icon:  "fa-user",
			}
			t, _ := template.ParseGlob("public/HTML/*.html")
			t.ExecuteTemplate(w, "posts.html", payload)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if method == "liked" {
		if isLoggedIn(r) {
			cookie, _ := r.Cookie("SESSION")
			username := databaseAPI.GetUser(database, cookie.Value)
			posts := databaseAPI.GetLikedPosts(database, username)
			payload := PostsPage{
				User:  User{IsLoggedIn: true},
				Title: "Posts liked by me",
				Posts: posts,
				Icon:  "fa-heart",
			}
			t, _ := template.ParseGlob("public/HTML/*.html")
			t.ExecuteTemplate(w, "posts.html", payload)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// NewPost displays the NewPost page
func NewPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	t, _ := template.ParseGlob("public/HTML/*.html")
	t.ExecuteTemplate(w, "createThread.html", nil)
}
