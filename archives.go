package main

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/markbates/pkger"
	yaml "gopkg.in/yaml.v2"
)

type archivesparser struct {
	file io.Reader
}

// ArchivesGenerator parses archivespace ead xml data.
type ArchivesGenerator struct {
	archivefile io.Reader
	rulesfile   string
}

// Generate a channel of Records.
func (m *ArchivesGenerator) Generate() <-chan Record {
	out := make(chan Record)
	p := archivesparser{file: m.archivefile}
	go p.parse(out)
	return out
}

// Streams the xml file and kicks off processing for each record found
func (m *archivesparser) parse(out chan Record) {
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
				processXMLRecord(se, decoder, out)
			}
		}
	}
	close(out)
}

// processXMLRecord handles the mapping from EAD to Record. More complex mappings split out into funcs
func processXMLRecord(se xml.StartElement, decoder *xml.Decoder, out chan Record) {
	var ar AspaceRecord
	decoder.DecodeElement(&ar, &se)

	r := Record{}

	// Citation field
	r.Citation = ar.Metadata.Ead.Archdesc.Prefercite.P.Text

	// ContentType field
	r.ContentType = "Archival " + ar.Metadata.Ead.Archdesc.Level

	// Contributor field
	if len(ar.Metadata.Ead.Archdesc.Did.Origination) > 0 {
		r.Contributor = eadContributors(ar)
	}

	//  Holdings field
	var h []Holding
	h = append(h, Holding{Location: ar.Metadata.Ead.Archdesc.Did.Physloc.Text})
	r.Holdings = h

	// Identifier field
	r.Identifier = "MIT:archivesspace:" + strings.Replace(ar.Metadata.Ead.Archdesc.Did.Unitid, " ", ".", -1)

	// Language field
	if len(ar.Metadata.Ead.Archdesc.Did.Langmaterial) > 0 {
		r.Language = eadLanguage(ar)
	}

	// Links field
	r.Links = eadLinks(ar)

	// Notes field
	r.Notes = eadNotes(ar)

	//  PhysicalDescription field
	if len(ar.Metadata.Ead.Archdesc.Did.Physdesc) > 0 {
		r.PhysicalDescription = eadPhysicalDescription(ar)
	}

	// PublicationDate field
	r.PublicationDate = eadPublicationDate(ar)

	// Source field
	r.Source = "MIT ArchivesSpace"

	// SourceLink field
	id := ar.Header.Identifier
	linkIdentifier := strings.Split(id, "oai:mit/")[1]
	r.SourceLink = "https://archivesspace.mit.edu" + linkIdentifier

	// Subject field
	r.Subject = eadSubjects(ar)

	// Summary field
	if len(ar.Metadata.Ead.Archdesc.Did.Abstract) > 0 {
		for _, a := range ar.Metadata.Ead.Archdesc.Did.Abstract {
			r.Summary = append(r.Summary, a.Text)
		}
	}

	// Title field
	r.Title = ar.Metadata.Ead.Archdesc.Did.Unittitle.Text

	out <- r
}

// AspaceCodesMap defines codes for parsing ASpace record fields
type AspaceCodesMap struct {
	Enumerations struct {
		LinkedAgentRelators map[string]string `yaml:"linked_agent_archival_record_relators"`
	} `yaml:"enumerations"`
}

func eadContributors(ar AspaceRecord) []*Contributor {
	var contribs []*Contributor
	var codes AspaceCodesMap

	file, err := pkger.Open("/config/aspace_code_mappings.yml")
	if err != nil {
		panic(err)
	}
	yamlFile, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &codes)
	if err != nil {
		panic(err)
	}

	for _, c := range ar.Metadata.Ead.Archdesc.Did.Origination {
		contrib := new(Contributor)
		switch {
		case c.Corpname.Text != "":
			contrib.Kind = codes.Enumerations.LinkedAgentRelators[c.Corpname.Role]
			contrib.Value = c.Corpname.Text
		case c.Famname.Text != "":
			contrib.Kind = codes.Enumerations.LinkedAgentRelators[c.Famname.Role]
			contrib.Value = c.Famname.Text
		case c.Persname.Text != "":
			contrib.Kind = codes.Enumerations.LinkedAgentRelators[c.Persname.Role]
			contrib.Value = c.Persname.Text
		}
		if contrib.Kind == "" {
			contrib.Kind = c.Label
		}

		contribs = append(contribs, contrib)
	}

	return contribs
}

func eadLanguage(ar AspaceRecord) []string {
	var lang []string

	for _, l := range ar.Metadata.Ead.Archdesc.Did.Langmaterial {
		lang = append(lang, l.Text)
		lang = append(lang, l.Language.Text)
	}
	return skipEmpty(lang)
}

func eadLinks(ar AspaceRecord) []Link {
	var links []Link

	dsc, _ := xmlquery.Parse(strings.NewReader(ar.Metadata.Ead.Archdesc.Dsc.Text))
	dao := xmlquery.Find(dsc, "//dao")

	for _, obj := range dao {
		link := Link{
			URL:  obj.SelectAttr("xlink:href"),
			Text: obj.SelectElement("daodesc").SelectElement("p").InnerText(),
			Kind: "Digital object",
		}

		// only keep links that start with http. This isn't ideal, but seems okay.
		if strings.HasPrefix(link.URL, "http") {
			links = append(links, link)
		}
	}
	return links
}

func eadNotes(ar AspaceRecord) []string {
	var notes []string

	if len(ar.Metadata.Ead.Archdesc.Accessrestrict) > 0 {
		for _, n := range ar.Metadata.Ead.Archdesc.Accessrestrict {
			note := n.Head + ": " + n.P
			notes = append(notes, note)
		}
	}

	if len(ar.Metadata.Ead.Archdesc.Userestrict) > 0 {
		for _, n := range ar.Metadata.Ead.Archdesc.Userestrict {
			note := n.Head + ": " + n.P
			notes = append(notes, note)
		}
	}

	if len(ar.Metadata.Ead.Archdesc.Bioghist) > 0 {
		for _, n := range ar.Metadata.Ead.Archdesc.Bioghist {
			if len(n.P) > 0 {
				var note []string
				note = append(note, n.Head)
				for _, t := range n.P {
					note = append(note, t)
				}
				notes = append(notes, strings.Join(note, "\n"))
			}
		}
	}
	return notes
}

func eadPhysicalDescription(ar AspaceRecord) string {
	var pd string
	for _, p := range ar.Metadata.Ead.Archdesc.Did.Physdesc {
		if len(p.Extent) > 0 {
			for _, e := range p.Extent {
				var joiner string
				if pd == "" {
					joiner = ""
				} else {
					joiner = "; "
				}
				pd = pd + joiner + strings.Trim(e.Text, "()")
			}
		}
	}
	return pd
}

func eadPublicationDate(ar AspaceRecord) string {
	var date []string

	if len(ar.Metadata.Ead.Archdesc.Did.Unitdate) > 0 {
		for _, d := range ar.Metadata.Ead.Archdesc.Did.Unitdate {
			date = append(date, d.Text)
		}
	}
	return strings.Join(date, ",")
}

func eadSubjects(ar AspaceRecord) []string {
	var subjects []string

	// Subject
	if len(ar.Metadata.Ead.Archdesc.Controlaccess.Subject) > 0 {
		for _, s := range ar.Metadata.Ead.Archdesc.Controlaccess.Subject {
			subjects = append(subjects, s.Text)
		}
	}

	// corpname
	if len(ar.Metadata.Ead.Archdesc.Controlaccess.Corpname) > 0 {
		for _, s := range ar.Metadata.Ead.Archdesc.Controlaccess.Corpname {
			subjects = append(subjects, s.Text)
		}
	}

	// famname
	if len(ar.Metadata.Ead.Archdesc.Controlaccess.Famname) > 0 {
		for _, s := range ar.Metadata.Ead.Archdesc.Controlaccess.Famname {
			subjects = append(subjects, s.Text)
		}
	}

	// geogname
	if len(ar.Metadata.Ead.Archdesc.Controlaccess.Geogname) > 0 {
		for _, s := range ar.Metadata.Ead.Archdesc.Controlaccess.Geogname {
			subjects = append(subjects, s.Text)
		}
	}

	// name
	subjects = append(subjects, ar.Metadata.Ead.Archdesc.Controlaccess.Title.Text)

	// persname
	if len(ar.Metadata.Ead.Archdesc.Controlaccess.Persname) > 0 {
		for _, s := range ar.Metadata.Ead.Archdesc.Controlaccess.Persname {
			subjects = append(subjects, s.Text)
		}
	}

	return skipEmpty(subjects)
}

func skipEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if strings.TrimSpace(str) != "" {
			r = append(r, str)
		}
	}
	return r
}
