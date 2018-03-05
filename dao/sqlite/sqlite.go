package sqlite

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/omgitsotis/backend-challenge/dao"
)

var logger = log.New(os.Stdout, "SQLite:", log.Ldate|log.Ltime|log.Lshortfile)

type SQLiteDAO struct {
	db *sql.DB
}

func (sql *SQLiteDAO) GetItemsByTerm(term string) ([]dao.Row, error) {
	// rows, err = sql.db.Query(
	//     "SELECT * FROM items WHERE item_name LIKE "%?%" COLLATE NOCASE",
	//     term
	// )
	//
	// defer rows.Close()
	// for rows.Next() {
	//
	// }
	return []dao.Row{}, nil
}

func (sql *SQLiteDAO) GetItemsByLocation(lang, long float64) ([]dao.Row, error) {
	return []dao.Row{}, nil
}

func (sql *SQLiteDAO) CloseDB() {
	sql.db.Close()
}

func (sql *SQLiteDAO) GetItemCount(term string, r dao.Radius) (int, error) {
	var count int
	row := sql.db.QueryRow(
		"SELECT COUNT(*) FROM items WHERE item_name LIKE '%?%' AND lat > ? AND lat < ? AND lng > ? AND < ?",
		term,
		r.MinLatitude,
		r.MaxLatitude,
		r.MinLongitude,
		r.MaxLongitude,
	)

	err := row.Scan(&count)
	if err != nil {
		log.Print(err)
		return 0, err
	}

	return count, nil
}

func NewSQLiteDAO(file string) (*SQLiteDAO, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		logger.Printf("Error opening database: [%s]", err.Error())
		return nil, err
	}

	dao := SQLiteDAO{db}
	return &dao, nil
}
