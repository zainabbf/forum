package databaseAPI

import (
	"database/sql"
	"fmt"

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

var DBPassword string // Database password

func SetDBPassword(password string) {
	DBPassword = password
}

// Connect to the database
func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:database.db?_auth&_auth_password=%s&_auth_crypt=sha256", DBPassword))
	if err != nil {
		return nil, err
	}
	return db, nil
}
