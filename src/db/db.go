package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(url string) (*Database, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Database{db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(query, args...)
}

func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(query, args...)
}

func (d *Database) Get(dest interface{}, query string, args ...interface{}) error {
	rows, err := d.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return sql.ErrNoRows
	}
	err = rows.Scan(dest)
	if err != nil {
		return err
	}
	return rows.Err()
}

func (d *Database) GetJSON(dest interface{}, query string, args ...interface{}) error {
	var jsonBytes []byte
	err := d.Get(&jsonBytes, query, args...)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, dest)
}

func (d *Database) Insert(table string, data map[string]interface{}) (int64, error) {
	columns := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	for column, value := range data {
		columns = append(columns, column)
		values = append(values, value)
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ", "), strings.Join(strings.Repeat("?", len(columns)), ", "))
	result, err := d.Exec(query, values...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (d *Database) Update(table string, data map[string]interface{}, where string, args ...interface{}) (int64, error) {
	columns := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	for column, value := range data {
		columns = append(columns, column)
		values = append(values, value)
	}
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(strings.Map(func(s string) string { return s + " = ?" }, columns), ", "), where)
	values = append(values, args...)
	result, err := d.Exec(query, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (d *Database) Delete(table string, where string, args ...interface{}) (int64, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)
	result, err := d.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
