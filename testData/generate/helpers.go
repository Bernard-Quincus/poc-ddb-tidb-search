package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"poc-ddb-tidb-search/pkg/models"
	"time"
)

type TestData struct {
	Job   *models.Job
	Count int
}

const dataDir = "../../testData/generate/data"

func DeepCopy(dst, src *models.Job) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func randomStatus() models.Status {
	rand.Seed(time.Now().Unix())
	status := []models.Status{
		models.StatusNew,
		models.StatusUnallocated,
		models.StatusDraft,
		models.StatusOnRoute,
		models.StatusCompleted,
		models.StatusForceCompleted,
		models.StatusFailed,
		models.StatusCancelled,
	}
	n := rand.Intn(len(status))
	return status[n]
}

func loadFile(path string) (*models.Job, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		log.Println(err)
	}

	job := new(models.Job)
	err = json.Unmarshal(data, job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func saveToFile(data []byte, name string) error {
	filePath := filepath.Join(dataDir, name)
	return os.WriteFile(filePath, data, 0755)
}
