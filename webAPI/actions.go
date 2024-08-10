package webAPI

import (
	"FORUM-GO/databaseAPI"
	"fmt"
	"log"
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

	// Check if content is empty or have a white space
	if strings.TrimSpace(content) == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Content cannot be empty or just whitespace"))
		return
	}
	// Check if content is empty or have a white space
	if strings.TrimSpace(title) == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Title cannot be empty or just whitespace"))
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

func handleCommentUpvote(w http.ResponseWriter, username string, commentId int, now string) {
	// Check if the user has already upvoted the comment
	upvoted := databaseAPI.HasCommentUpvoted(database, username, commentId)
	if upvoted {
		if err := databaseAPI.RemoveCommentVote(database, commentId, username); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			fmt.Printf("Error removing comment vote: %v\n", err)
			return
		}

		if err := databaseAPI.DecreaseCommentUpvotes(database, commentId); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			fmt.Printf("Error decreasing comment upvotes: %v\n", err)
			return
		}

		fmt.Printf("Removed upvote from %s on comment %d at %s\n", username, commentId, now)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Check if the user has downvoted the comment
	downvoted := databaseAPI.HasCommentDownvoted(database, username, commentId)

	if downvoted {
		if err := databaseAPI.DecreaseCommentDownvotes(database, commentId); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			fmt.Printf("Error decreasing comment downvotes: %v\n", err)
			return
		}

		if err := databaseAPI.IncreaseCommentUpvotes(database, commentId); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			fmt.Printf("Error increasing comment upvotes: %v\n", err)
			return
		}

		if err := databaseAPI.UpdateCommentVote(database, commentId, username, 1); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			fmt.Printf("Error updating comment vote: %v\n", err)
			return
		}

		fmt.Printf("%s upvoted on comment %d at %s\n", username, commentId, now)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Otherwise, add an upvote
	if err := databaseAPI.IncreaseCommentUpvotes(database, commentId); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Printf("Error increasing comment upvotes: %v\n", err)
		return
	}

	if err := databaseAPI.AddCommentVote(database, commentId, username, 1); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Printf("Error adding comment vote: %v\n", err)
		return
	}

	fmt.Printf("%s upvoted on comment %d at %s\n", username, commentId, now)
	w.WriteHeader(http.StatusOK)
}
func handleCommentDownvote(w http.ResponseWriter, username string, commentId int, now string) {
	// Check if the user has already downvoted the comment
	downvoted := databaseAPI.HasCommentDownvoted(database, username, commentId)

	if downvoted {
		if err := databaseAPI.RemoveCommentVote(database, commentId, username); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			fmt.Printf("Error removing comment vote: %v\n", err)
			return
		}

		if err := databaseAPI.DecreaseCommentDownvotes(database, commentId); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			fmt.Printf("Error decreasing comment downvotes: %v\n", err)
			return
		}

		fmt.Printf("Removed downvote from %s on comment %d at %s\n", username, commentId, now)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Check if the user has upvoted the comment
	upvoted := databaseAPI.HasCommentUpvoted(database, username, commentId)

	if upvoted {
		if err := databaseAPI.DecreaseCommentUpvotes(database, commentId); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			fmt.Printf("Error decreasing comment upvotes: %v\n", err)
			return
		}

		if err := databaseAPI.IncreaseCommentDownvotes(database, commentId); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			fmt.Printf("Error increasing comment downvotes: %v\n", err)
			return
		}

		if err := databaseAPI.UpdateCommentVote(database, commentId, username, -1); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			fmt.Printf("Error updating comment vote: %v\n", err)
			return
		}

		fmt.Printf("%s downvoted on comment %d at %s\n", username, commentId, now)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Otherwise, add a downvote
	if err := databaseAPI.IncreaseCommentDownvotes(database, commentId); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Printf("Error increasing comment downvotes: %v\n", err)
		return
	}

	if err := databaseAPI.AddCommentVote(database, commentId, username, -1); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Printf("Error adding comment vote: %v\n", err)
		return
	}

	fmt.Printf("%s downvoted on comment %d at %s\n", username, commentId, now)
	w.WriteHeader(http.StatusOK)
}

// CommentVoteHandler handles the voting on a post via API
func CommentVoteHandler(w http.ResponseWriter, r *http.Request) {
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
	commentIdStr := r.FormValue("commentId")
	voteStr := r.FormValue("vote")
	log.Printf("comment %s: voteStr: %s", commentIdStr, voteStr)
	commentId, err := strconv.Atoi(commentIdStr)
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
		handleCommentUpvote(w, username, commentId, now)
	case -1:
		handleCommentDownvote(w, username, commentId, now)
	}
}

// func CommentVoteHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	commentIdStr := vars["id"]
// 	action := vars["action"]

// 	commentId, err := strconv.Atoi(commentIdStr)
// 	if err != nil {
// 		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
// 		return
// 	}

// 	cookie, err := r.Cookie("SESSION")
// 	if err != nil {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	username := databaseAPI.GetUser(database, cookie.Value)
// 	if username == "" {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}
// 	now := time.Now().Format("2006-01-02 15:04:05")
// 	switch action {
// 	case "upvote":
// 		handleCommentUpvote(w, username, commentId, now)
// 	case "downvote":
// 		handleCommentDownvote(w, username, commentId, now)

// 	default:
// 		http.Error(w, "Invalid action", http.StatusBadRequest)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }
