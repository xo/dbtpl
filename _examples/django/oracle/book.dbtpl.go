package oracle

// Code generated by dbtpl. DO NOT EDIT.

import (
	"context"
	"database/sql"
	"time"
)

// Book represents a row from 'django.books'.
type Book struct {
	BookID            int64          `json:"book_id"`              // book_id
	ISBN              sql.NullString `json:"isbn"`                 // isbn
	BookType          int64          `json:"book_type"`            // book_type
	Title             sql.NullString `json:"title"`                // title
	Year              int64          `json:"year"`                 // year
	Available         time.Time      `json:"available"`            // available
	BooksAuthorIDFkey int64          `json:"books_author_id_fkey"` // books_author_id_fkey
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the [Book] exists in the database.
func (b *Book) Exists() bool {
	return b._exists
}

// Deleted returns true when the [Book] has been marked for deletion
// from the database.
func (b *Book) Deleted() bool {
	return b._deleted
}

// Insert inserts the [Book] to the database.
func (b *Book) Insert(ctx context.Context, db DB) error {
	switch {
	case b._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case b._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO django.books (` +
		`isbn, book_type, title, year, available, books_author_id_fkey` +
		`) VALUES (` +
		`:1, :2, :3, :4, :5, :6` +
		`) RETURNING book_id INTO :7`
	// run
	logf(sqlstr, b.ISBN, b.BookType, b.Title, b.Year, b.Available, b.BooksAuthorIDFkey)
	var id int64
	if _, err := db.ExecContext(ctx, sqlstr, b.ISBN, b.BookType, b.Title, b.Year, b.Available, b.BooksAuthorIDFkey, sql.Out{Dest: &id}); err != nil {
		return logerror(err)
	} // set primary key
	b.BookID = int64(id)
	// set exists
	b._exists = true
	return nil
}

// Update updates a [Book] in the database.
func (b *Book) Update(ctx context.Context, db DB) error {
	switch {
	case !b._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case b._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE django.books SET ` +
		`isbn = :1, book_type = :2, title = :3, year = :4, available = :5, books_author_id_fkey = :6 ` +
		`WHERE book_id = :7`
	// run
	logf(sqlstr, b.ISBN, b.BookType, b.Title, b.Year, b.Available, b.BooksAuthorIDFkey, b.BookID)
	if _, err := db.ExecContext(ctx, sqlstr, b.ISBN, b.BookType, b.Title, b.Year, b.Available, b.BooksAuthorIDFkey, b.BookID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the [Book] to the database.
func (b *Book) Save(ctx context.Context, db DB) error {
	if b.Exists() {
		return b.Update(ctx, db)
	}
	return b.Insert(ctx, db)
}

// Upsert performs an upsert for [Book].
func (b *Book) Upsert(ctx context.Context, db DB) error {
	switch {
	case b._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `MERGE django.bookst ` +
		`USING (` +
		`SELECT :1 book_id, :2 isbn, :3 book_type, :4 title, :5 year, :6 available, :7 books_author_id_fkey ` +
		`FROM DUAL ) s ` +
		`ON s.book_id = t.book_id ` +
		`WHEN MATCHED THEN ` +
		`UPDATE SET ` +
		`t.isbn = s.isbn, t.book_type = s.book_type, t.title = s.title, t.year = s.year, t.available = s.available, t.books_author_id_fkey = s.books_author_id_fkey ` +
		`WHEN NOT MATCHED THEN ` +
		`INSERT (` +
		`isbn, book_type, title, year, available, books_author_id_fkey` +
		`) VALUES (` +
		`s.isbn, s.book_type, s.title, s.year, s.available, s.books_author_id_fkey` +
		`);`
	// run
	logf(sqlstr, b.BookID, b.ISBN, b.BookType, b.Title, b.Year, b.Available, b.BooksAuthorIDFkey)
	if _, err := db.ExecContext(ctx, sqlstr, b.BookID, b.ISBN, b.BookType, b.Title, b.Year, b.Available, b.BooksAuthorIDFkey); err != nil {
		return logerror(err)
	}
	// set exists
	b._exists = true
	return nil
}

// Delete deletes the [Book] from the database.
func (b *Book) Delete(ctx context.Context, db DB) error {
	switch {
	case !b._exists: // doesn't exist
		return nil
	case b._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM django.books ` +
		`WHERE book_id = :1`
	// run
	logf(sqlstr, b.BookID)
	if _, err := db.ExecContext(ctx, sqlstr, b.BookID); err != nil {
		return logerror(err)
	}
	// set deleted
	b._deleted = true
	return nil
}

// BookByBookID retrieves a row from 'django.books' as a [Book].
//
// Generated from index 'books_book_id_idx'.
func BookByBookID(ctx context.Context, db DB, bookID int64) (*Book, error) {
	// query
	const sqlstr = `SELECT ` +
		`book_id, isbn, book_type, title, year, available, books_author_id_fkey ` +
		`FROM django.books ` +
		`WHERE book_id = :1`
	// run
	logf(sqlstr, bookID)
	b := Book{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, bookID).Scan(&b.BookID, &b.ISBN, &b.BookType, &b.Title, &b.Year, &b.Available, &b.BooksAuthorIDFkey); err != nil {
		return nil, logerror(err)
	}
	return &b, nil
}

// BooksByBooksAuthorIDFkey retrieves a row from 'django.books' as a [Book].
//
// Generated from index 'books_books_auth_73ac0c26'.
func BooksByBooksAuthorIDFkey(ctx context.Context, db DB, booksAuthorIDFkey int64) ([]*Book, error) {
	// query
	const sqlstr = `SELECT ` +
		`book_id, isbn, book_type, title, year, available, books_author_id_fkey ` +
		`FROM django.books ` +
		`WHERE books_author_id_fkey = :1`
	// run
	logf(sqlstr, booksAuthorIDFkey)
	rows, err := db.QueryContext(ctx, sqlstr, booksAuthorIDFkey)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// process
	var res []*Book
	for rows.Next() {
		b := Book{
			_exists: true,
		}
		// scan
		if err := rows.Scan(&b.BookID, &b.ISBN, &b.BookType, &b.Title, &b.Year, &b.Available, &b.BooksAuthorIDFkey); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &b)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// Author returns the Author associated with the [Book]'s (BooksAuthorIDFkey).
//
// Generated from foreign key 'books_books_aut_73ac0c26_f'.
func (b *Book) Author(ctx context.Context, db DB) (*Author, error) {
	return AuthorByAuthorID(ctx, db, b.BooksAuthorIDFkey)
}
