package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"poc-ddb-tidb-search/pkg/models"

	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
)

var (
	recordSent = 0

	// starts at 0, change this to continue where you left off, or you can one shot generate 100K test files
	// this var is used for creating filename eg. Job-1.json in the generate/data folder
	fileCount = 0 //add by +10K
)

func main() {

	jobBase, err := loadFile("../../testData/job.json")
	if err != nil {
		fmt.Println("failed to load file")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startTime := time.Now()
	fmt.Println("Started sending sample data...")

	// will send 10k total test data, 1k per batch, or adjust according to your liking
	for j := 0; j < 10; j++ {
		for res := range send(ctx, generateTestData(ctx, jobBase, 1000)) { // change number of test data to generate
			if res.StatusCode != 200 {
				fmt.Println("got error", res.StatusCode)
				continue
			}

			recordSent++
			fmt.Println("Total record sent: ", recordSent)
		}
		fmt.Println("Let server breathe for 20 secs")
		time.Sleep(20 * time.Second)
	}

	fmt.Println("Elapsed time:", time.Since(startTime))
	fmt.Println("done")
}

func send(ctx context.Context, jobStream <-chan *TestData) <-chan HttpResponse {
	resultStream := make(chan HttpResponse)
	var wg sync.WaitGroup

	sender := func() {
		defer wg.Done()

		client := NewHttpClient()
		for job := range jobStream {
			jobJson, _ := json.MarshalIndent(job.Job, "", "   ")
			res, err := client.Post(jobJson)
			if err != nil {
				fmt.Println("got error while sending request:", err.Error())
			}

			select {
			case <-ctx.Done():
				return
			case resultStream <- res:
				if res.StatusCode == 200 {
					fileCount++
					fname := fmt.Sprintf("Job-%d.json", fileCount)
					saveToFile(jobJson, fname)
				}
			}
		}
	}

	wg.Add(10)
	for j := 0; j < 10; j++ {
		go sender()
	}

	go func() {
		wg.Wait()
		close(resultStream)
	}()

	return resultStream
}

func generateTestData(ctx context.Context, jobBase *models.Job, num int) <-chan *TestData {
	jobStream := make(chan *TestData)

	go func() {
		defer close(jobStream)
		for i := 0; i < num; i++ {
			job := new(models.Job)
			DeepCopy(job, jobBase)
			makeRandomData(job)

			time.Sleep(time.Millisecond * 100) //adjust accordingly, make sure don't spam the server

			select {
			case <-ctx.Done():
				return
			case jobStream <- &TestData{Job: job, Count: i}:
			}
		}
	}()
	return jobStream
}

// makeRandomData mutates job
func makeRandomData(job *models.Job) {
	jobID := uuid.NewString()
	job.ID = jobID
	job.JobID = jobID
	job.SegmentID = uuid.NewString()
	job.RefShipmentID = faker.UUIDDigit()
	job.OrderPayload.OrderID = uuid.NewString()
	job.Status = randomStatus()

	job.OrderPayload.ConsigneeInfo.Name = faker.Name()
	job.OrderPayload.ConsigneeInfo.Email = faker.Email()
	job.OrderPayload.ConsigneeInfo.Phone = faker.Phonenumber()

	job.OrderPayload.ShipperInfo.Name = faker.Name()
	job.OrderPayload.ShipperInfo.Email = faker.Email()
	job.OrderPayload.ShipperInfo.Phone = faker.Phonenumber()

	pickdate := faker.Date()

	layout := "2006-02-01"
	pd, _ := time.Parse(layout, pickdate)
	delDate := pd.AddDate(0, 0, 1) // add 1 day

	job.PickupDate = pd.Format("02/01/2006")        // dd/mm/yyyy
	job.DeliveryDate = delDate.Format("02/01/2006") // dd/mm/yyyy

	job.PickupStartTime = faker.TimeString()
	job.DeliveryStartTime = faker.TimeString()
}
