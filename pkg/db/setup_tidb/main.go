package main

import (
	"context"

	"poc-ddb-tidb-search/pkg/db"

	logger "github.com/sirupsen/logrus"
)

// Only run this once, run it manually via terminal

func main() {
	tiDB, err := db.NewTiDB("test") // use default db name first
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "TiDBErr",
		}).Error("failed to connect to TiDB instance")
	}
	defer tiDB.Close()

	if err := setupDB(tiDB); err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "TiDBErr",
		}).Error("failed to setup TiDB database")
		return
	}

	logger.Info("TiDB Database setup was successful")
}

func setupDB(tiDB db.DB) error {

	tidb := tiDB.GetTiDBConn()
	ctx := context.Background()

	_, err := tidb.ExecContext(ctx, "SET GLOBAL tidb_multi_statement_mode='ON'")
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "TiDBErr",
		}).Error("failed to turn on multi statement")
		return err
	}

	// tidb.ExecContext(ctx, "USE dispatchDB; DROP TABLE jobs;")
	// tidb.ExecContext(ctx, "USE dispatchDB; DROP TABLE jobs_reference;")
	// return nil

	// dropIndices(ctx, tiDB)
	// return nil

	if err := setupTables(ctx, tiDB); err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "TiDBErr",
		}).Error("failed to setup TiDB tables and indices")
		return err
	}

	_, err = tidb.ExecContext(ctx, "SET GLOBAL tidb_multi_statement_mode='OFF'")
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "TiDBErr",
		}).Error("failed to turn off multi statement")
		return err
	}

	return nil
}

func setupTables(ctx context.Context, tiDB db.DB) error {
	tidb := tiDB.GetTiDBConn()

	createTableJobs := `
		USE dispatchDB;

		CREATE TABLE IF NOT EXISTS jobs (
			uuid Varchar(36) PRIMARY KEY,
			org_id Varchar(50),
			shipment_id Varchar(100),
			job_id Varchar(100),
			order_id Varchar(100),
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, 
			status Varchar(50),
			start_time DATETIME,
			commit_time DATETIME,
			detail MEDIUMBLOB
		);
	`

	createTableJobsRefs := `
		USE dispatchDB;

		CREATE TABLE IF NOT EXISTS jobs_reference (
			uuid Varchar(36) PRIMARY KEY,
			org_id Varchar(50),
			shipment_tags Varchar(500),
			order_refids Varchar(100),
			shipment_ref_ids Varchar(100),
			assigned_vendor Varchar(500),
			assigned_facility Varchar(500),
			job_postal_code Varchar(100),
			job_city Varchar(500),
			job_street Varchar(500),
			customer_account_name Varchar(500),
			sender_name Varchar(500),
			consignee_name Varchar(500),
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);
	`
	_, err := tidb.ExecContext(ctx, "SET GLOBAL tidb_multi_statement_mode='ON'")
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "TiDBErr",
		}).Error("failed to create database")
		return err
	}

	_, err = tidb.ExecContext(ctx, createTableJobs)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "TiDBErr",
		}).Error("failed to create jobs table")
		return err
	}

	_, err = tidb.ExecContext(ctx, createTableJobsRefs)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "TiDBErr",
		}).Error("failed to create jobs_reference table")
		return err
	}

	return createIndices(ctx, tiDB)
}

func createIndices(ctx context.Context, tiDB db.DB) error {
	tidb := tiDB.GetTiDBConn()

	createJobxIndices := `
		USE dispatchDB;

		CREATE INDEX shpID_idx ON jobs (
			org_id,shipment_id
		);

		CREATE INDEX jobID_idx ON jobs (
			org_id,job_id
		);
		
		CREATE INDEX orderID_idx ON jobs (
			org_id,order_id
		);

		CREATE INDEX dateStatus_idx ON jobs (
			org_id,updated_at,status
		);

		CREATE INDEX start_time_idx ON jobs (
			org_id,start_time
		);

		CREATE INDEX commit_time_idx ON jobs (
			org_id,commit_time
		);
	`

	createRefsIndices := `
		USE dispatchDB;

		CREATE INDEX shpTags_idx ON jobs_reference (
			org_id,shipment_tags
		);

		CREATE INDEX ordRefs_idx ON jobs_reference (
			org_id,order_refids
		);

		CREATE INDEX shpRefs_idx ON jobs_reference (
			org_id,shipment_ref_ids
		);

		CREATE INDEX dateVendor_idx ON jobs_reference (
			org_id,updated_at,assigned_vendor
		);

		CREATE INDEX dateFacility_idx ON jobs_reference (
			org_id,updated_at,assigned_facility
		);

		CREATE INDEX datePostal_idx ON jobs_reference (
			org_id,updated_at,job_postal_code
		);

		CREATE INDEX dateCity_idx ON jobs_reference (
			org_id,updated_at,job_city
		);

		CREATE INDEX dateStreet_idx ON jobs_reference (
			updated_at,job_street
		);

		CREATE INDEX dateAccountName_idx ON jobs_reference (
			org_id,updated_at,customer_account_name
		);

		CREATE INDEX dateSender_idx ON jobs_reference (
			org_id,updated_at,sender_name
		);

		CREATE INDEX dateConsignee_idx ON jobs_reference (
			org_id,updated_at,consignee_name
		);
	`
	_, err := tidb.ExecContext(ctx, "SET GLOBAL tidb_multi_statement_mode='ON'")
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "TiDBErr",
		}).Error("failed to create database")
		return err
	}

	_, err = tidb.ExecContext(ctx, createJobxIndices)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "TiDBErr",
		}).Error("failed to create index on jobs table")
		return err
	}

	_, err = tidb.ExecContext(ctx, createRefsIndices)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "TiDBErr",
		}).Error("failed to create index on jobs_reference table")
		return err
	}

	return nil
}

func dropIndices(ctx context.Context, tiDB db.DB) error {
	tidb := tiDB.GetTiDBConn()

	dropIndices := `
		USE dispatchDB;

		DROP INDEX shpID_idx ON jobs;

		DROP INDEX jobID_idx ON jobs;

		DROP INDEX orderID_idx ON jobs;

		DROP INDEX dateStatus_idx ON jobs;

		DROP INDEX datePriority_idx ON jobs;

		DROP INDEX start_time_idx ON jobs;

		DROP INDEX commit_time_idx ON jobs;
		`

	dropRefsIndices := `
		USE dispatchDB;

		DROP INDEX shpTags_idx ON jobs_reference;

		DROP INDEX ordRefs_idx ON jobs_reference;

		DROP INDEX shpRefs_idx ON jobs_reference;

		DROP INDEX dateVendor_idx ON jobs_reference;

		DROP INDEX dateFacility_idx ON jobs_reference;

		DROP INDEX datePostal_idx ON jobs_reference;

		DROP INDEX dateCity_idx ON jobs_reference;

		DROP INDEX dateStreet_idx ON jobs_reference;

		DROP INDEX dateAccountName_idx ON jobs_reference;

		DROP INDEX dateSender_idx ON jobs_reference;

		DROP INDEX dateConsignee_idx ON jobs_reference;
		`

	_, err := tidb.ExecContext(ctx, "SET GLOBAL tidb_multi_statement_mode='ON'")
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "TiDBErr",
		}).Error("failed to create database")
		return err
	}

	_, err = tidb.ExecContext(ctx, dropIndices)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "TiDBErr",
		}).Error("failed to create index on jobs table")
		return err
	}

	_, err = tidb.ExecContext(ctx, dropRefsIndices)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "TiDBErr",
		}).Error("failed to create index on jobs_reference table")
		return err
	}

	return nil
}
