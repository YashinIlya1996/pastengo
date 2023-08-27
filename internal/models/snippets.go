package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	var userId int
	stmt := `
	insert into snippets (title, content, expires) 
	values ($1, $2, now() + interbal '$3 days')
	returning ID;
	`
	err := m.DB.QueryRow(stmt, title, content, expires).Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `
	select id, title, content, created, expires 
	from snippets
	where expires > now() and id = $1;
	`
	var res *Snippet = &Snippet{}
	err := m.DB.QueryRow(stmt, id).Scan(&res.ID, &res.Title, &res.Content, &res.Created, &res.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return res, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `
	select id, title, content, created, expires
	from snippets
	where expires > now() order by created desc 
	limit 10
	`
	res := []*Snippet{}
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
