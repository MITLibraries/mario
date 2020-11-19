package generator

import (
	"encoding/xml"
	"io"
	"regexp"
	"strings"

	"github.com/mitlibraries/mario/pkg/record"
)

type dspaceparser struct {
	file io.Reader
}

// DspaceGenerator parses dspace oai_dc xml data.
type DspaceGenerator struct {
	Dspacefile io.Reader
}

// Generate a channel of Records.
func (m *DspaceGenerator) Generate() <-chan record.Record {
	out := make(chan record.Record)
	p := dspaceparser{file: m.Dspacefile}
	go p.parse(out)
	return out
}

// Streams the xml file and kicks off processing for each record found
func (m *dspaceparser) parse(out chan record.Record) {
	decoder := xml.NewDecoder(m.file)

	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			// If we just read a StartElement token named "record"
			if se.Name.Local == "record" {
				processMETSRecord(se, decoder, out)
			}
		}
	}
	close(out)
}

// processMETSRecord handles the mapping from OAI-PMH harvested METS to Record.
func processMETSRecord(se xml.StartElement, decoder *xml.Decoder, out chan record.Record) {
	var dr DspaceRecord
	decoder.DecodeElement(&dr, &se)

	mods := dr.Metadata.Mets.DmdSec.MdWrap.XMLData.Mods

	r := record.Record{}

	// ContentType field
	r.ContentType = mods.Genre

	// Contributor field
	if len(mods.Contributors) > 0 {
		for _, c := range mods.Contributors {
			contrib := new(record.Contributor)
			contrib.Kind = c.Role.RoleTerm.RoleName
			contrib.Value = c.NamePart.Name
			r.Contributor = append(r.Contributor, contrib)
		}
	}

	// Identifier field
	r.Identifier = strings.ReplaceAll(dr.Header.Identifier, "/", "-")

	// Citation, DOI, ISBN, ISSN, and OCLC fields
	var dois, isbns, issns, oclcs []string
	for _, i := range mods.Identifiers {
		switch {
		case i.Type == "citation":
			r.Citation = i.Value
		case i.Type == "doi":
			dois = append(dois, i.Value)
		case i.Type == "isbn":
			isbns = append(isbns, i.Value)
		case i.Type == "issn":
			issns = append(issns, i.Value)
		case i.Type == "oclc":
			oclcs = append(oclcs, i.Value)
		}
	}
	r.Doi = skipEmpty(dois)
	r.Isbn = skipEmpty(isbns)
	r.Issn = skipEmpty(issns)
	r.OclcNumber = skipEmpty(oclcs)

	// Language field
	var languages []string
	for _, l := range mods.Languages {
		languages = append(languages, l.LanguageTerm)
	}
	r.Language = skipEmpty(languages)

	// Links field
	for _, i := range mods.Identifiers {
		if i.Type == "uri" {
			uri := record.Link{
				Kind: "Digital object URI",
				URL:  i.Value,
				Text: "Digital object URI",
			}
			r.Links = append(r.Links, uri)
		}
	}

	// PublicationDate field
	r.PublicationDate = mods.OriginInfo.DateIssued

	// Source field
	r.Source = "DSpace@MIT"

	// SourceLink field
	id := dr.Header.Identifier
	linkIdentifier := strings.Split(id, "dspace.mit.edu:")[1]
	r.SourceLink = "https://hdl.handle.net/" + linkIdentifier

	// Subject field
	var subjects []string
	for _, s := range mods.Subjects {
		subjects = append(subjects, s.Topic)
	}
	r.Subject = skipEmpty(subjects)

	// Summary field (remove extra whitespace that sometimes appears in abstracts)
	space := regexp.MustCompile(`\s+`)
	abstract := space.ReplaceAllString(mods.Abstract, " ")
	r.Summary = append(r.Summary, abstract)

	// Title field
	r.Title = mods.TitleInfo.Title

	out <- r
}
