package db

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"poc-ddb-tidb-search/pkg/models"
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
