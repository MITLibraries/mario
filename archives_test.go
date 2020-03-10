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

	// Test Citation field
	if record.Citation != "Kevin Lynch papers, MC 208, box X. Massachusetts Institute of Technology, Department of Distinctive Collections, Cambridge, Massachusetts." {
		t.Error("Expected match, got", record.Citation)
	}

	// Test ContentType field
	if record.ContentType != "Archival collection" {
		t.Error("Expected match, got", record.ContentType)
	}

	// Test Contributor field
	if record.Contributor[0].Value != "Lynch, Kevin, 1918-1984" {
		t.Error("Expected match, got", record.Contributor[0].Value)
	}

	if record.Contributor[0].Kind != "Person" {
		t.Error("Expected match, got", record.Contributor[0].Kind)
	}

	// Test Holdings field
	if record.Holdings[0].Location != "Materials are stored off-site. Advance notice is required for use." {
		t.Error("Expected match, got", record.Holdings[0].Location)
	}

	// Test Identifier field
	if record.Identifier != "MIT:archivesspace:MC.0208" {
		t.Error("Expected match, got", record.Identifier)
	}

	// Test Links field
	if len(record.Links) != 234 {
		t.Error("Expected 234, got", len(record.Links))
	}

	// Test PublicationDate field
	if record.PublicationDate != "1934-1988" {
		t.Error("Expected match, got", record.PublicationDate)
	}

	// Test SourceLink field
	if record.SourceLink != "https://archivesspace.mit.edu/repositories/2/resources/739" {
		t.Error("Expected match, got", record.SourceLink)
	}

	// Test Subject field
	if len(record.Subject) != 8 {
		t.Error("Expected match, got", len(record.Subject))
	}

	// Test Title field
	if record.Title != "Kevin Lynch papers" {
		t.Error("Expected match, got", record.Title)
	}
}
