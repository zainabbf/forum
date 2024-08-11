package databaseAPI

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// HasUpvoted check if user has upvoted a post
func HasUpvoted(database *sql.DB, username string, postId int) bool {
	rows, _ := database.Query("SELECT vote FROM votes WHERE username = ? AND post_id = ? ", username, postId)
	vote := 0
	for rows.Next() {
		rows.Scan(&vote)
	}
	if vote == 1 {
		return true
	}
	return false
}

// HasDownvoted check if user has downvoted a post
func HasDownvoted(database *sql.DB, username string, postId int) bool {
	rows, _ := database.Query("SELECT vote FROM votes WHERE username = ? AND post_id = ? ", username, postId)
	vote := 0
	for rows.Next() {
		rows.Scan(&vote)
	}
	if vote == -1 {
		return true
	}
	return false
}

// RemoveVote removes a vote from a post
func RemoveVote(database *sql.DB, postId int, username string) {
	statement, _ := database.Prepare("DELETE FROM votes WHERE post_id = ? AND username = ?")
	statement.Exec(postId, username)
}

func HasCommentUpvoted(database *sql.DB, username string, commentId int) (bool) {
	rows, err := database.Query("SELECT vote FROM comment_votes WHERE username = ? AND comment_id = ?", username, commentId)
	if err != nil {
		return false
	}
	defer rows.Close()

	vote := 0
	for rows.Next() {
		if err := rows.Scan(&vote); err != nil {
			return false
		}
	}
	if vote==1{
		return true
	}else{
		return false
	}
    
}

func HasCommentDownvoted(database *sql.DB, username string, commentId int) (bool) {
	rows, err := database.Query("SELECT vote FROM comment_votes WHERE username = ? AND comment_id = ?", username, commentId)
	if err != nil {
		return false
	}
	defer rows.Close()

	vote := 0
	for rows.Next() {
		if err := rows.Scan(&vote); err != nil {
			return false
		}
	}
	if vote==-1{
		return true
	}else{
		return false
	}

}

func AddCommentVote(database *sql.DB, commentId int, username string, vote int) error {
	statement, err := database.Prepare("INSERT INTO comment_votes (username, comment_id, vote) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(username, commentId, vote)
	return err
}

// Remove a vote from the comment
func RemoveCommentVote(database *sql.DB, commentId int, username string) error {
	statement, err := database.Prepare("DELETE FROM comment_votes WHERE comment_id = ? AND username = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(commentId, username)
	return err
}

// DecreaseUpvotes decreases the upvotes of a post by 1
func DecreaseUpvotes(database *sql.DB, postId int) {
	statement, _ := database.Prepare("UPDATE posts SET upvotes = upvotes - 1 WHERE id = ?")
	statement.Exec(postId)
}

// DecreaseDownvotes decreases the downvotes of a post by 1
func DecreaseDownvotes(database *sql.DB, postId int) {
	statement, _ := database.Prepare("UPDATE posts SET downvotes = downvotes - 1 WHERE id = ?")
	statement.Exec(postId)
}

// IncreaseUpvotes increases the upvotes of a post by 1
func IncreaseUpvotes(database *sql.DB, postId int) {
	statement, _ := database.Prepare("UPDATE posts SET upvotes = upvotes + 1 WHERE id = ?")
	statement.Exec(postId)
}

// IncreaseDownvotes increases the downvotes of a post by 1
func IncreaseDownvotes(database *sql.DB, postId int) {
	statement, _ := database.Prepare("UPDATE posts SET downvotes = downvotes + 1 WHERE id = ?")
	statement.Exec(postId)
}

// AddVote adds a vote to the database
func AddVote(database *sql.DB, postId int, username string, vote int) {
	statement, _ := database.Prepare("INSERT INTO votes (username, post_id, vote) VALUES (?, ?, ?)")
	statement.Exec(username, postId, vote)
}

// UpdateVote updates the vote of a user for a postS
func UpdateVote(database *sql.DB, postId int, username string, vote int) {
	statement, _ := database.Prepare("UPDATE votes SET vote = ? WHERE post_id = ? AND username = ?")
	statement.Exec(vote, postId, username)
}

// IncreaseCommentUpvotes increases the upvotes of a comment by 1

func IncreaseCommentUpvotes(database *sql.DB, commentId int) error {
	statement, err := database.Prepare("UPDATE comments SET upvotes = upvotes + 1 WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(commentId)
	return err
}

// DecreaseCommentUpvotes decreases the upvotes of a comment by 1
func DecreaseCommentUpvotes(database *sql.DB, commentId int) error {
	statement, err := database.Prepare("UPDATE comments SET upvotes = upvotes - 1 WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(commentId)
	return err
}

// IncreaseCommentDownvotes increases the downvotes of a comment by 1
func IncreaseCommentDownvotes(database *sql.DB, commentId int) error {
	statement, err := database.Prepare("UPDATE comments SET downvotes = downvotes + 1 WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(commentId)
	return err
}

// DecreaseCommentDownvotes decreases the downvotes of a comment by 1
func DecreaseCommentDownvotes(database *sql.DB, commentId int) error {
	statement, err := database.Prepare("UPDATE comments SET downvotes = downvotes - 1 WHERE id = ? AND downvotes > 0")
	if err != nil {
		return err
	}
	_, err = statement.Exec(commentId)
	return err
}

// UpdateVote updates the vote of a user for a postS
func UpdateCommentVote(database *sql.DB, commentId int, username string, vote int) error {
	statement, err := database.Prepare("UPDATE comment_votes SET vote = ? WHERE comment_id = ? AND username = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(vote, commentId, username)
	return err
}
