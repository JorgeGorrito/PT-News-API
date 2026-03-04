package mappers

import (
	"database/sql"
	"time"

	"github.com/JorgeGorrito/PT-News-API/internal/domain/entities"
)

func ScanAuthor(rows *sql.Rows) (*entities.Author, error) {
	var id int64
	var name, email, biography string
	var createdAt time.Time

	if err := rows.Scan(&id, &name, &email, &biography, &createdAt); err != nil {
		return nil, err
	}

	author, err := entities.NewAuthor(name, email)
	if err != nil {
		return nil, err
	}

	if err := author.SetBiography(biography); err != nil {
		return nil, err
	}

	author.SetID(id)

	return author, nil
}

func ScanAuthorRow(row *sql.Row) (*entities.Author, error) {
	var id int64
	var name, email, biography string
	var createdAt time.Time

	if err := row.Scan(&id, &name, &email, &biography, &createdAt); err != nil {
		return nil, err
	}

	author, err := entities.NewAuthor(name, email)
	if err != nil {
		return nil, err
	}

	if err := author.SetBiography(biography); err != nil {
		return nil, err
	}

	author.SetID(id)

	return author, nil
}
