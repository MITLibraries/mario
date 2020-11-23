package generator

import (
	"os"
	"testing"

	"github.com/mitlibraries/mario/pkg/record"
)

func TestDspaceProcessesAllRecords(t *testing.T) {
	dspace, err := os.Open("../../fixtures/dspace_samples.xml")
	if err != nil {
		t.Error(err)
	}

	out := make(chan record.Record)
	p := dspaceparser{file: dspace}
	go p.parse(out)

	var chanLength int
	for range out {
		chanLength++
	}

	if chanLength != 226 {
		t.Error("Expected 226, got", chanLength)
	}
}

func TestDspaceRecordParsing(t *testing.T) {
	dspace, err := os.Open("../../fixtures/dspace_samples.xml")
	if err != nil {
		t.Error(err)
	}

	out := make(chan record.Record)
	p := dspaceparser{file: dspace}
	go p.parse(out)

	record := <-out

	// Test Citation field
	if record.Citation != "Price, Max et al. \"Fodder, pasture, and the development of complex society in the Chalcolithic: isotopic perspectives on animal husbandry.\" Archaeological and Anthropological Sciences 12, 4 (March 2020): 95 © 2020 Springer-Verlag" {
		t.Error("Expected 'Price, Max et al. \"Fodder, pasture, and the development of complex society in the Chalcolithic: isotopic perspectives on animal husbandry.\" Archaeological and Anthropological Sciences 12, 4 (March 2020): 95 © 2020 Springer-Verlag', got", record.Citation)
	}

	// Test Collection field
	if len(record.Collection) != 1 {
		t.Error("Expected 1, got", record.Collection)
	}
	if record.Collection[0] != "MIT Open Access Articles" {
		t.Error("Expected 'MIT Open Access Articles', got", record.Collection[0])
	}

	// Test ContentType field
	if record.ContentType != "Article" {
		t.Error("Expected 'Article', got", record.ContentType)
	}

	// Test Contributor field
	if len(record.Contributor) != 5 {
		t.Error("Expected 5, got", len(record.Contributor))
	}
	if record.Contributor[0].Kind != "author" {
		t.Error("Expected 'author', got", record.Contributor[0].Kind)
	}
	if record.Contributor[0].Value != "Price, Max D" {
		t.Error("Expected 'Price, Max D', got", record.Contributor[0].Value)
	}

	// Test DOI field
	if len(record.Doi) != 1 {
		t.Error("Expected 1, got", len(record.Doi))
	}
	if record.Doi[0] != "10.2514/6.2020-4181" {
		t.Error("Expected '10.2514/6.2020-4181', got", record.Doi[0])
	}

	// Test Identifier field
	if record.Identifier != "oai:dspace.mit.edu:1721.1-128382" {
		t.Error("Expected 'oai:dspace.mit.edu:1721.1-128382', got", record.Identifier)
	}

	// Test ISBN field
	if len(record.Isbn) != 1 {
		t.Error("Expected 1, got", len(record.Isbn))
	}
	if record.Isbn[0] != "9781624106088" {
		t.Error("Expected '9781624106088', got", record.Isbn[0])
	}

	// Test ISSN field
	if len(record.Issn) != 1 {
		t.Error("Expected 1, got", len(record.Issn))
	}
	if record.Issn[0] != "1866-9557" {
		t.Error("Expected '1866-9557', got", record.Issn[0])
	}

	// Test Language field
	if len(record.Language) != 1 {
		t.Error("Expected 1, got", len(record.Language))
	}
	if record.Language[0] != "en" {
		t.Error("Expected 'en', got", record.Language[0])
	}

	// Test Links field
	if len(record.Links) != 1 {
		t.Error("Expected 1, got", len(record.Links))
	}
	if record.Links[0].URL != "https://hdl.handle.net/1721.1/128382" {
		t.Error("Expected 'https://hdl.handle.net/1721.1/128382', got", record.Links[0].URL)
	}
	if record.Links[0].Text != "Digital object URI" {
		t.Error("Expected 'Digital object URI', got", record.Links[0].Text)
	}
	if record.Links[0].Kind != "Digital object URI" {
		t.Error("Expected 'Digital object URI', got", record.Links[0].Kind)
	}

	// Test OCLC field
	if len(record.OclcNumber) != 1 {
		t.Error("Expected 1, got", len(record.OclcNumber))
	}
	if record.OclcNumber[0] != "1202775217" {
		t.Error("Expected '1202775217', got", record.OclcNumber[0])
	}

	// Test PublicationDate field
	if record.PublicationDate != "2020-03" {
		t.Error("Expected '2020-03', got", record.PublicationDate)
	}

	// Test SourceLink field
	if record.SourceLink != "https://hdl.handle.net/1721.1/128382" {
		t.Error("Expected 'https://hdl.handle.net/1721.1/128382', got", record.SourceLink)
	}

	// Test Subject field
	if len(record.Subject) != 3 {
		t.Error("Expected 3, got", len(record.Subject))
	}
	if record.Subject[0] != "Engineering Systems Division." {
		t.Error("Expected 'Engineering Systems Division.', got", record.Subject[0])
	}

	// Test Summary field
	if record.Summary[0] != "The emergence of social complexity in the Southern Levant during the Chalcolithic (c. 4500–3600 cal. BC) was intimately tied to intensification in animal management. For the first time, secondary products such as milk and wool were intensively exploited, supplying communities with increasingly diverse foodstuffs and raw materials for craft production and exchange, but the precise herding practices underlying these new production strategies are unknown. Here, we explore the role of multi-species livestock pasturing through carbon and nitrogen isotopic analysis of animal bones from Marj Rabba (Har Hasha’avi, West) in the Lower Galilee (ca. 4600–4200 cal. BC). Isotopic results suggest different pasturing/foddering of sheep compared with goats. Cattle were largely pastured locally, but high δ13C values in some animals indicate access to the Jordan River Valley (the Ghor in Arabic), where major Chalcolithic settlements were situated. This may indicate some cattle were moved along regional Chalcolithic exchange networks established for other prestige objects, such as copper. Finally, we provide evidence for moderate 15N enrichment in pigs relative to herbivorous livestock indicates. Possible interpretations include consumption of nuts (esp. acorns), household refuse containing animal protein, and/or fattening pigs on grain. Although an interpretation that requires further exploration, grain foddering of pigs would complement the zooarchaeological data for early slaughter, which suggests intensive meat production at Marj Rabba. It might also help explain why pig husbandry, as a drain on grain stockpiles, was gradually abandoned during the Bronze Age. Taken together, the isotopic and zooarchaeological data indicate an economy in transition from a non-specialized, household-based Neolithic economy to one in which the production of agrarian wealth, including animal secondary products, was beginning to emerge." {
		t.Error("Expected 'The emergence of social complexity in the Southern Levant during the Chalcolithic (c. 4500–3600 cal. BC) was intimately tied to intensification in animal management. For the first time, secondary products such as milk and wool were intensively exploited, supplying communities with increasingly diverse foodstuffs and raw materials for craft production and exchange, but the precise herding practices underlying these new production strategies are unknown. Here, we explore the role of multi-species livestock pasturing through carbon and nitrogen isotopic analysis of animal bones from Marj Rabba (Har Hasha’avi, West) in the Lower Galilee (ca. 4600–4200 cal. BC). Isotopic results suggest different pasturing/foddering of sheep compared with goats. Cattle were largely pastured locally, but high δ13C values in some animals indicate access to the Jordan River Valley (the Ghor in Arabic), where major Chalcolithic settlements were situated. This may indicate some cattle were moved along regional Chalcolithic exchange networks established for other prestige objects, such as copper. Finally, we provide evidence for moderate 15N enrichment in pigs relative to herbivorous livestock indicates. Possible interpretations include consumption of nuts (esp. acorns), household refuse containing animal protein, and/or fattening pigs on grain. Although an interpretation that requires further exploration, grain foddering of pigs would complement the zooarchaeological data for early slaughter, which suggests intensive meat production at Marj Rabba. It might also help explain why pig husbandry, as a drain on grain stockpiles, was gradually abandoned during the Bronze Age. Taken together, the isotopic and zooarchaeological data indicate an economy in transition from a non-specialized, household-based Neolithic economy to one in which the production of agrarian wealth, including animal secondary products, was beginning to emerge.', got", record.Summary[0])
	}

	// Test Title field
	if record.Title != "Fodder, pasture, and the development of complex society in the Chalcolithic: isotopic perspectives on animal husbandry at Marj Rabba" {
		t.Error("Expected 'Fodder, pasture, and the development of complex society in the Chalcolithic: isotopic perspectives on animal husbandry at Marj Rabba', got", record.Title)
	}
}
