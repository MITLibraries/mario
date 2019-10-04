package main

import (
	"os"
	"testing"
)

func TestArchivesProcessesAllRecords(t *testing.T) {
	ead, err := os.Open("fixtures/aspace_samples.xml")
	if err != nil {
		t.Error(err)
	}

	out := make(chan Record)
	p := archivesparser{file: ead}
	go p.parse(out)

	var chanLength int
	for range out {
		chanLength++
	}

	if chanLength != 11 {
		t.Error("Expected match, got", chanLength)
	}
}

func TestArchivesRecordParsing(t *testing.T) {
	ead, err := os.Open("fixtures/aspace_samples.xml")
	if err != nil {
		t.Error(err)
	}

	out := make(chan Record)
	p := archivesparser{file: ead}
	go p.parse(out)

	record := <-out

	if record.Contributor[0].Value != "Lynch, Kevin, 1918-1984" {
		t.Error("Expected match, got", record.Contributor[0].Value)
	}

	if record.Contributor[0].Kind != "Person" {
		t.Error("Expected match, got", record.Contributor[0].Kind)
	}

	if record.Identifier != "MIT:archivespace:MC.0208" {
		t.Error("Expected match, got", record.Identifier)
	}

	if record.Title != "Kevin Lynch papers" {
		t.Error("Expected match, got", record.Title)
	}

	if len(record.Subject) != 8 {
		t.Error("Expected match, got", len(record.Subject))
	}

	if record.PublicationDate != "1934-1988" {
		t.Error("Expected match, got", record.PublicationDate)
	}

	if len(record.Links) != 234 {
		t.Error("Expected 234, got", len(record.Links))
	}

	if record.Holdings[0].Location != "Materials are stored off-site. Advance notice is required for use." {
		t.Error("Expected match, got", record.Holdings[0].Location)
	}
}
