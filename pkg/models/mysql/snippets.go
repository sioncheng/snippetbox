package mysql

import (
	"database/sql"
	"errors"

	"github.com/sioncheng/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `insert into snippets (title, content, created, expires) 
		values (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), interval ? day) )
	`
	result, err := m.DB.Exec(stmt)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `select id, title, content, created, expires from snippets
		where expires > UTC_TIMESTAMP() and id = ?
	`
	row := m.DB.QueryRow(stmt, id)
	s := &models.Snippet{}
	err := row.Scan(&s.Id, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	} else {
		return s, nil
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {

	stmt := `select id, title, content, created, expires from snippets
		where expires > UTC_TIMESTAMP() order by created desc limit 10
	`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.Id, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
