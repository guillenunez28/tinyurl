package tinyurl

import (
	"database/sql"
	"errors"
	"time"
)

type dbUrl struct {
	ShortUrl       string    `db:"short_url"`
	LongUrl        string    `db:"long_url"`
	ExpirationTime time.Time `db:"expiration_date"`
}

func createDbResource(db *sql.DB, tUrl TinyUrl) error {
	var args []any
	args = append(args, tUrl.ShortUrl)
	args = append(args, tUrl.LongUrl)
	var expirationTime sql.NullTime
	expirationTime.Valid = false // Set Valid to false to represent NULL
	stm := "INSERT INTO urls (short_url, long_url, expiration_date) VALUES (?, ?, ?)"

	if tUrl.ExpirationTime != nil {
		args = append(args, tUrl.ExpirationTime)
	} else {
		args = append(args, time.Time{})
	}

	_, err := db.Exec(stm, args...)
	if err != nil {
		return err
	}
	return nil
}

func deleteDbResource(db *sql.DB, shortUrl string) error {
	stm := "DELETE FROM urls WHERE short_url = ?"
	rows, err := db.Exec(stm, shortUrl)
	rowsAffects, err := rows.RowsAffected()
	if rowsAffects == 0 {
		return errors.New("No resource found")
	}
	if err != nil {
		return err
	}
	return nil
}

func getDbResource(db *sql.DB, shortUrl string) (TinyUrl, error) {
	var tinyResource dbUrl
	stm := "SELECT short_url, long_url, expiration_date FROM urls WHERE short_url = ?"

	row := db.QueryRow(stm, shortUrl)
	var expirationTime string
	err := row.Scan(&tinyResource.ShortUrl, &tinyResource.LongUrl, &expirationTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return TinyUrl{}, errors.New("Not Found")
		}
		return TinyUrl{}, err
	}
	if expirationTime != "0000-00-00 00:00:00" {
		timeExpired, err := time.Parse("2006-01-02 15:04:05", expirationTime)
		if err != nil {
			return TinyUrl{}, errors.New("error parsing time")
		}
		tinyResource.ExpirationTime = timeExpired
		if time.Now().After(tinyResource.ExpirationTime) {
			return TinyUrl{}, errors.New("Resource has expired.")
		}
	}
	apiTinyUrl := TinyUrl{
		ShortUrl:       tinyResource.ShortUrl,
		LongUrl:        tinyResource.LongUrl,
		ExpirationTime: &tinyResource.ExpirationTime,
	}
	return apiTinyUrl, nil
}

func addDbUrlHit(db *sql.DB, shortUrl string) error {
	stm := "INSERT INTO urls_hits (short_url) VALUES (?)"

	_, err := db.Exec(stm, shortUrl)
	if err != nil {
		return err
	}
	return nil
}

func getDbStats(db *sql.DB, stm, shortUrl string) (int, error) {
	var result int
	row := db.QueryRow(stm, shortUrl)
	err := row.Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("Not Found")
		}
		return 0, err
	}
	return result, nil
}

func getDbResourceStats(db *sql.DB, shortUrl string) (Stats, error) {
	stm := "SELECT COUNT(*) FROM urls_hits where short_url = ? AND created_at >= NOW() - INTERVAL 24 HOUR"
	last24Hours, err := getDbStats(db, stm, shortUrl)
	if err != nil {
		return Stats{}, err
	}

	stm = "SELECT COUNT(*) FROM urls_hits where short_url = ? AND created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)"
	last7Days, err := getDbStats(db, stm, shortUrl)
	if err != nil {
		return Stats{}, err
	}

	stm = "SELECT COUNT(*) FROM urls_hits where short_url = ?"
	allTime, err := getDbStats(db, stm, shortUrl)
	if err != nil {
		return Stats{}, err
	}
	response := Stats{
		HitsAllTime:     allTime,
		HitsLast24Hours: last24Hours,
		HitsLastWeek:    last7Days,
	}

	return response, nil
}

func getDbResources(db *sql.DB) ([]TinyUrl, error) {
	stm := "SELECT short_url, long_url, expiration_date FROM urls LIMIT 10"
	rows, err := db.Query(stm)
	if err != nil {
		return []TinyUrl{}, err
	}
	defer rows.Close()
	result := []TinyUrl{}
	for rows.Next() {
		var row TinyUrl
		var expiredTime string
		err := rows.Scan(&row.ShortUrl, &row.LongUrl, &expiredTime)
		if err != nil {
			return []TinyUrl{}, err
		}
		if expiredTime != "0000-00-00 00:00:00" {
			parseTimeExpired, err := time.Parse("2006-01-02 15:04:05", expiredTime)
			if err != nil {
				return []TinyUrl{}, err
			}
			row.ExpirationTime = &parseTimeExpired
		}
		result = append(result, row)
	}

	return result, nil
}
