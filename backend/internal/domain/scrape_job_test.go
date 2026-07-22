package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestScrapeJob_SetStatus(t *testing.T) {
	job := NewScrapeJob(uuid.New(), uuid.New())

	if got := job.Status; got != ScrapeJobStatusPending {
		t.Fatalf("NewScrapeJob status = %q, want %q", got, ScrapeJobStatusPending)
	}

	if err := job.SetStatus(ScrapeJobStatusRunning); err != nil {
		t.Fatalf("SetStatus(running) returned error: %v", err)
	}
	if job.Status != ScrapeJobStatusRunning {
		t.Fatalf("Status = %q, want %q", job.Status, ScrapeJobStatusRunning)
	}

	if err := job.SetStatus("bogus"); err == nil {
		t.Fatal("SetStatus(\"bogus\") did not return an error")
	}
	if job.Status != ScrapeJobStatusRunning {
		t.Fatalf("Status changed after rejected SetStatus: got %q", job.Status)
	}
}

func TestScrapeJobStatus_Valid(t *testing.T) {
	valid := []ScrapeJobStatus{ScrapeJobStatusPending, ScrapeJobStatusRunning, ScrapeJobStatusSucceeded, ScrapeJobStatusFailed}
	for _, s := range valid {
		if !s.Valid() {
			t.Errorf("%q.Valid() = false, want true", s)
		}
	}
	if ScrapeJobStatus("unknown").Valid() {
		t.Error(`"unknown".Valid() = true, want false`)
	}
}
