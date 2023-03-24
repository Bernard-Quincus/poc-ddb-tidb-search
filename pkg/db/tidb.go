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

type TiDBRow struct {
	UUID   string
	Detail string
}

type TiDBResult struct {
	TotalItems int
	Details    []*TiDBRow
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

func (tidb *tiDB) Search(input ...any) (any, error) {
	result := make([]*TiDBRow, 0)
	totalItems := 0

	for i, sql := range input {
		if i == 0 {
			res, err := tidb.execQuery(sql.(string))
			if err != nil {
				return nil, err
			}
			result = append(result, res...)
			continue
		}

		count, err := tidb.execTotalCountQuery(sql.(string))
		if err != nil {
			return nil, err
		}
		totalItems = count
	}

	return &TiDBResult{TotalItems: totalItems, Details: result}, nil
}

func (tidb *tiDB) Close() error {
	return tidb.db.Close()
}

func (tidb *tiDB) GetTiDBConn() *sql.DB {
	return tidb.db
}

func (tidb *tiDB) execQuery(q string) ([]*TiDBRow, error) {

	rows, err := tidb.db.QueryContext(context.Background(), q)
	defer rows.Close()

	var result = make([]*TiDBRow, 0)

	// temp
	fmt.Println("rows==>", rows)

	for rows.Next() {
		var r = new(TiDBRow)
		if err := rows.Scan(&r.UUID, &r.Detail); err != nil {
			return nil, err
		}

		// temp
		fmt.Println("r==>", r.Detail, r.UUID)

		result = append(result, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	for rows.NextResultSet() {
		for rows.Next() {
			var r = new(TiDBRow)
			if err := rows.Scan(&r.UUID, &r.Detail); err != nil {
				return nil, err
			}

			// temp
			fmt.Println("r2 ==>", r.Detail, r.UUID)
			result = append(result, r)
		}
	}

	return result, nil
}

func (tidb *tiDB) execTotalCountQuery(q string) (int, error) {

	rows, err := tidb.db.QueryContext(context.Background(), q)
	defer rows.Close()

	var count = 0
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}
	if err = rows.Err(); err != nil {
		return 0, err
	}

	return count, nil
}
