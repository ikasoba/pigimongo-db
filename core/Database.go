package core

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"strings"

	_ "modernc.org/sqlite"
)

type Database struct {
	sql_db *sql.DB
}

type EqualPair struct {
	query string
	value any
}

/**
 * @params path - `:memory:` であればインメモリに保存
 */
func NewDatabase(path string) (*Database, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	db.Exec("CREATE TABLE IF NOT EXISTS pigimongo (data text);")

	return &Database{
		sql_db: db,
	}, nil
}

func (db *Database) Insert(value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = db.sql_db.Exec("INSERT INTO pigimongo VALUES (?);", string(data))
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) FindEquals(out any, items ...EqualPair) error {
	query := "SELECT data FROM pigimongo WHERE "
	expr := []string{}
	values := []any{}

	for _, item := range items {
		// jqっぽくなるように
		q, v := "$"+item.query, item.value

		data, err := json.Marshal(v)
		if err != nil {
			return err
		}

		expr = append(expr, "data -> ? = ?")
		values = append(values, q, string(data))
	}

	query += strings.Join(expr, " AND ")

	log.Print(query, values)

	rows := db.sql_db.QueryRow(query, values...)

	var res any
	if err := rows.Scan(&res); err != nil {
		return err
	}

	if s, ok := res.(string); ok {
		if err := json.Unmarshal([]byte(s), out); err != nil {
			return err
		}

		return nil
	} else {
		return errors.New("query response is not string.")
	}
}
