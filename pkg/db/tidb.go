package db

import "errors"

type tiDB struct {
	tableName string
}

type tiDBResult struct {
}

func NewTiDB() DB {
	return &tiDB{}
}

func (tidb *tiDB) SetTableName(name string) error {
	if name == "" {
		return errors.New("empty table name")
	}
	tidb.tableName = name
	return nil
}

func (tidb *tiDB) Put(input any) (any, error) {

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

func (tidb *tiDB) Close(input any) error {

	return nil
}
