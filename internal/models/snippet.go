package models

import (
	"database/sql"
	"time"
)

// A Struct for the data in each snippet in our database,
// The fields of the struct corrosponde to the fields in our MySQL snippets
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Struct wrapping sql.DB method
type snippetModel struct {
	DB *sql.DB
}

// snippetModel method to insert a new snippet to the database
func (m *snippetModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

// snippetModel method to get a snippet base on ID
func (m *snippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

// snippetModel method to get the latest 10 snippets
func (m *snippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
