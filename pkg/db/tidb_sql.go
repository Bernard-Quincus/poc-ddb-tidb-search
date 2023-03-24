package db

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"poc-ddb-tidb-search/pkg/models"
	"poc-ddb-tidb-search/pkg/query"
	"strings"
	"time"
)

var (
	InsertJobSQL     = `Insert into jobs %s values (%s)`
	InsertJobRefSQL  = `Insert into jobs_reference %s values (%s)`
	UpdateJobSQL     = `Update jobs set %s`           // deferred
	UpdateJobsRefSQL = `Update jobs_reference set %s` // deferred
)

func MakeInsertJobSQLStatement(job *models.Job, uuid string) (string, error) {
	keys := "(`uuid`,`org_id`,`shipment_id`,`job_id`,`order_id`,`status`,`start_time`,`commit_time`,`detail`)"
	vals := `"%s","%s","%s","%s","%s","%s","%s","%s","%s"`

	jobJson, err := json.Marshal(job)
	if err != nil {
		return "", err
	}
	jobBlob := base64.StdEncoding.EncodeToString(jobJson)

	layout := "02/01/2006 5:04:05" // dd/mm/yyyy time from jobs payload
	startTime := job.PickupDate + " " + job.PickupStartTime
	st, _ := time.Parse(layout, startTime)

	commitTime := job.DeliveryDate + " " + job.DeliveryCommitTime
	ct, _ := time.Parse(layout, commitTime)

	vals = fmt.Sprintf(vals, uuid, job.OrgID2, job.RefShipmentID, job.ID, job.OrderPayload.OrderID, job.Status,
		st.Format(time.RFC3339), ct.Format(time.RFC3339), jobBlob)

	return fmt.Sprintf(InsertJobSQL, keys, vals), nil
}

func MakeInsertJobReferenceSQLStatement(job *models.Job, uuid string) (string, error) {
	keys := "(`uuid`,`org_id`,`shipment_tags`,`order_refids`,`shipment_ref_ids`,`assigned_vendor`,`assigned_facility`,`job_postal_code`,`job_city`,`job_street`,`customer_account_name`,`sender_name`,`consignee_name`)"
	vals := `"%s","%s","%s","%s","%s","%s","%s","%s","%s","%s","%s","%s","%s"`

	tags := ""
	if len(job.PackageTags) > 0 {
		tags = strings.Join(job.PackageTags, ",")
	}

	facility := ""
	if job.FromFacility != nil {
		facility = job.FromFacility.FacilityName
	}

	senderName := ""
	if job.OrderPayload != nil && job.OrderPayload.ShipperInfo != nil {
		senderName = job.OrderPayload.ShipperInfo.Name
	}

	consigneeName := ""
	if job.OrderPayload != nil && job.OrderPayload.ConsigneeInfo != nil {
		consigneeName = job.OrderPayload.ConsigneeInfo.Name
	}

	vals = fmt.Sprintf(vals, uuid, job.OrgID2, tags, job.RefOrderID, job.RefShipmentID, job.PartnerName,
		facility, job.DeliveryPostcode, job.DeliveryCity, job.DeliveryAddress, "none yet", senderName, consigneeName)

	return fmt.Sprintf(InsertJobRefSQL, keys, vals), nil
}

func MakeSearchSQLStatements(params *query.JobSearchParams, orgID string) []string {
	stmts := make([]string, 0)
	kv := make([]string, 0)
	joinPrefixA := "a."
	joinPrefixB := "b."
	needToJoin := false

	// check if we got jobs_reference table field to query
	if params.ConsigneeName != "" || params.FacilityName != "" || params.OrderRefTags != "" || params.SenderName != "" ||
		params.ShipmentTags != "" || params.VendorName != "" {
		needToJoin = true
	}

	if !needToJoin {
		joinPrefixA = ""
		joinPrefixB = ""
	}

	// jobs table
	if params.JobID != "" {
		kv = append(kv, fmt.Sprintf("%sjob_id='%s'", joinPrefixA, params.JobID))
	}
	if params.OrderID != "" {
		kv = append(kv, fmt.Sprintf("%sorder_id='%s'", joinPrefixA, params.OrderID))
	}
	if params.ShipmentID != "" {
		kv = append(kv, fmt.Sprintf("%shipment_id='%s'", joinPrefixA, params.ShipmentID))
	}
	if params.Status != "" {
		kv = append(kv, fmt.Sprintf("%sstatus='%s'", joinPrefixA, params.Status))
	}
	if params.StartTime != "" {
		kv = append(kv, fmt.Sprintf("cast(%sstart_time as date)='%s'", joinPrefixA, params.StartTime)) // yyyy-mm-dd
	}
	if params.CommitTime != "" {
		kv = append(kv, fmt.Sprintf("cast(%scommit_time as date)='%s'", joinPrefixA, params.CommitTime)) // yyyy-mm-dd
	}

	// job_refs table
	if params.ShipmentTags != "" {
		kv = append(kv, fmt.Sprintf("%sshipment_tags like '%s'", joinPrefixB, "%"+params.ShipmentTags+"%"))
	}
	if params.OrderRefTags != "" {
		kv = append(kv, fmt.Sprintf("%sorder_refids like '%s'", joinPrefixB, "%"+params.OrderRefTags+"%"))
	}
	if params.VendorName != "" {
		kv = append(kv, fmt.Sprintf("%sassigned_vendor='%s'", joinPrefixB, params.VendorName))
	}
	if params.FacilityName != "" {
		kv = append(kv, fmt.Sprintf("%sassigned_facility='%s'", joinPrefixB, params.FacilityName))
	}
	if params.SenderName != "" {
		kv = append(kv, fmt.Sprintf("%ssender_name='%s'", joinPrefixB, params.SenderName))
	}
	if params.ConsigneeName != "" {
		kv = append(kv, fmt.Sprintf("%sconsignee_name='%s'", joinPrefixB, params.ConsigneeName))
	}

	qfields := strings.Join(kv, " and ")

	q := fmt.Sprintf("Select uuid, detail from %s.jobs where org_id='%s' and %s limit %d, %d", TiDB_DatabaseName, orgID, qfields, params.PageNumber*params.PageSize, params.PageSize)
	q2 := fmt.Sprintf("Select count(*) as totalrec from %s.jobs where org_id='%s' and %s", TiDB_DatabaseName, orgID, qfields)

	if needToJoin {
		q = fmt.Sprintf("Select a.uuid, a.detail from %s.jobs a left join %s.jobs_reference b on a.uuid = b.uuid where a.org_id='%s' and %s limit %d, %d",
			TiDB_DatabaseName, TiDB_DatabaseName, orgID, qfields, params.PageNumber*params.PageSize, params.PageSize)

		q2 = fmt.Sprintf("Select count(*) as totalrec from %s.jobs a left join %s.jobs_reference b on a.uuid = b.uuid where a.org_id='%s' and %s",
			TiDB_DatabaseName, TiDB_DatabaseName, orgID, qfields)
	}

	stmts = append(stmts, q, q2)

	return stmts
}
