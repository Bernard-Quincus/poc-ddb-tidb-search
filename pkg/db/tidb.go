package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"crypto/tls"

	"github.com/go-sql-driver/mysql"
	logger "github.com/sirupsen/logrus"
)

func init() {
	logger.SetFormatter(&logger.JSONFormatter{})
}

const TiDB_DatabaseName = "dispatchDB"

type tiDB struct {
	db        *sql.DB
	dbName    string
	dsn       string
	tableName string
}

type tiDBResult struct {
}

func NewTiDB(dbName string) (DB, error) {
	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: "gateway01.ap-southeast-1.prod.aws.tidbcloud.com",
	})

	if dbName == "" {
		dbName = "test"
	}

	dsn := fmt.Sprintf("anN5DThCfX3tUZr.root:Q9qyi4i6Zy7XL3Xw@tcp(gateway01.ap-southeast-1.prod.aws.tidbcloud.com:4000)/%s?tls=tidb", dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &tiDB{
		db: db,
	}, nil
}

func (tidb *tiDB) SetTableName(name string) error {
	if name == "" {
		return errors.New("empty table name")
	}
	tidb.tableName = name
	return nil
}

func (tidb *tiDB) Put(input ...any) (any, error) {

	tx, err := tidb.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		logger.Info("tx:", err.Error())
		return nil, err
	}

	for _, sql := range input {
		_, execErr := tx.Exec(sql.(string))
		if execErr != nil {
			logger.Info("execErr:", execErr.Error())
			_ = tx.Rollback()
			return nil, execErr
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Info("commit err:", err.Error())
		return nil, err
	}

	return nil, nil
}

func (tidb *tiDB) Get(input any) (any, error) {

	return nil, nil
}

func (tidb *tiDB) Delete(input any) error {

	return nil
}

func (tidb *tiDB) Search(input any) (any, error) {

	return nil, nil
}

func (tidb *tiDB) Close() error {
	return tidb.db.Close()
}

func (tidb *tiDB) GetTiDBConn() *sql.DB {
	return tidb.db
}
