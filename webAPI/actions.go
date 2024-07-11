package webAPI

import (
	"FORUM-GO/databaseAPI"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Vote struct {
	PostId int
	Vote   int
}

func CreatePostApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	cookie, _ := r.Cookie("SESSION")
	username := databaseAPI.GetUser(database, cookie.Value)
	title := r.FormValue("title")
	content := r.FormValue("content")
	categories := r.Form["categories[]"]

	// Check if content is empty
	if content == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Content cannot be empty"))
		return
	}
	// Check if content is empty
	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("title cannot be empty"))
		return
	}

	// If no categories are chosen, add the post to the "other" category
	if len(categories) == 0 {
		categories = append(categories, "Other")
	}

	//print the current working directory using
	wd, _ := os.Getwd()
	fmt.Println("Current working directory:", wd)

	validCategories := databaseAPI.GetCategories(database)
	for _, category := range categories {
		if !inArray1(category, validCategories) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid category : " + category))
			return
		}
	}
	stringCategories := strings.Join(categories, ",")
	now := time.Now()
	databaseAPI.CreatePost(database, username, title, stringCategories, content, now) // Update this line to save the image URL
	fmt.Println("Post created by " + username + " with title " + title + " at " + now.Format("2006-01-02 15:04:05"))
	http.Redirect(w, r, "/filter?by=myposts", http.StatusFound)
}

func inArray1(input string, array []string) bool {
	for _, v := range array {
		if v == input {
			return true
		}
	}
	return false
}

// CommentsApi creates a comment
func CommentsApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	cookie, _ := r.Cookie("SESSION")
	username := databaseAPI.GetUser(database, cookie.Value)
	postId := r.FormValue("postId")
	content := r.FormValue("content")
	now := time.Now()
	postIdInt, _ := strconv.Atoi(postId)
	databaseAPI.AddComment(database, username, postIdInt, content, now)
	fmt.Println("Comment created by " + username + " on post " + postId + " at " + now.Format("2006-01-02 15:04:05"))
	http.Redirect(w, r, "/post?id="+postId, http.StatusFound)
}

// VoteApi handles the voting on a post via API
func VoteApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !isLoggedIn(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cookie, err := r.Cookie("SESSION")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username := databaseAPI.GetUser(database, cookie.Value)
	postIdStr := r.FormValue("postId")
	voteStr := r.FormValue("vote")

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	vote, err := strconv.Atoi(voteStr)
	if err != nil || (vote != 1 && vote != -1) {
		http.Error(w, "Invalid vote value", http.StatusBadRequest)
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	switch vote {
	case 1:
		handleUpvote(w, username, postId, now)
	case -1:
		handleDownvote(w, username, postId, now)
	}
}

func handleUpvote(w http.ResponseWriter, username string, postId int, now string) {
	if databaseAPI.HasUpvoted(database, username, postId) {
		databaseAPI.RemoveVote(database, postId, username)
		databaseAPI.DecreaseUpvotes(database, postId)
		fmt.Printf("Removed upvote from %s on post %d at %s\n", username, postId, now)
		w.WriteHeader(http.StatusOK)
		return
	}

	if databaseAPI.HasDownvoted(database, username, postId) {
		databaseAPI.DecreaseDownvotes(database, postId)
		databaseAPI.IncreaseUpvotes(database, postId)
		databaseAPI.UpdateVote(database, postId, username, 1)
		fmt.Printf("%s upvoted on post %d at %s\n", username, postId, now)
		w.WriteHeader(http.StatusOK)
		return
	}

	databaseAPI.IncreaseUpvotes(database, postId)
	databaseAPI.AddVote(database, postId, username, 1)
	fmt.Printf("%s upvoted on post %d at %s\n", username, postId, now)
	w.WriteHeader(http.StatusOK)
}

func handleDownvote(w http.ResponseWriter, username string, postId int, now string) {
	if databaseAPI.HasDownvoted(database, username, postId) {
		databaseAPI.RemoveVote(database, postId, username)
		databaseAPI.DecreaseDownvotes(database, postId)
		fmt.Printf("Removed downvote from %s on post %d at %s\n", username, postId, now)
		w.WriteHeader(http.StatusOK)
		return
	}

	if databaseAPI.HasUpvoted(database, username, postId) {
		databaseAPI.DecreaseUpvotes(database, postId)
		databaseAPI.IncreaseDownvotes(database, postId)
		databaseAPI.UpdateVote(database, postId, username, -1)
		fmt.Printf("%s downvoted on post %d at %s\n", username, postId, now)
		w.WriteHeader(http.StatusOK)
		return
	}

	databaseAPI.IncreaseDownvotes(database, postId)
	databaseAPI.AddVote(database, postId, username, -1)
	fmt.Printf("%s downvoted on post %d at %s\n", username, postId, now)
	w.WriteHeader(http.StatusOK)
}
