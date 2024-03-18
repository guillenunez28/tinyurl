package db

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() (*sql.DB, error) {
	var err error
	db, err := sql.Open("mysql", "root:mypass@tcp(localhost:3306)/")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	sqlStmt := `
	CREATE DATABASE IF NOT EXISTS TinyUrls;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}
	sqlStmt = `
	USE TinyUrls;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	// Create table for storing URLs
	sqlStmt = `
	create table if not exists urls (
		short_url varchar(255) not null primary key, 
		long_url varchar(255),
		expiration_date DATETIME
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	// Create table for storing URL hits over time
	sqlStmt = `
	create table if not exists urls_hits (
		short_url varchar(255) not null, 
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		INDEX index_short (short_url)
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return db, nil
}
