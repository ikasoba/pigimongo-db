package core

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/ikasoba/pigimongo-db/query"
	"github.com/rs/xid"

	_ "modernc.org/sqlite"
)

type Database struct {
	sql_db *sql.DB
}

/**
 * @params path - `:memory:` であればインメモリに保存
 */
func NewDatabase(path string) (*Database, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	db.Exec("CREATE TABLE IF NOT EXISTS pigimongo (id text, data text);")

	return &Database{
		sql_db: db,
	}, nil
}

func (db *Database) Add(value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	id := xid.New()
	_, err = db.sql_db.Exec("INSERT INTO pigimongo VALUES (?, json_set(?, '$.Id_', ?));", id.String(), string(data), id.String())
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) Update(in any, pigimongo_query string, values ...any) error {
	data, err := json.Marshal(in)
	if err != nil {
		return err
	}

	sql_query := "UPDATE pigimongo SET data = json_set(data, '$', ?) WHERE "

	tree, err := query.ParseQuery(pigimongo_query)
	if err != nil {
		return err
	}

	ctx := query.NewBuildContext(values...)
	err = ctx.BuildQueryToWhere(tree)
	if err != nil {
		return err
	}

	sql_query += ctx.Query

	ctx.Values = append([]any{string(data)}, ctx.Values...)

	_, err = db.sql_db.Exec(sql_query, ctx.Values...)

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) Find(out any, pigimongo_query string, values ...any) error {
	sql_query := "SELECT id, data FROM pigimongo WHERE "

	tree, err := query.ParseQuery(pigimongo_query)
	if err != nil {
		return err
	}

	ctx := query.NewBuildContext(values...)
	err = ctx.BuildQueryToWhere(tree)
	if err != nil {
		return err
	}

	sql_query += ctx.Query

	rows := db.sql_db.QueryRow(sql_query, ctx.Values...)

	res := []any{nil, nil}
	if err := rows.Scan(&res[0], &res[1]); err != nil {
		return err
	}

	if s, ok := res[1].(string); ok {
		if err := json.Unmarshal([]byte(s), out); err != nil {
			return err
		}

		return nil
	} else {
		return errors.New("query response is not string.")
	}
}

func (db *Database) Remove(pigimongo_query string, values ...any) error {
	sql_query := "DELETE FROM pigimongo WHERE "

	tree, err := query.ParseQuery(pigimongo_query)
	if err != nil {
		return err
	}

	ctx := query.NewBuildContext(values...)
	err = ctx.BuildQueryToWhere(tree)
	if err != nil {
		return err
	}

	sql_query += ctx.Query

	rows, err := db.sql_db.Exec(sql_query, ctx.Values...)

	if err != nil {
		return err
	}

	count, err := rows.RowsAffected()
	if count <= 0 {
		return errors.New("No matching items were found.")
	} else if err != nil {
		return err
	}

	return nil
}
