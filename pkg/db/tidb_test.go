package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTiDB(t *testing.T) {
	db, err := NewTiDB(TiDB_DatabaseName)
	assert.NoError(t, err)

	defer db.Close()

	tidb := db.GetTiDBConn()

	var dbName string
	err = tidb.QueryRow("SELECT DATABASE();").Scan(&dbName)
	if err != nil {
		t.Log("failed to execute query", err)
	}
	assert.Equal(t, dbName, "dispatchDB")

	rows, err := tidb.QueryContext(context.Background(), "SHOW TABLES")
	assert.NoError(t, err)

	defer rows.Close()

	type tableName struct {
		name string
	}
	var names []tableName

	for rows.Next() {
		var tn tableName
		if err := rows.Scan(&tn.name); err != nil {
			t.Log(err)
		}
		names = append(names, tn)
	}
	if err = rows.Err(); err != nil {
		t.Log(err)
	}

	for rows.NextResultSet() {
		for rows.Next() {
			var tn tableName
			if err := rows.Scan(&tn.name); err != nil {
				t.Log(err)
			}
			names = append(names, tn)
		}
	}

	t.Log("table names: ", names)
	assert.Greater(t, len(names), 0)
}

func TestFetchRecords(t *testing.T) {
	db, err := NewTiDB(TiDB_DatabaseName)
	assert.NoError(t, err)

	defer db.Close()

	tidb := db.GetTiDBConn()

	rows, err := tidb.QueryContext(context.Background(), "Select uuid, shipment_id from jobs")
	assert.NoError(t, err)

	defer rows.Close()

	type jobRec struct {
		uuid        string
		shipment_id string
	}
	var recs []jobRec

	for rows.Next() {
		var tn jobRec
		if err := rows.Scan(&tn.uuid, &tn.shipment_id); err != nil {
			t.Log(err)
		}
		recs = append(recs, tn)
	}
	if err = rows.Err(); err != nil {
		t.Log(err)
	}

	for rows.NextResultSet() {
		for rows.Next() {
			var tn jobRec
			if err := rows.Scan(&tn.uuid, &tn.shipment_id); err != nil {
				t.Log(err)
			}
			recs = append(recs, tn)
		}
	}

	t.Log("job records: ", recs, len(recs))
	assert.Greater(t, len(recs), 0)
}
