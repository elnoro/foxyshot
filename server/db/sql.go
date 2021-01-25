package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const table = `foxyshot_screenshots`
const timeout = 1 * time.Second

func NewSqlDb(dsn string) (*SqlDb, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// TODO pass as parameters
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Println("Connected to:", version)

	insertStmt, err := db.Prepare("INSERT INTO " + table + " (path, description) VALUES (?, ?)")
	if err != nil {
		db.Close()

		return nil, err
	}
	deleteStmt, err := db.Prepare("DELETE FROM " + table + " WHERE id = ?")
	if err != nil {
		insertStmt.Close()
		db.Close()

		return nil, err
	}

	return &SqlDb{db: db, insert: insertStmt, delete: deleteStmt}, nil
}

type SqlDb struct {
	db     *sql.DB
	insert *sql.Stmt
	delete *sql.Stmt
}

func (d *SqlDb) Close() error {
	d.insert.Close()
	d.delete.Close()

	return d.db.Close()
}

func (d *SqlDb) Add(path, desc string) (string, error) {
	r, err := d.insert.Exec(path, desc)
	if err != nil {
		return "", err
	}
	insertId, err := r.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(insertId)), nil
}

func (d *SqlDb) Remove(id string) error {
	removeId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	_, err = d.delete.Exec(removeId)
	if err != nil {
		return err
	}
	return nil
}

func (d *SqlDb) FindByDesc(desc string) (records []*ImageRecord, err error) {
	return d.findBySQL("select id, path, description from "+table+" where description like ?", `%`+desc+`%`)
}

func (d *SqlDb) All() ([]*ImageRecord, error) {
	return d.findBySQL("select id, path, description from " + table)
}

func (d *SqlDb) findBySQL(sql string, args ...interface{}) (records []*ImageRecord, err error) {
	rows, err := d.db.Query(sql, args...)
	if err != nil {
		return records, err
	}

	defer rows.Close()
	for rows.Next() {
		ir := &ImageRecord{Data: &ImageData{}}
		err := rows.Scan(&ir.Id, &ir.Data.Path, &ir.Data.Desc)
		if err != nil {
			return records, err
		}
		records = append(records, ir)
	}
	err = rows.Err()

	return records, err

}
