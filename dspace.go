package main

import (
	"encoding/xml"
	"io"
	"strings"
)

type dspaceparser struct {
	file io.Reader
}

//DspaceGenerator parses binary archivespace data.
type DspaceGenerator struct {
	file      io.Reader
	rulesfile string
}

//Generate a channel of Records.
func (m *DspaceGenerator) Generate() <-chan Record {
	out := make(chan Record)
	p := dspaceparser{file: m.file}
	go p.parse(out)
	return out
}

func (m *dspaceparser) parse(out chan Record) {
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
				var dr DspaceRecord
				decoder.DecodeElement(&dr, &se)

				r := Record{}
				id := dr.Header.Identifier
				r.Identifier = "MIT.dspace." + id

				r.Source = "MIT DSpace"
				linkIdentifier := strings.Split(id, "oai:dspace.mit.edu:")[1]
				r.SourceLink = "https://dspace.mit.edu/handle/" + linkIdentifier

				// todo: confirm if this is the date to use
				r.PublicationDate = dr.Metadata.Qualifieddc.Issued

				r.ContentType = dr.Metadata.Qualifieddc.Type

				// todo: I'm not sure if Publisher belongs in imprint. It feels weird to not have a top level Record.Publisher which makes me concnered we overly MARCd this name
				var pub []string
				pub = append(pub, dr.Metadata.Qualifieddc.Publisher)
				r.Imprint = pub

				r.Title = dr.Metadata.Qualifieddc.Title

				// todo: check with metadata as to whether we use relation to mean alternate title or related items and use it appropriately once determined

				r.Contributor = dspaceContribs(dr.Metadata.Qualifieddc.Creator, dr.Metadata.Qualifieddc.Contributor)

				var summary []string
				summary = append(summary, dr.Metadata.Qualifieddc.Abstract)
				r.Summary = summary

				out <- r
			}
		}
	}
	close(out)
}

func dspaceContribs(creator []string, contribs []string) []*Contributor {
	var c []*Contributor

	for _, a := range creator {
		y := new(Contributor)
		y.Kind = "Creator"
		y.Value = a

		if y.Value != "" {
			c = append(c, y)
		}
	}

	for _, a := range contribs {
		y := new(Contributor)
		y.Kind = "Contributor"
		y.Value = a

		if y.Value != "" {
			c = append(c, y)
		}
	}

	return c
}

// DspaceRecord maps data from oai-pmh harvests in XML to Go
type DspaceRecord struct {
	XMLName xml.Name `xml:"record"`
	Text    string   `xml:",chardata"`
	Header  struct {
		Text       string   `xml:",chardata"`
		Identifier string   `xml:"identifier"`
		Datestamp  string   `xml:"datestamp"`
		SetSpec    []string `xml:"setSpec"`
	} `xml:"header"`
	Metadata struct {
		Text        string `xml:",chardata"`
		Qualifieddc struct {
			Text           string   `xml:",chardata"`
			Qdc            string   `xml:"qdc,attr"`
			Doc            string   `xml:"doc,attr"`
			Xsi            string   `xml:"xsi,attr"`
			Dcterms        string   `xml:"dcterms,attr"`
			Dc             string   `xml:"dc,attr"`
			SchemaLocation string   `xml:"schemaLocation,attr"`
			Title          string   `xml:"title"`
			Creator        []string `xml:"creator"`
			Contributor    []string `xml:"contributor"`
			Abstract       string   `xml:"abstract"`
			DateAccepted   string   `xml:"dateAccepted"`
			Available      string   `xml:"available"`
			Created        string   `xml:"created"`
			Issued         string   `xml:"issued"`
			Type           string   `xml:"type"`
			Identifier     string   `xml:"identifier"`
			Relation       string   `xml:"relation"`
			Publisher      string   `xml:"publisher"`
		} `xml:"qualifieddc"`
	} `xml:"metadata"`
}
