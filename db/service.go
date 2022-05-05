package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DBService struct {
	DB *sql.DB
}

func Init(password string, url string, dbName string) (*DBService, error) {
	dsn := "root:" + password + "@tcp(" + url + ")/" + dbName
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &DBService{DB: db}, nil
}

func (db *DBService) GetEventTraceHeight(chainId int) (int64, error) {
	sqlStr := ""
	row := db.DB.QueryRow(sqlStr, chainId)
	height := int64(0)
	err := row.Scan(&height)
	if err != nil {
		return 0, err
	}
	return height, err
}

func (db *DBService) PutEventTraceHeight(chainId int, height int64) error {
	sqlStr := ""
	_, err := db.DB.Exec(sqlStr, height, chainId)
	return err
}
