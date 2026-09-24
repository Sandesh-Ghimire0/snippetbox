package mysql

import (
	"database/sql"

	"github.com/Sandesh-Ghimire0/snippetbox/pkg/models"
)

// contain the code specifically for working with the snippets in our MySQL database.

type SnippetModel struct {
	DB *sql.DB
}

// inserts the new data
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	return 0, nil
}

// return the snippet with specific id
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
} 

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
