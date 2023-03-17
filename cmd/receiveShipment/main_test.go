package main

import (
	"os"
	"poc-ddb-tidb-search/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBody(t *testing.T) {
	testPayload, err := os.ReadFile("../../testData/job.json")
	assert.NoError(t, err)

	job, err := parseBody(string(testPayload))
	assert.NoError(t, err)

	// just check key fields
	assert.Equal(t, "SHIPMENT#default_9655995f-c682-44e4-9ebd-888a404b7b15_0e7c5b46-653e-4447-9257-0dd00e333df5_delivery", job.PK)
	assert.Equal(t, "default_9655995f-c682-44e4-9ebd-888a404b7b15_0e7c5b46-653e-4447-9257-0dd00e333df5_delivery", job.ID)
	assert.Equal(t, models.Status("unallocated"), job.Status)
	assert.Equal(t, "POC-TEST-ORGID-001122", job.OrgID2)
	assert.Equal(t, job.ID, job.DocID)
}
