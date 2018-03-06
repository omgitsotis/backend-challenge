package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/omgitsotis/backend-challenge/dao"
)

var logger = log.New(os.Stdout, "SQLite:", log.Ldate|log.Ltime|log.Lshortfile)

type SQLiteDAO struct {
	db *sql.DB
}

func (sql *SQLiteDAO) GetItems(term string, r dao.Radius) ([]dao.Row, error) {
	results := make([]dao.Row, 0)

	q := fmt.Sprintf(
		"SELECT * FROM items WHERE item_name LIKE '%%%s%%' AND lat > %f AND lat < %f AND lng > %f AND lng < %f COLLATE NOCASE",
		term,
		r.MinLatitude,
		r.MaxLatitude,
		r.MinLongitude,
		r.MaxLongitude,
	)
	rows, err := sql.db.Query(q)

	if err != nil {
		log.Printf("Error executing query: %s", err.Error())
		return results, err
	}

	defer rows.Close()

	for rows.Next() {
		var itemID, itemURL, imgURLs string
		var lat, lng float64
		if err := rows.Scan(&itemID, &lat, &lng, &itemURL, &imgURLs); err != nil {
			log.Printf("Error scaning row: %s", err.Error())
			return results, err
		}

		// Format the img urls field to an array of strings
		imgURLs = strings.Replace(imgURLs, "[", "", -1)
		imgURLs = strings.Replace(imgURLs, "]", "", -1)
		imgURLs = strings.Replace(imgURLs, "\"", "", -1)

		imgArr := strings.Split(imgURLs, ",")

		xCoor := math.Pow((r.CenterLatitude - lat), 2)
		yCoor := math.Pow((r.CenterLongitude - lng), 2)

		dis := math.Sqrt(xCoor + yCoor)

		row := dao.Row{itemID, lat, lng, itemURL, imgArr, dis}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error fetching row %s", err.Error())
		return results, err
	}

	return results, nil
}

func (sql *SQLiteDAO) CountItems(term string, r dao.Radius) (int, error) {
	var count int
	q := fmt.Sprintf(
		"SELECT COUNT(*) FROM items WHERE item_name LIKE '%%%s%%' AND lat > %f AND lat < %f AND lng > %f AND lng < %f COLLATE NOCASE",
		term,
		r.MinLatitude,
		r.MaxLatitude,
		r.MinLongitude,
		r.MaxLongitude,
	)

	log.Println(q)

	row := sql.db.QueryRow(q)

	err := row.Scan(&count)
	if err != nil {
		log.Printf("Error performing scan %s", err)
		return 0, err
	}

	return count, nil
}

func (sql *SQLiteDAO) CloseDB() {
	sql.db.Close()
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
