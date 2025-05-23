package postgres

// Code generated by dbtpl. DO NOT EDIT.

import (
	"context"

	"github.com/lib/pq"
)

// AuthorBookResult is the result of a search.
type AuthorBookResult struct {
	AuthorID   int            `json:"author_id"`   // author_id
	AuthorName string         `json:"author_name"` // author_name
	BookID     int            `json:"book_id"`     // book_id
	BookISBN   string         `json:"book_isbn"`   // book_isbn
	BookTitle  string         `json:"book_title"`  // book_title
	BookTags   pq.StringArray `json:"book_tags"`   // book_tags
}

// AuthorBookResultsByTags runs a custom query, returning results as [AuthorBookResult].
func AuthorBookResultsByTags(ctx context.Context, db DB, tags pq.StringArray) ([]*AuthorBookResult, error) {
	// query
	const sqlstr = `SELECT ` +
		`a.author_id, ` + // ::integer AS author_id
		`a.name, ` + // ::text AS author_name
		`b.book_id, ` + // ::integer AS book_id
		`b.isbn, ` + // ::text AS book_isbn
		`b.title, ` + // ::text AS book_title
		`b.tags::text[] AS book_tags ` +
		`FROM books b ` +
		`JOIN authors a ON a.author_id = b.author_id ` +
		`WHERE b.tags && $1::varchar[]`
	// run
	logf(sqlstr, tags)
	rows, err := db.QueryContext(ctx, sqlstr, tags)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// load results
	var res []*AuthorBookResult
	for rows.Next() {
		var abr AuthorBookResult
		// scan
		if err := rows.Scan(&abr.AuthorID, &abr.AuthorName, &abr.BookID, &abr.BookISBN, &abr.BookTitle, &abr.BookTags); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &abr)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}
