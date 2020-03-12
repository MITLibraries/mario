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
		t.Error("Expected 'Kevin Lynch papers, MC 208, box X. Massachusetts Institute of Technology, Department of Distinctive Collections, Cambridge, Massachusetts.', got", record.Citation)
	}

	// Test ContentType field
	if record.ContentType != "Archival collection" {
		t.Error("Expected 'Archival collection', got", record.ContentType)
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
		t.Error("Expected 'Materials are stored off-site. Advance notice is required for use.', got", record.Holdings[0].Location)
	}

	// Test Identifier field
	if record.Identifier != "MIT:archivesspace:MC.0208" {
		t.Error("Expected 'MIT:archivesspace:MC.0208', got", record.Identifier)
	}

	// Test Language field
	if len(record.Language) != 1 {
		t.Error("Expected 1, got", len(record.Language))
	}

	if record.Language[0] != "English" {
		t.Error("Expected 'English', got", record.Language[0])
	}

	// Test Links field
	if len(record.Links) != 234 {
		t.Error("Expected 234, got", len(record.Links))
	}

	if record.Links[0].URL != "http://hdl.handle.net/1721.3/35646" {
		t.Error("Expected 'http://hdl.handle.net/1721.3/35646', got", record.Links[0].URL)
	}

	if record.Links[0].Text != "K.L. 3-8-55: 1955 March 8" {
		t.Error("Expected 'K.L. 3-8-55: 1955 March 8', got", record.Links[0].Text)
	}

	if record.Links[0].Kind != "Digital object" {
		t.Error("Expected 'Digital object', got", record.Links[0].Kind)
	}

	//Test PhysicalDescription field
	if record.PhysicalDescription != "16.5 Cubic Feet; 12 record cartons, 8 manuscript boxes, 1 half manuscript box, 4 medium flat boxes, 1 large flat box, 3 small media boxes, 1 slide box and 2 loose drawings, 10 oversize folders" {
		t.Error("Expected '16.5 Cubic Feet; 12 record cartons, 8 manuscript boxes, 1 half manuscript box, 4 medium flat boxes, 1 large flat box, 3 small media boxes, 1 slide box and 2 loose drawings, 10 oversize folders', got", record.PublicationDate)
	}

	// Test PublicationDate field
	if record.PublicationDate != "1934-1988" {
		t.Error("Expected '1934-1988', got", record.PublicationDate)
	}

	// Test SourceLink field
	if record.SourceLink != "https://archivesspace.mit.edu/repositories/2/resources/739" {
		t.Error("Expected 'https://archivesspace.mit.edu/repositories/2/resources/739', got", record.SourceLink)
	}

	// Test Subject field
	if len(record.Subject) != 8 {
		t.Error("Expected 8, got", len(record.Subject))
	}

	if record.Subject[0] != "Urban ecology" {
		t.Error("Expected 'Urban ecology', got", record.Subject[0])
	}

	// Test Title field
	if record.Title != "Kevin Lynch papers" {
		t.Error("Expected match, got", record.Title)
	}
}
