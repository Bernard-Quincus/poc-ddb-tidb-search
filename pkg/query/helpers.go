package query

import (
	"errors"
	"fmt"
	"poc-ddb-tidb-search/pkg/models"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

type JobSearchParams struct {
	//jobs table
	ShipmentID string
	OrderID    string
	JobID      string
	Status     string
	StartTime  string // cast to date
	CommitTime string // cast to date

	//ref table
	ShipmentTags  string // like
	OrderRefTags  string // like
	ConsigneeName string
	SenderName    string
	VendorName    string
	FacilityName  string

	PageSize   int
	PageNumber int
}

type JobRow struct {
	UUID string      `json:"uuid"`
	Job  *models.Job `json:"job"`
}

type JobSearchResult struct {
	Data         []*JobRow `json:"data"`
	PageSize     int       `json:"page_size"`
	TotalItems   int       `json:"total_items"`
	ResponseTime string    `json:"response_time"`
}

var allQueryParameters = []string{
	"shipment_id",
	"order_id",
	"job_id",
	"status",
	"start_time",
	"commit_time",
	"shipment_tags",
	"order_tags",
	"consignee_name",
	"sender_name",
	"vendor_name",
	"facility_name",
	"page_size",
	"page_number",
}

const pageSize = 20

func ParametersFromRequest(request *events.APIGatewayProxyRequest) (*JobSearchParams, error) {
	p := &JobSearchParams{PageNumber: 0, PageSize: pageSize}

	if err := p.readQueryParameters(request); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *JobSearchParams) readQueryParameters(request *events.APIGatewayProxyRequest) error {
	hasParamValue := false

	for _, paramName := range allQueryParameters {
		param, exists := request.QueryStringParameters[paramName]
		if !exists {
			continue
		}

		if param != "" {
			hasParamValue = true
		}

		switch paramName {
		case "shipment_id":
			p.ShipmentID = param
		case "order_id":
			p.OrderID = param
		case "job_id":
			p.JobID = param
		case "status":
			p.Status = param
		case "start_time":
			p.StartTime = param
		case "commit_time":
			p.CommitTime = param
		case "shipment_tags":
			p.ShipmentTags = param
		case "order_tags":
			p.OrderRefTags = param
		case "consignee_name":
			p.ConsigneeName = param
		case "sender_name":
			p.SenderName = param
		case "vendor_name":
			p.VendorName = param
		case "facility_name":
			p.FacilityName = param
		case "page_size":
			pageSize, err := strconv.ParseInt(param, 10, 32)
			if err != nil {
				return fmt.Errorf("page size is not a valid integer: %v", err)
			}
			p.PageSize = int(pageSize)
		case "page_number":
			pageNum, err := strconv.ParseInt(param, 10, 32)
			if err != nil {
				return fmt.Errorf("page number is not a valid integer: %v", err)
			}
			p.PageNumber = int(pageNum)
		}
	}

	if !hasParamValue {
		return errors.New("no valid params found")
	}
	return nil
}

func GetOrgID(request *events.APIGatewayProxyRequest) string {
	orgID, ok := request.Headers["ORGID"]
	if !ok {
		orgID, ok = request.Headers["orgid"]
		if !ok {
			return ""
		}
	}

	return orgID
}
