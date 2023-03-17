package models

import (
	"encoding/json"
	"strconv"
	"strings"
)

// ServiceType is the service type of shipment
type ServiceType string

const (
	ServiceTypeUnknown ServiceType = "service type unknown"

	ServiceTypePickUp   ServiceType = "pickup"
	ServiceTypeDelivery ServiceType = "delivery"
)

var (
	// StrWarpServiceType is string service type warp to ServiceType const
	StrWarpServiceType = map[string]ServiceType{
		"pickup":   ServiceTypePickUp,
		"delivery": ServiceTypeDelivery,
	}
)

// Status is shipment status
type Status string

const (
	StatusUnknown Status = "status unknown"

	StatusNew                  Status = "new"
	StatusUnallocated          Status = "unallocated"
	StatusDraft                Status = "draft"
	StatusOnRoute              Status = "onroute"
	StatusCompleted            Status = "completed"
	StatusForceCompleted       Status = "force completed"
	StatusFailed               Status = "failed"
	StatusCancelled            Status = "cancelled"
	StatusOnVehicleForDelivery Status = "on vehicle for delivery"
	StatusArrivedInHub         Status = "arrived in hub"
	StatusSkipShipment         Status = "skip shipment"
	StatusBeenConsolidated     Status = "been consolidated"
	StatusDriverArrived        Status = "driver arrived"

	StatusDeleted Status = "deleted"
	StatusInvalid Status = "invalid" // invalid status to describe the record invalid
)

var (
	// StrWarpShipmentStatus is string status warp Status const
	StrWarpShipmentStatus = map[string]Status{
		"new":                     StatusNew,
		"unallocated":             StatusUnallocated,
		"draft":                   StatusDraft,
		"onroute":                 StatusOnRoute,
		"completed":               StatusCompleted,
		"force completed":         StatusForceCompleted,
		"failed":                  StatusFailed,
		"cancelled":               StatusCancelled,
		"on vehicle for delivery": StatusOnVehicleForDelivery,
		"arrived in hub":          StatusArrivedInHub,
		"skip shipment":           StatusSkipShipment,
		"been consolidated":       StatusBeenConsolidated,
		"driver arrived":          StatusDriverArrived,
	}
)

// MeasurementUnitInfo contains the information about the
// measurement units (weights & dimensions)
type MeasurementUnitInfo struct {
	Dimensions string `json:"dimensions"`
	Weight     string `json:"weight"`
	Volume     string `json:"volume"`
}

// HandleFacilityInfo is similary but detail is differnt from "Facility Info" (based on DynamoDB), therefore, I have to divide they
type HandleFacilityInfo struct {
	ID           string  `json:"id"`
	Name         string  `json:"facility_name"`
	CityID       string  `json:"city_id"`
	CityName     string  `json:"city_name"`
	StateID      string  `json:"state_id"`
	StateName    string  `json:"state_name"`
	CountryID    string  `json:"country_id"`
	CountryName  string  `json:"country_name"`
	Address      string  `json:"address"`
	PostCode     string  `json:"postcode"`
	StartLat     float64 `json:"start_latitude"`
	StartLong    float64 `json:"start_longitude"`
	DepotLat     float64 `json:"depot_latitude"`
	DepotLong    float64 `json:"depot_longitude"`
	IANATimezone string  `json:"iana_timezone"`
	UTCTimezone  string  `json:"utc_timezone"`
}

// FacilityInfo contains information about a facility (hub in Shipment)
type FacilityInfo struct {
	ID             string                      `json:"id"`
	Address        string                      `json:"address"`
	City           string                      `json:"city"`
	State          string                      `json:"state"`
	Country        string                      `json:"country"`
	Postcode       string                      `json:"postal_code"`
	Address1       string                      `json:"address_line1"`
	Address2       string                      `json:"address_line2"`
	Address3       string                      `json:"address_line3"`
	Latitude       float64                     `json:"latitude"`
	Longitude      float64                     `json:"longitude"`
	CountryID      string                      `json:"country_id"`
	StateID        string                      `json:"state_id"`
	CityID         string                      `json:"city_id"`
	PersonInCharge *FacilityInChargePeopleInfo `json:"person_in_charge"`
}

// Job is shipment's dispatch job information
type Job struct {
	OrgID2 string `json:"orgID"` // test
	DocID  string `json:"docID"` // test
	PK     string `json:"PK"`
	SK     string `json:"SK"`
	GSI1PK string `json:"GSI1PK,omitempty"`
	GSI1SK string `json:"GSI1SK,omitempty"`
	GSI2PK string `json:"GSI2PK,omitempty"`
	GSI2SK string `json:"GSI2SK,omitempty"`

	// GSI3PK follows the format `ORG#{org_id}#ORDER#{order_QR_id}` for the query that gets list of jobs by for an order
	GSI3PK string `json:"GSI3PK,omitempty"`

	// GSI3SK follows the format `SEGMENT#{segment_id}#PACKAGE#{package_id}` for the query that get list of jobs for a particular segment in Segment Journey.
	// For deleted jobs, key format of the GSI3SK will be `DELETED#SEGMENT#{segment_id}#PACKAGE#{package_id}`
	GSI3SK string `json:"GSI3SK,omitempty"`

	ItemType ItemType `json:"entity_type"`

	// Reference IDs from other module
	// UUID format
	RefOrderID        string `json:"ref_order_id"`
	RefShipmentID     string `json:"ref_shipment_id"`
	RefConsolidatedID string `json:"ref_consolidated_shipment_id"`
	RefPackageID      string `json:"ref_package_id"`
	ReferenceID       string `json:"reference_id"`
	// Human Readable IDs
	RefOrderLabel        string `json:"ref_order_label"`
	RefShipmentLabel     string `json:"ref_shipment_label"`
	RefPackageLabel      string `json:"ref_package_label"`
	RefConsolidatedLabel string `json:"ref_consolidated_shipment_label"`
	RefShipSeqNo         string `json:"ref_ship_seq_no"`

	ID                    string              `json:"shipment_id"` // note: this is actual JOB_ID what we expected
	OriginalShipmentID    string              `json:"original_shipment_id"`
	Code                  string              `json:"shipment_code"`
	Currency              string              `json:"currency_type"`
	CustomerMobile        string              `json:"customer_mobile_no"`
	ClientOrderID         string              `json:"client_order_id"`
	ClientOrderCode       string              `json:"client_order_code"`
	ETA                   int64               `json:"eta"`
	JobID                 string              `json:"job_id"` // this is legacy naming which is dispatch ID
	JobName               string              `json:"job_name"`
	LocationID            string              `json:"location_id"` // handled by which facility (this is facility ID)
	HandleFacilityInfo    *HandleFacilityInfo `json:"handle_facility_detail,omitempty"`
	RouteType             *RoutingType        `json:"route_type"`
	PackageCount          int64               `json:"number_of_packages"`
	OrderDate             string              `json:"order_date"`
	OrderTime             string              `json:"order_time"`
	OrgID                 string              `json:"organisation_id"`
	OrgUnitID             string              `json:"organisation_unit_id"`
	PackageID             string              `json:"template_id"`
	PackageCode           string              `json:"package_code"` // package code comes from csv ID column
	TrackingID            string              `json:"tracking_id"`  // tracking id is reference id to shipment module id
	PackageType           string              `json:"package_type"`
	PackagingType         *PackagingType      `json:"packaging_type"`
	Commodity             string              `json:"commodity"`
	RecipientLine         string              `json:"recipient_line"`
	ShipperLine           string              `json:"shipper_line"`
	Region                string              `json:"region"`
	ServiceType           ServiceType         `json:"service_type"`
	TimeWindowStart       int64               `json:"timewindowstart"`
	TimeWindowEnd         int64               `json:"timewindowend"`
	UpdateAtUnixsec       int64               `json:"updated_at"`
	UpdateBy              string              `json:"updated_by"`
	CreatedBy             string              `json:"created_by"`
	Volume                float64             `json:"volume"`
	Weight                float64             `json:"weight"`
	AlgoID                string              `json:"algo_id"`
	AlgoType              string              `json:"algo_type"`
	BookedAt              string              `json:"booked_at"`
	RouteNo               string              `json:"route_no"`
	SectionNumber         int64               `json:"section_number"`
	SequenceNo            int64               `json:"sequence_no"`
	Status                Status              `json:"status"`
	StatusUpdateAtUnixsec int64               `json:"status_updated_at"`
	StatusGPSLat          float64             `json:"status_lat"`
	StatusGPSLon          float64             `json:"status_lon"`

	PathSequence int `json:"path_sequence"`

	Comment string `json:"comment"`

	GPSLat float64 `json:"lat"`
	GPSLon float64 `json:"lon"`

	ActualFailedTimeUnixsec int64             `json:"actual_failed_time"`
	FailedGPSLat            float64           `json:"failed_lat"`
	FailedGPSLon            float64           `json:"failed_lon"`
	FailedReport            *JobFailureReport `json:"failed_report"`

	ActualDeliveryTimeUnixsec int64   `json:"actual_delivery_time"`
	CompletedGPSLat           float64 `json:"completed_lat"`
	CompletedGPSLon           float64 `json:"completed_lon"`

	DriverArrived              bool            `json:"driver_arrived"`
	DriverArrivalTimeUnixSec   int64           `json:"driver_arrival_time"`
	DriverCompletedTimeUnixSec int64           `json:"driver_completed_time"`
	StatusRecords              []*StatusRecord `json:"status_records"`

	FromLocationID                string           `json:"from_location_id"`
	FromFacility                  *HubLocationInfo `json:"from_facility_detail,omitempty"`
	ToLocationID                  string           `json:"to_location_id"`
	ToFacility                    *HubLocationInfo `json:"to_facility_detail,omitempty"`
	IsIDRequiredForTargetFacility *bool            `json:"is_id_required_for_target_facility,omitempty"`

	// Shipment label/tags
	// IsODS (based on organisation-on-demand setting logic)
	IsODS             bool             `json:"is_ods" dynamodbav:"is_ods_shipment"`
	ODSMetric         *ODSMetric       `json:"ods_metric,omitempty"`
	ODSResponses      ODSDriverHistory `json:"ods_response_log,omitempty"`
	ODSLatestResponse *ODSAssignment   `json:"ods_latest_response,omitempty"`

	MeasurementUnits    *MeasurementUnitInfo `json:"measurement_units,omitempty"`
	ShipmentServiceType string               `json:"shipment_service_type,omitempty"`

	// Consolidated shipment related fields
	IsConsolidated          bool                `json:"is_consolidated"`
	ConsolidatedOrigin      *FacilityInfo       `json:"consolidated_origin,omitempty"`
	ConsolidatedDestination *FacilityInfo       `json:"consolidated_destination,omitempty"`
	ConsolidatedItemList    []*ConsolidatedItem `json:"consolidated_items,omitempty"`

	// Partner information
	PartnerID   string `json:"partner_id"`
	PartnerName string `json:"partner_name"`

	// Pickup information
	PickupGPSLat            float64 `json:"pickup_latitude"`
	PickupGPSLon            float64 `json:"pickup_longitude"`
	PickupPostcode          string  `json:"pickup_postcode"`
	PickupCountry           string  `json:"pickup_country"`
	PickupCountryID         string  `json:"pickup_country_id"`
	PickupState             string  `json:"pickup_state"`
	PickupStateID           string  `json:"pickup_state_id"`
	PickupCity              string  `json:"pickup_city"`
	PickupCityID            string  `json:"pickup_city_id"`
	PickupAddress           string  `json:"pickup_address"`
	PickupDate              string  `json:"pickup_date"`
	PickupCommitTime        string  `json:"pickup_commit_time"`
	PickupCommitTimeUnixsec int64   `json:"pickup_commit_time_ts"`
	PickupCommitTimestamp   string  `json:"pickup_commit_timestamp"`
	PickupServiceTime       int64   `json:"pickup_service_time"`
	PickupStartTime         string  `json:"pickup_start_time"`
	PickupStartTimezone     string  `json:"pickup_start_time_timezone"`
	PickupStartTimeUnixsec  int64   `json:"pickup_start_time_ts"`
	PickupStartTimestamp    string  `json:"pickup_start_timestamp"`

	// Delivery information
	DeliveryGPSLat            float64 `json:"delivery_latitude"`
	DeliveryGPSLon            float64 `json:"delivery_longitude"`
	DeliveryPostcode          string  `json:"delivery_postcode"`
	DeliveryCountry           string  `json:"delivery_country"`
	DeliveryCountryID         string  `json:"delivery_country_id"`
	DeliveryState             string  `json:"delivery_state"`
	DeliveryStateID           string  `json:"delivery_state_id"`
	DeliveryCity              string  `json:"delivery_city"`
	DeliveryCityID            string  `json:"delivery_city_id"`
	DeliveryAddress           string  `json:"delivery_address"`
	DeliveryDate              string  `json:"delivery_date"`
	DeliveryServiceTime       int64   `json:"delivery_service_time"`
	DeliveryStartTime         string  `json:"delivery_start_time"`
	DeliveryStartTimestamp    string  `json:"delivery_start_timestamp"`
	DeliveryStartTimezone     string  `json:"delivery_start_time_timezone"`
	DeliveryStartTimeUnixsec  int64   `json:"delivery_start_time_ts"`
	DeliveryCommitTime        string  `json:"delivery_commit_time"`
	DeliveryCommitTimeUnixsec int64   `json:"delivery_commit_time_ts"`
	DeliveryCommitTimestamp   string  `json:"delivery_commit_timestamp"`

	// Upload files (POD)
	PODfilename                      string     `json:"file_name"`
	PODfileCaptureTimeUnixsec        int64      `json:"file_captured_at"`
	PODfileUploadTimeUnixsec         int64      `json:"file_upload_time"`
	PODSignatureFilename             string     `json:"signature_file_name"`
	PODSignatureCaptureTimeUnixsec   int64      `json:"signature_captured_at"`
	PODSignatureCaptureUploadUnixsec int64      `json:"signature_upload_time"`
	PODPresignedURL                  string     `json:"pod_url"`
	PODSignaturePresignedURL         string     `json:"pod_signature_url"`
	PODAssets                        *PODAssets `json:"pod_assets"`

	// Inbound scan
	InboundScanRequired bool               `json:"inbound_scan_required"`
	InboundScanResult   *InboundScanResult `json:"inbound_scan_result"`

	// POD verification scan
	PODScanRequired bool           `json:"pod_scan_required"`
	PODScanResult   *PODScanResult `json:"pod_scan_result"`

	// Package report results
	PackageReport *PackageReport `json:"package_report"`

	// Driver record
	DriverID      string `json:"driver_id"`
	DriverName    string `json:"driver_name"`
	DriverPic     string `json:"driver_pic"`
	VehicleID     string `json:"vehicle_id"`
	VehicleName   string `json:"vehicle_name"`
	VehicleNumber string `json:"vehicle_number"`
	VehicleType   string `json:"vehicle_type"`

	DepotGPSLat float64 `json:"depot_latitude"`
	DepotGSPLon float64 `json:"depot_longitude"`
	Distance    string  `json:"distance"`

	MissedReason string `json:"missed_reason"`

	CODEnabled bool    `json:"cod_enabled"`
	CODAmount  float64 `json:"cod_amount"`

	MilestoneCode string `json:"milestone_code"`

	// NOTE: P2P/H&S related fields
	OpsType            *string `json:"ops_type"` // NOW: HAS TWO TYPE Hub == Service, P2P == Shipment
	PickupFacilityID   *string `json:"pickup_facility_id"`
	DeliveryFacilityID *string `json:"delivery_facility_id"`

	// R&R related fields
	RRModified        bool   `json:"rr_been_modified"`
	RRModifiedISOTime string `json:"rr_modified_time,omitempty"`

	OrderNote      *string       `json:"order_note,omitempty"`
	PackageNote    *string       `json:"package_note,omitempty"`
	PackageTags    []string      `json:"package_tags,omitempty"`
	OrderTagList   []string      `json:"order_tag_list,omitempty"`
	AttachmentList []*Attachment `json:"order_attachment_list,omitempty"`

	PickupInstruction   *string `json:"pickup_instruction,omitempty"`
	DeliveryInstruction *string `json:"delivery_instruction,omitempty"`

	OrderPayload *OrderPayload `json:"order_payload,omitempty"`

	// ReceiverName will contain the receiver's name when an alternative recipient receives the package
	ReceiverName string `json:"receiver_name"`

	ETAStatus string `json:"eta_status"`

	// FlightDepartureInfo is for ground to airport driver usage
	FlightDepartureInfo *FlightDepartureInfo `json:"flight_departure_info,omitempty"`
	// FlightArrivalInfo is for next mile ground to airport pickup usage
	FlightArrivalInfo *FlightArrivalInfo `json:"flight_arrival_info,omitempty"`

	// SegmentID records the ID of a segment in Segment Journey
	SegmentID string `json:"segment_id"`
	// SegmentSequenceNo records the RefID of a segment in segment journey
	SegmentSequenceNo uint32 `json:"segment_sequence_no"`

	// TriggerCode R&R notification temporary variable. store each job notification type for trigger checking
	TriggerCode TriggerCode `json:"-"`
	// OldDeliveryStartTime R&R notification temporary variable. store the old delivery start time
	OldDeliveryStartTime string `json:"-"`
	// OldDeliveryCommitTime R&R notification temporary variable. store the old delivery commit time
	OldDeliveryCommitTime string `json:"-"`
	// OldPickupStartTime R&R notification temporary variable. store the old pickup start time
	OldPickupStartTime string `json:"-"`
	// OldPickupCommitTime R&R notification temporary variable. store the old pickup commit time
	OldPickupCommitTime string `json:"-"`

	// Fields for unoptimized jobs
	IsUnoptimizedJob   bool                `json:"is_unoptimized_job"`
	UnoptimizedJobInfo *UnoptimizedJobInfo `json:"unoptimized_job_info,omitempty"`
}

// UnoptimizedJobInfo describes unoptimized job data required to create the dispatch without GE & LE
type UnoptimizedJobInfo struct {
	PartnerID       string   `json:"partner_id"`
	DriverID        string   `json:"driver_id"`
	VehicleID       string   `json:"vehicle_id"`
	JobETA          string   `json:"job_eta"` // ISO times string
	SegmentDistance *float64 `json:"segment_distance,omitempty"`
}

// StatusRecord contains data when a status was updated
type StatusRecord struct {
	UpdatedAtUnixMs int64   `json:"updated_at_unix_ms"`
	UpdatedBy       string  `json:"updated_by"`
	Lat             float64 `json:"lat"`
	Long            float64 `json:"long"`
	Remarks         string  `json:"remarks"`
}

// FlightDepartureInfo is for ground to airport driver usage
type FlightDepartureInfo struct {
	AirLineNum            string `json:"air_number,omitempty"`
	AirLineCode           string `json:"air_line_code,omitempty"`
	AirLine               string `json:"air_line,omitempty"`
	MasterAirwayBillLabel string `json:"master_airwaybill_label,omitempty"` // AWB for air segments
	DepartedTime          string `json:"departed_time,omitempty"`
	LockoutTime           string `json:"lockout_time,omitempty"`
}

// FlightArrivalInfo is for next mile ground to airport pickup usage
type FlightArrivalInfo struct {
	AirLineNum            string `json:"air_number,omitempty"`
	AirLineCode           string `json:"air_line_code,omitempty"`
	AirLine               string `json:"air_line,omitempty"`
	MasterAirwayBillLabel string `json:"master_airwaybill_label,omitempty"` // AWB for air segments
	ArrivalTime           string `json:"arrival_time,omitempty"`
	RecoverTime           string `json:"recover_time,omitempty"`
}

// RRChangeType describe what kind of change to job after resequenced & reassignment
type RRChangeType string

const (
	// Reassigned means this job as new assignment to the driver
	Reassigned RRChangeType = "this job is new assignment"
	// AffectedJob means this job ETA changed
	AffectedJob RRChangeType = "this job's ETA been updated"
	// ResequencedJob means this job been resequenced
	ResequencedJob RRChangeType = "resequenced job"
	// NewAssignedJob means this job been append from unassigned job
	NewAssignedJob RRChangeType = "new assigned job"
)

// RearrangeJobInfo describe resequenced/reassignment job detail  info
type RearrangeJobInfo struct {
	DispatchID       string       `json:"dispatch_id"`
	JobID            string       `json:"job_id"`
	DriverID         string       `json:"driver_id"`
	DriverName       string       `json:"driver_name"`
	VehicleID        string       `json:"vehicle_id"`
	VehicleName      string       `json:"vehicle_name"`
	VehicleType      string       `json:"vehicle_type"`
	VehicleRegNum    string       `json:"vehicle_number"`
	ETA              string       `json:"eta"`
	ETAEpoch         int64        `json:"eta_epoch_time"`
	DistanceInKM     float64      `json:"distance"`
	SequenceNo       int          `json:"sequence_no"`
	SectionNo        int          `json:"section_no"`
	RRChangedType    RRChangeType `json:"rr_changed_type"`
	RRChangedISOTime string       `json:"rr_changed_iso_time"`
	/// below fields are for user story: DSPTCH-2027
	NewPickupAddress   *AddressInfo           `json:"rr_pickup_address,omitempty"`
	NewDeliveryAddress *AddressInfo           `json:"rr_delivery_address,omitempty"`
	PickupTimeWindow   *ServiceCommitTimeInfo `json:"rr_pickup_commit_timewindow,omitempty"`
	DeliveryTimeWindow *ServiceCommitTimeInfo `json:"rr_delivery_commit_timewindow,omitempty"`
}

// AddressInfo describe address information in detail,
type AddressInfo struct {
	GPSLat    float64 `json:"latitude"`
	GPSLon    float64 `json:"longitude"`
	Postcode  string  `json:"postcode"`
	Country   string  `json:"country"`
	CountryID string  `json:"country_id"`
	State     string  `json:"state"`
	StateID   string  `json:"state_id"`
	City      string  `json:"city"`
	CityID    string  `json:"city_id"`
	Address   string  `json:"address"`
	Address2  string  `json:"address_2,omitempty"`
	Address3  string  `json:"address_3,omitempty"`
}

// ServiceCommitTimeInfo describe commmit time window and related info (either pickup or delivery)
type ServiceCommitTimeInfo struct {
	JobDate           string `json:"job_date"`        // format: DD/MM/YYY
	StartTime         string `json:"start_time"`      // format: HH:ii:ss
	StartISOTime      string `json:"start_iso_time"`  // format: 2006-01-02T15:04:05-07:00
	StartTimeUnixsec  int64  `json:"start_epoch"`     // format: unix time, at time.Location (based on timezone)
	CommitTime        string `json:"commit_time"`     // format: HH:ii:ss
	CommitISOTime     string `json:"commit_iso_time"` // format: 2006-01-02T15:04:05-07:00
	CommitTimeUnixsec int64  `json:"commit_epoch"`    // format: unix time, at time.Location (based on timezone)
}

// HubLocationInfo describe job From/To Location detail
type HubLocationInfo struct {
	// PartnerID      string  `json:"partner_id" dynamodbav:"partner_id"`
	// PartnerName    string  `json:"partner_name" dynamodbav:"partner_name"`
	FacilityID       string            `json:"id" dynamodbav:"id"`
	FacilityName     string            `json:"facility_name" dynamodbav:"facility_name"`
	CityID           string            `json:"city_id" dynamodbav:"city_id"`
	CityName         string            `json:"city_name" dynamodbav:"city_name"`
	StateID          string            `json:"state_id" dynamodbav:"state_id"`
	StateName        string            `json:"state_name" dynamodbav:"state_name"`
	CountryID        string            `json:"country_id" dynamodbav:"country_id"`
	CountryName      string            `json:"country" dynamodbav:"country"`
	Address          string            `json:"address" dynamodbav:"address"`
	Postcode         string            `json:"postcode" dynamodbav:"postcode"`
	StartLatitude    float64           `json:"start_latitude" dynamodbav:"start_latitude"`
	StartLongitude   float64           `json:"start_longitude" dynamodbav:"start_longitude"`
	DepotLatitude    float64           `json:"depot_latitude" dynamodbav:"depot_latitude"`
	DepotLongitude   float64           `json:"depot_longitude" dynamodbav:"depot_longitude"`
	EndLatitude      float64           `json:"end_latitude" dynamodbav:"end_latitude"`
	EndLongitude     float64           `json:"end_longitude" dynamodbav:"end_longitude"`
	IANATimezone     string            `json:"iana_timezone" dynamodbav:"iana_timezone"`
	UTCTimezone      string            `json:"utc_timezone" dynamodbav:"utc_timezone"`
	ContactNumber    string            `json:"contact_number"`
	ContactPhoneCode *ContactPhoneCode `json:"contact_phonecode,omitempty"`
	ContactName      string            `json:"contact_name,omitempty"`
	Email            string            `json:"email,omitempty"`
}

// ContactPhoneCode describe contact info of location
type ContactPhoneCode struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	LocationCode string `json:"location_code"`
	PhoneCode    string `json:"phonecode"`
	MinLength    int64  `json:"min_length"`
}

// GetOrderLabel returns the order label
func (job *Job) GetOrderLabel() string {
	orderLabel := ""

	if job.ClientOrderCode != "" {
		orderLabel = job.ClientOrderCode
	}

	if job.RefOrderLabel != "" {
		orderLabel = job.RefOrderLabel
	}

	if job.ReferenceID != "" {
		orderLabel = job.ReferenceID
	}

	return orderLabel
}

// GetShipmentLabel to return shipment ID human readable (in their own system definition)
func (job *Job) GetShipmentLabel() string {
	var shipLabel string

	// RefShipmentLabel is obtain/async by OrderModule/Shipment Async
	if job.RefShipmentLabel != "" {
		shipLabel = job.RefShipmentLabel
	}

	return shipLabel
}

// GetPackageLabel returns the package label
func (job *Job) GetPackageLabel() string {
	packageLabel := ""
	if job.PackageCode != "" {
		packageLabel = job.PackageCode
	}

	if job.RefPackageLabel != "" {
		packageLabel = job.RefPackageLabel
	}

	return packageLabel
}

// IsDone returns true if the given job is done even it's failed
func (job Job) IsDone() bool {
	switch job.Status {
	case StatusCompleted,
		StatusForceCompleted,
		StatusFailed:
		return true
	}

	return false
}

// GetCommitTimeTS returns commit time based on service type of the job
func (job Job) GetCommitTimeTS() int64 {
	var commitTime int64
	switch job.ServiceType {
	case ServiceTypePickUp:
		commitTime = job.PickupCommitTimeUnixsec
	case ServiceTypeDelivery:
		commitTime = job.DeliveryCommitTimeUnixsec
	}

	return commitTime
}

// IsNotUnoptimized Validate Unoptimised job
func (job Job) IsNotUnoptimized() bool {
	if job.OrgID == "" || job.RefShipmentID == "" || job.RefPackageID == "" {
		return true
	}
	if !job.IsUnoptimizedJob || job.UnoptimizedJobInfo == nil {
		return true
	}
	if job.UnoptimizedJobInfo.VehicleID == "" || job.UnoptimizedJobInfo.PartnerID == "" || job.UnoptimizedJobInfo.DriverID == "" || job.UnoptimizedJobInfo.JobETA == "" {
		return true
	}

	return false
}

// ShipmentList represents a list of jobs
type ShipmentList []*Job

// NOTE: these func are for sort.Sort use case, please do not remove them
func (sl ShipmentList) Len() int {
	return len(sl)
}

func (sl ShipmentList) Less(i, j int) bool {
	if sl[i].SequenceNo == sl[j].SequenceNo {
		return sl[i].SectionNumber < sl[j].SectionNumber
	}
	return sl[i].SequenceNo < sl[j].SequenceNo
}

func (sl ShipmentList) Swap(i, j int) {
	sl[i], sl[j] = sl[j], sl[i]
}

// OrderPayload contains the data of the original requested payload from Shipment
type OrderPayload struct {
	ItemType           ItemType       `json:"entity_type"`
	OrgID              string         `json:"organisation_id"`
	ServiceType        string         `json:"service_type"`
	OrderID            string         `json:"order_id"`
	OrderLabel         string         `json:"order_id_label"`
	BookedAt           string         `json:"booked_at"`
	ConsigneeInfo      *ConsigneeInfo `json:"consignee_info"`
	DeliveryAddress    string         `json:"delivery_address"`
	DeliveryCity       string         `json:"delivery_city"`
	DeliveryCityID     string         `json:"delivery_city_id"`
	DeliveryCommitTime string         `json:"delivery_commit_time"`
	DeliveryCountry    string         `json:"delivery_country"`
	DeliveryCountryID  string         `json:"delivery_country_id"`
	DeliveryPostcode   string         `json:"delivery_postcode"`
	// DeliveryServiceTime int64          `json:"delivery_service_time"` -> Not available
	DeliveryStartTime string  `json:"delivery_start_time"`
	DeliveryState     string  `json:"delivery_state"`
	DeliveryStateID   string  `json:"delivery_state_id"`
	DeliveryTimezone  string  `json:"delivery_timezone"`
	ID                string  `json:"id"`
	MaxSLA            float64 `json:"max_sla"`
	MinSLA            float64 `json:"min_sla"`
	OriginAddress     string  `json:"origin_address"`
	OriginCity        string  `json:"origin_city"`
	OriginCityID      string  `json:"origin_city_id"`
	OriginCountry     string  `json:"origin_country"`
	OriginCountryID   string  `json:"origin_country_id"`
	OriginPostcode    string  `json:"origin_postcode"`
	OriginState       string  `json:"origin_state"`
	OriginStateID     string  `json:"origin_state_id"`

	PackageNum       int32                `json:"no_of_package"`
	Packages         []*Package           `json:"packages"`
	PartnerID        string               `json:"partner_id"`
	PickupCommitTime string               `json:"pickup_commit_time"`
	PickupStartTime  string               `json:"pickup_start_time"`
	PickupTimezone   string               `json:"pickup_timezone"`
	ShipmentCode     string               `json:"shipment_code"`
	AllocationMatrix []*RouteFacilityInfo `json:"shipment_path"`
	ShipperInfo      *ShipperInfo         `json:"shipper_info"`

	BookingMode        string               `json:"booking_mode"`
	ReferenceID        *string              `json:"reference_id"`
	Attachments        []*Attachment        `json:"attachments,omitempty"`
	CustomerReferences []*CustomerReference `json:"customer_references,omitempty"`
	TagList            []string             `json:"tag_list,omitempty"`

	PickupInstruction   *string       `json:"pickup_instruction"`
	DeliveryInstruction *string       `json:"delivery_instruction"`
	OperationType       *string       `json:"ops_type"` // set it is pointer is for legacy data has no operation type
	HandleFacilityID    *string       `json:"handle_facility_id"`
	HandleFacilityInfo  *FacilityInfo `json:"handle_facility_info,omitempty"`
	Note                *string       `json:"note"`

	// Consolidated shipment related fields
	RefConsolidatedID       string              `json:"ref_consolidated_shipment_id"`
	IsConsolidated          bool                `json:"is_consolidated,omitempty"`
	ConsolidatedOrigin      *FacilityInfo       `json:"consolidated_origin,omitempty"`
	ConsolidatedDestination *FacilityInfo       `json:"consolidated_destination,omitempty"`
	ConsolidatedItemList    []*ConsolidatedItem `json:"consolidated_items,omitempty"`

	MeasurementUnits    *MeasurementUnitInfo `json:"measurement_units,omitempty"`
	ShipmentServiceType string               `json:"shipment_service_type,omitempty"`

	Origin      *OrderAddress `json:"origin"`
	Destination *OrderAddress `json:"destination"`

	PricingInfo *PricingInfo `json:"pricing_info"`

	IsManualSetRun  *bool            `json:"is_manual_set_run,omitempty"`
	OrderSetRunInfo *OrderSetRunInfo `json:"order_set_run_info,omitempty"`

	CreatedBy string `json:"created_by"`
}

// OrderSetRunInfo describes manual set run data input by user from OM
type OrderSetRunInfo struct {
	ID              string `json:"id"`
	PartnerID       string `json:"partner_id"`
	DriverID        string `json:"driver_id"`
	VehicleID       string `json:"vehicle_id"`
	PickupETA       string `json:"pickup_eta"`
	DeliveryETA     string `json:"delivery_eta"`
	SegmentDistance string `json:"segment_distance"`
}

// Attachment represents the attachment(s) in the order
type Attachment struct {
	ID         string  `json:"id"`
	FileURL    string  `json:"file_url"`
	FileName   string  `json:"file_name"`
	FileSize   float64 `json:"file_size"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
	RefreshURL string  `json:"refresh_url,omitempty"`
}

// CustomerReference contains data of a customer reference
type CustomerReference struct {
	ID        string `json:"id"`
	IDLabel   string `json:"id_label"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// OrderAddress describe shipper/consignee address information
type OrderAddress struct {
	ID           string `json:"id"`
	PostalCode   string `json:"postal_code"`
	CityID       string `json:"city_id"`
	City         string `json:"city"`
	StateID      string `json:"state_id"`
	State        string `json:"state"`
	CountryID    string `json:"country_id"`
	Country      string `json:"country"`
	Address      string `json:"address"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	AddressLine3 string `json:"address_line3"`
	// NOTE: set it as string is due to legacy code save order payload as string type
	// json.Unmarshal json.Number with empty string will get error (since GoLang ver. 1.14)
	Lat              string `json:"latitude"`
	Long             string `json:"longitude"`
	ManualCoordinate bool   `json:"manual_coordinates"`
}

// PricingInfo describe order pricing related information
type PricingInfo struct {
	ID                string        `json:"id"`
	COD               float64       `json:"cod"`
	TAX               float64       `json:"tax"`
	Total             float64       `json:"total"`
	CurrencyDetail    *CurrencyInfo `json:"currency"`
	Discount          float64       `json:"discount"`
	Surcharge         float64       `json:"surcharge"`
	BaseTariff        float64       `json:"base_tariff"`
	CurrencyCode      string        `json:"currency_code"`
	InsuranceCharge   float64       `json:"insurance_charge"`
	ServiceTypeCharge float64       `json:"service_type_charge"`
}

// CurrencyInfo describe currency detail information
type CurrencyInfo struct {
	ID                string  `json:"id"`
	Code              string  `json:"code"`
	Name              string  `json:"name"`
	Deleted           bool    `json:"deleted"`
	CreateAt          string  `json:"created_at"` // format: 2021-09-14T09:10:26.897Z
	UpdateAt          string  `json:"updated_at"`
	ExchangeRate      float64 `json:"exchange_rate"`
	IsDefaultCurrency bool    `json:"is_default_currency"`
	OrgID             *string `json:"organisation_id"`
}

// FacilityInChargePeopleInfo describe facility person in charge contact information
// NOTE: it is duplicated is due to it been used and applied as api payload, but the root cause is due to consiginee/shipeer Info on the field PhoneCode using differnt naming but same struct
type FacilityInChargePeopleInfo struct {
	PhoneCode *PhoneCode `json:"consignee_phone_code_id"`
	Name      string     `json:"name"`
	Phone     string     `json:"phone"`
	Email     string     `json:"email,omitempty"`
}

// ConsigneeInfo contains information about a person
type ConsigneeInfo struct {
	PhoneCode *PhoneCode `json:"consignee_phone_code_id"`
	Name      string     `json:"name"`
	Phone     string     `json:"phone"`
	Email     string     `json:"email,omitempty"`
}

// ShipperInfo contains information about a person
type ShipperInfo struct {
	PhoneCode *PhoneCode `json:"shipper_phone_code_id"`
	Name      string     `json:"name"`
	Phone     string     `json:"phone"`
	Email     string     `json:"email,omitempty"`
}

// PhoneCode describe Phone Country Code
type PhoneCode struct {
	Code        string `json:"code"`
	ID          string `json:"id"`
	CountryName string `json:"name"`
}

// Package describes a package
type Package struct {
	ID            string         `json:"id"`
	Code          string         `json:"code"`
	Commodity     string         `json:"commodity"`
	GrossWeight   float64        `json:"gross_weight"`
	Height        float64        `json:"height"`
	ItemCount     int64          `json:"item_count"`
	InsuranceType *string        `json:"insurance_type"`
	Lat           json.Number    `json:"lat"`
	Long          json.Number    `json:"lon"`
	Length        float64        `json:"length"`
	PackageType   string         `json:"package_type"`
	PackagingType *PackagingType `json:"packaging_type"`
	Width         float64        `json:"width"`
	Note          *string        `json:"note"`
	Tags          []string       `json:"tags_list"`
}

// RouteFacilityInfo describes a node in an allocation matrix
type RouteFacilityInfo struct {
	HubID             string `json:"hub_id"`
	MaxSLA            int64  `json:"max_sla"`
	MinSLA            int64  `json:"min_sla"`
	PartnerID         string `json:"partner_id"`
	Position          int32  `json:"position"`
	TransportType     string `json:"transport_type"`
	TransportCategory string `json:"transport_category"` // Air, Ground
}

// GetServiceLatLong returns the service lat, long depending on service type
func (job *Job) GetServiceLatLong() (lat, long float64) {
	lat, long = job.PickupGPSLat, job.PickupGPSLon
	if job.ServiceType == ServiceTypeDelivery {
		lat, long = job.DeliveryGPSLat, job.DeliveryGPSLon
	}
	return lat, long
}

// GetServiceConsignee returns the consignee depending on service type
func (job *Job) GetServiceConsignee() string {
	consignee := job.ShipperLine
	if job.ServiceType == ServiceTypeDelivery {
		consignee = job.RecipientLine
	}
	return consignee
}

// ETAStatusTTLRecord contains the details about the TTL record for target ETA status event
type ETAStatusTTLRecord struct {
	OrgID      string
	JobID      string
	RecordType ETAStatusTTLRecordType
	CreatedAt  int64 // unix timestamp
	StartAt    int64 // unix timestamp
	DispatchID string
}

// ToJSON returns the encoded string for this struct
func (r ETAStatusTTLRecord) ToJSON() string {
	data, _ := json.Marshal(r)
	return string(data)
}

// JobFailureReport contains reporting details of job failure
type JobFailureReport struct {
	Reason      string        `json:"reason"`
	Remarks     string        `json:"remarks"`
	PhotoProofs []*PhotoProof `json:"photo_proofs"`
}

type OnDemandStatus string

const (
	NewAssignedODS OnDemandStatus = "new"
	AcceptedODS    OnDemandStatus = "accepted"
	RejectedODS    OnDemandStatus = "rejected"
	TimeoutODS     OnDemandStatus = "timeout"
)

var (
	StrConvertOnDemandStatus = map[string]OnDemandStatus{
		"new":      NewAssignedODS,
		"accepted": AcceptedODS,
		"rejected": RejectedODS,
		"timeout":  TimeoutODS,
	}
)

// ODSMetric shows ODS response or timeout metrics
type ODSMetric struct {
	Rejected int64 `json:"rejected"`
	Timeout  int64 `json:"timeout"`
}

// ODSDriverHistory is the assignment and responding log to a shipment order,
// map key string is "driverID"
type ODSDriverHistory map[string]*ODSAssignment

type ODSAssignment struct {
	DriverID         string         `json:"driver_id"`
	AssignedUnixtime int64          `json:"assigned_time_ts"`
	AssignedISOTime  string         `json:"assigned_time"`
	ExpiredUnixtime  int64          `json:"expired_time_ts"`
	ExpiredISOTime   string         `json:"expired_time"`
	Timezone         string         `json:"timezone"`
	ResponseStatus   OnDemandStatus `json:"response_status"`
	ResponseUnixtime int64          `json:"response_time_ts"`
	ResponseISOTime  string         `json:"response_time"`
}

// ConsolidatedItemType represents the consolidated item type
// It should be either shipment or package
type ConsolidatedItemType string

const (
	// ConsolidatedItemTypeShipment - shipment type
	ConsolidatedItemTypeShipment ConsolidatedItemType = "shipment"
	// ConsolidatedItemTypePackage - package type
	ConsolidatedItemTypePackage ConsolidatedItemType = "package"
)

// ConsolidatedItem contains information of an item
// in a consolidated shipment
type ConsolidatedItem struct {
	ID         string               `json:"id"`                    // UUID
	Label      string               `json:"label"`                 // Human readable ID (shipment/package label for shipment/package item type)
	Type       ConsolidatedItemType `json:"type"`                  // Available types: (shipment,package)
	PackageIDs []string             `json:"package_ids,omitempty"` // Packages in a shipment, empty if item type is package
}

// Dimensions3D contains 3D dimensions size & unit
type Dimensions3D struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Length float64 `json:"length"`
	Unit   string  `json:"metric"`
}

// Dimensions3DOptional contains 3D dimensions size & unit (optional)
// The dimensions size could be empty string
type Dimensions3DOptional struct {
	Width  NumberStr `json:"width"`
	Height NumberStr `json:"height"`
	Length NumberStr `json:"length"`
	Unit   string    `json:"metric"`
}

// Dimensions2DOptional contains 2D dimensions size & unit (optional)
// The dimensions size could be empty string
type Dimensions2DOptional struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Unit   string  `json:"metric"`
}

type PackagingType struct {
	ID                  string                `json:"id"`
	Name                string                `json:"name"`
	OrgID               string                `json:"organisation_id"`
	Description         string                `json:"description"`
	IsOrderPackaging    bool                  `json:"order_packaging"`
	IsShipmentPackaging bool                  `json:"shipment_packaging"`
	Dimensions          *Dimensions3D         `json:"dimension,omitempty"`          // Units: Centimeters,Milimeter,Meter,Inches
	InternalDimensions  *Dimensions3DOptional `json:"internal_dimension,omitempty"` // Units: Centimeters,Milimeter,Meter,Inches
	DoorAperture        *Dimensions2DOptional `json:"door_aperture,omitempty"`      // Units: Centimeters,Milimeter,Meter,Inches
	MaxGrossWeight      float64               `json:"max_gross_weight"`
	MaxGrossWeightUnit  string                `json:"max_gross_weight_unit"` // Kilogram,Pound
	MaxNetWeight        float64               `json:"max_net_weight"`
	MaxNetWeightUnit    string                `json:"max_net_weight_unit"` // Kilogram,Pound
	TareWeight          float64               `json:"tare_weight"`
	TareWeightUnit      string                `json:"tare_weight_unit"` // Kilogram,Pound
	InternalVolume      float64               `json:"internal_volume"`
	InternalVolumeUnit  string                `json:"internal_volume_unit"` // Cubic Meter,Liter,Cubic Feet
	DeletedAt           string                `json:"deleted_at"`
	SyncAt              string                `json:"sync_at"`
	CreatedAt           string                `json:"created_at"`
	UpdatedAt           string                `json:"updated_at"`
}

type NumberStr string

// UnmarshalJSON customizes json.Unmarshal method for NumberStr
func (n *NumberStr) UnmarshalJSON(data []byte) error {
	var numStr string
	err := json.Unmarshal(data, &numStr)
	if err != nil {
		numStr = string(data)
	}
	numStr = strings.TrimSpace(numStr)
	if numStr != "" {
		_, parseFloatErr := strconv.ParseFloat(numStr, 64)
		if parseFloatErr != nil {
			return parseFloatErr
		}
		*n = NumberStr(numStr)
	}
	return nil
}

// MarshalJSON customizes json.Marshal method for NumberStr
func (n NumberStr) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(n))
}

func (n NumberStr) String() string {
	return string(n)
}

// Float64 returns the float 64 bit value of the number
func (n NumberStr) Float64() (num float64, err error) {
	return strconv.ParseFloat(string(n), 64)
}

// PhotoProof represents a photo proof object (filename, uploaded url, ...)
type PhotoProof struct {
	Filename               string `json:"filename"`
	FileCaptureTimeUnixsec int64  `json:"file_capture_time_unixsec"`
	FileUploadTimeUnixsec  int64  `json:"file_upload_time_unixsec"`
	PresignedURL           string `json:"presigned_url"`
}

// ETAStatus is used to track if a job is going to delayed based on commit time and ETA time
type ETAStatus int32

const (
	ETAStatus_ETA_STATUS_UNKNOWN        ETAStatus = 0
	ETAStatus_ETA_STATUS_NOT_APPLICABLE ETAStatus = 1
	ETAStatus_ETA_STATUS_ON_TIME        ETAStatus = 2
	ETAStatus_ETA_STATUS_ETA_ELAPSED    ETAStatus = 3
	ETAStatus_ETA_STATUS_DELAYED        ETAStatus = 4
)

// Enum value maps for ETAStatus.
var (
	ETAStatus_name = map[int32]string{
		0: "ETA_STATUS_UNKNOWN",
		1: "ETA_STATUS_NOT_APPLICABLE",
		2: "ETA_STATUS_ON_TIME",
		3: "ETA_STATUS_ETA_ELAPSED",
		4: "ETA_STATUS_DELAYED",
	}
	ETAStatus_value = map[string]int32{
		"ETA_STATUS_UNKNOWN":        0,
		"ETA_STATUS_NOT_APPLICABLE": 1,
		"ETA_STATUS_ON_TIME":        2,
		"ETA_STATUS_ETA_ELAPSED":    3,
		"ETA_STATUS_DELAYED":        4,
	}
)

func (x ETAStatus) Enum() *ETAStatus {
	p := new(ETAStatus)
	*p = x
	return p
}

// PODAssets contains information about the POD assets of a job
type PODAssets struct {
	Document  []PODItem `json:"document"`
	Signature []PODItem `json:"signature"`
	JobIDs    []string  `json:"job_ids"`
}

// PODItem represents an POD item
type PODItem struct {
	Filename     string  `json:"filename"`
	CapturedAt   string  `json:"captured_at"` // ISO time string "2006-01-02T15:04:05-07:00"
	Type         DocType `json:"type"`
	UploadedAt   string  `json:"uploaded_at"` // ISO time string "2006-01-02T15:04:05-07:00"
	PresignedURL string  `json:"presigned_url,omitempty"`
}

// DocType represents the document type
type DocType string

const (
	// DocTypePOD represent the POD image type
	DocTypePOD DocType = "document"
	// DocTypePODSignature represents the POD signature image type
	DocTypePODSignature DocType = "signature"
	// DocTypePackageReportPhotoProof represents the package report photo proof type
	DocTypePackageReportPhotoProof DocType = "package_report_photo_proof"
	// DocTypeJobFailureReport represents the image for the job failure report
	DocTypeJobFailureReport DocType = "job_failure_report"
)

// InboundScanResult represent the inbound scan result
type InboundScanResult struct {
	ScannedText string `json:"scanned_text"`
	// Matched determines whether the scanned text matches the job details or not
	Matched    bool    `json:"matched"`
	ScannedAt  string  `json:"scanned_at"`  // ISO time string
	VerifiedAt string  `json:"verified_at"` // ISO time string
	Lat        float64 `json:"lat"`
	Long       float64 `json:"long"`
}

// UpdateInboundScanResultInput represent the inbound scan update input
type UpdateInboundScanResultInput struct {
	Job         *Job
	ScannedText string
	ScannedAt   string // ISO time string
	Lat         float64
	Long        float64
}

// InboundScanFailedItem represent the inbound scan update input
type InboundScanFailedItem struct {
	JobID  string `json:"job_id"`
	Reason string `json:"reason"`
}

// PODScanResult represent the POD scan result
type PODScanResult struct {
	ScannedText string `json:"scanned_text"`
	// Matched determines whether the scanned text matches the job details or not
	Matched    bool    `json:"matched"`
	ScannedAt  string  `json:"scanned_at"`  // ISO time string
	VerifiedAt string  `json:"verified_at"` // ISO time string
	Lat        float64 `json:"lat"`
	Long       float64 `json:"long"`
}

// UpdatePODScanResultInput represent the POD scan update input
type UpdatePODScanResultInput struct {
	Job         *Job
	ScannedText string
	ScannedAt   string // ISO time string
	Lat         float64
	Long        float64
}

// PODScanFailedItem represent the POD scan update input
type PODScanFailedItem struct {
	JobID  string `json:"job_id"`
	Reason string `json:"reason"`
}

// PackageReportType represents types of package reports
type PackageReportType string

// current package report types
const (
	PackageReportTypeDamaged PackageReportType = "damaged"
	PackageReportTypeLost    PackageReportType = "lost"
)

// PackageReportReason represents reasons of package report type "damaged"
type PackageReportReason string

// current reasons of package report type "damaged"
const (
	PackageReportReasonWaterDamage PackageReportReason = "water_damage"
	PackageReportReasonDented      PackageReportReason = "dented"
	PackageReportReasonOpened      PackageReportReason = "opened"
)

// PackageState represents current package state at report time
type PackageState string

// current states for package report
const (
	PackageStateAcceptable   PackageState = "Acceptable"
	PackageStateUnacceptable PackageState = "Unacceptable"
)

// PackageReport contains reporting details of job package status
type PackageReport struct {
	State                 PackageState        `json:"state"`
	ReportType            PackageReportType   `json:"report_type"`
	Reason                PackageReportReason `json:"reason"`
	Remarks               string              `json:"remarks"`
	PhotoProofs           []*PhotoProof       `json:"photo_proofs"`
	ReportedTimeInUnixSec int64               `json:"reported_time_in_unix_sec"`
}

// RoutingType dispatch of route type can take
type RoutingType string

const (
	// NotDefined -
	NotDefined RoutingType = ""
	// StraightLine means operation type is StraightLine for drones
	StraightLine RoutingType = "straightLine"
	// HubAndSpoke means operation type is Hub
	HubAndSpoke RoutingType = "service"
	// Point2Point means operation type is P2P
	Point2Point RoutingType = "shipment"
)

var (
	// RoutingTypeName mapping of routing type from order
	RoutingTypeName = map[string]RoutingType{
		"StraightLine": StraightLine,
		"Hub":          HubAndSpoke,
		"P2P":          Point2Point,
	}
)

// ItemType is to define data item type in dynamodb
// the field name of item in DynamoDB is "entity_type"
type ItemType string

const (
	UnknownType       ItemType = ""
	JobDriverType     ItemType = "jobdriver" // Dispatch Job Driver Entity type in dynamoDB
	JobType           ItemType = "job"       // Dispatch Job Entity type in dynamoDB
	ShipmentType      ItemType = "shipment"  // Shipment entity type in dynamoDB
	PubSubType        ItemType = "pubsub"
	ODSTTLType        ItemType = "odsttl"        // ODS Driver Assignment TTL entity item
	ShipmentOrderType ItemType = "shipmentOrder" // Order from SHP

	RRMission         ItemType = "rr_mission"
	RRMissionRoute    ItemType = "rr_mission_route"
	RRMissionJob      ItemType = "rr_mission_job"
	RRMissionETAParam ItemType = "rr_mission_eta_param"

	// ETAStatusTTLRecordType helps us detect the record for TTL record model in DynamoDB
	//TAStatusTTLRecordType ItemType = "eta_status_ttl_record"
)

// ETAStatusTTLRecordType describes the event of TTL record
type ETAStatusTTLRecordType int32

var (
	StrConvertItemType = map[string]ItemType{
		"jobdriver":            JobDriverType,
		"job":                  JobType,
		"shipment":             ShipmentType,
		"odsttl":               ODSTTLType,
		"rr_mission":           RRMission,
		"rr_mission_route":     RRMissionRoute,
		"rr_mission_job":       RRMissionJob,
		"rr_mission_eta_param": RRMissionETAParam,
		//"eta_status_ttl_record": ETAStatusTTLRecordType,
		"shipmentOrder": ShipmentOrderType,
	}
)

// TriggerCode defines default trigger code type
type TriggerCode string

const (
	// TriggerCodeUnknown unknown trigger code
	TriggerCodeUnknown TriggerCode = ""

	// DispatchCreateDriverAssignedToRoute Dispatch created and scheduled (dispatch status = "scheduled") which meaning driver assign to a route
	DispatchCreateDriverAssignedToRoute TriggerCode = "1601"

	// ReSequenceAndReassignmentDriverReAssignedToRoute driver re-assigned to route through R&R
	ReSequenceAndReassignmentDriverReAssignedToRoute TriggerCode = "1602"

	// ReSequenceAndReassignmentDriverRemoveFromRoute driver removed from route through R&R
	ReSequenceAndReassignmentDriverRemoveFromRoute TriggerCode = "1603"

	//ReSequenceAndReassignmentDriverJobRescheduled driver's job was rescheduled
	ReSequenceAndReassignmentDriverJobRescheduled TriggerCode = "1607"

	// ReSequenceAndReassignmentDriverAssignedToRoute driver assigned to route through R&R or whose route added current route throught R&R
	ReSequenceAndReassignmentDriverAssignedToRoute TriggerCode = "1608"
)
