package main

import (
	"encoding/xml"
	"io"
	"strings"
)

type archivesparser struct {
	file io.Reader
}

//ArchivesGenerator parses binary archivespace data.
type ArchivesGenerator struct {
	archivefile io.Reader
	rulesfile   string
}

//Generate a channel of Records.
func (m *ArchivesGenerator) Generate() <-chan Record {
	out := make(chan Record)
	p := archivesparser{file: m.archivefile}
	go p.parse(out)
	return out
}

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
				var ar AspaceRecord
				decoder.DecodeElement(&ar, &se)

				r := Record{}
				r.Identifier = ar.Header.Identifier
				r.Source = "MIT ArchiveSpace"
				linkIdentifier := strings.Split(r.Identifier, "oai:mit/")[1]
				r.SourceLink = "https://emmas-lib.mit.edu" + linkIdentifier

				r.Title = ar.Metadata.Mods.Titleinfo.Title

				r.Summary = ar.Metadata.Mods.Abstract

				r = gatherNotes(ar.Metadata.Mods.Note, r)

				r.PhysicalDescription = gatherPD(ar.Metadata.Mods.PhysicalDescription)

				r.Holdings = asHoldings(ar.Metadata.Mods.Location)

				r.Subject = asSubjects(ar.Metadata.Mods.Subject)

				r.Language = asLanguages(ar.Metadata.Mods.Language)

				out <- r
			}
		}
	}
	close(out)
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" && str != "\n          \n        " {
			r = append(r, str)
		}
	}
	return r
}

func asLanguages(langs struct {
	Text         string "xml:\",chardata\""
	LanguageTerm struct {
		Text      string "xml:\",chardata\""
		Authority string "xml:\"authority,attr\""
	} "xml:\"languageTerm\""
}) []string {

	var l []string

	l = append(l, langs.LanguageTerm.Text)
	l = append(l, langs.Text)

	return deleteEmpty(l)
}

func asSubjects(subjects []struct {
	Text       string "xml:\",chardata\""
	Topic      string "xml:\"topic\""
	Genre      string "xml:\"genre\""
	Geographic string "xml:\"geographic\""
	Name       struct {
		Text string "xml:\",chardata\""
		Type string "xml:\"type,attr\""
	} "xml:\"name\""
}) []string {
	var c []string

	for _, y := range subjects {
		c = append(c, y.Topic)
		c = append(c, y.Text)
		c = append(c, y.Genre)
		c = append(c, y.Geographic)
	}

	return deleteEmpty(c)
}

func asHoldings(location struct {
	Text             string "xml:\",chardata\""
	PhysicalLocation string "xml:\"physicalLocation\""
}) []Holding {

	var h []Holding

	var h1 Holding
	h1.Location = location.PhysicalLocation

	h = append(h, h1)

	return h
}

func gatherPD(pd []struct {
	Text   string "xml:\",chardata\""
	Extent string "xml:\"extent\""
}) string {
	var c []string

	for _, y := range pd {
		c = append(c, y.Extent)
	}

	return strings.Join(c, " || ")
}

func gatherNotes(notes []struct {
	Text string "xml:\",chardata\""
	Type string "xml:\"type,attr\""
}, r Record) Record {
	var c []string

	for _, y := range notes {
		if y.Type == "preferredcitation" {
			r.Citation = y.Text
		} else {
			c = append(c, y.Text)
		}
	}

	r.Notes = c

	return r
}

//AspaceRecord from XML
type AspaceRecord struct {
	XMLName xml.Name `xml:"record"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	Xsi     string   `xml:"xsi,attr"`
	Header  struct {
		Text       string `xml:",chardata"`
		Identifier string `xml:"identifier"`
		Datestamp  string `xml:"datestamp"`
	} `xml:"header"`
	Metadata struct {
		Text string `xml:",chardata"`
		Mods struct {
			Text           string `xml:",chardata"`
			Xmlns          string `xml:"xmlns,attr"`
			Xlink          string `xml:"xlink,attr"`
			SchemaLocation string `xml:"schemaLocation,attr"`
			Location       struct {
				Text             string `xml:",chardata"`
				PhysicalLocation string `xml:"physicalLocation"`
			} `xml:"location"`
			Identifier string `xml:"identifier"`
			Titleinfo  struct {
				Text  string `xml:",chardata"`
				Title string `xml:"title"`
			} `xml:"titleinfo"`
			OriginInfo struct {
				Text        string `xml:",chardata"`
				DateCreated struct {
					Text     string `xml:",chardata"`
					Encoding string `xml:"encoding,attr"`
				} `xml:"dateCreated"`
			} `xml:"originInfo"`
			PhysicalDescription []struct {
				Text   string `xml:",chardata"`
				Extent string `xml:"extent"`
			} `xml:"physicalDescription"`
			Language struct {
				Text         string `xml:",chardata"`
				LanguageTerm struct {
					Text      string `xml:",chardata"`
					Authority string `xml:"authority,attr"`
				} `xml:"languageTerm"`
			} `xml:"language"`
			AccessCondition []struct {
				Text string `xml:",chardata"`
				Type string `xml:"type,attr"`
			} `xml:"accessCondition"`
			Note []struct {
				Text string `xml:",chardata"`
				Type string `xml:"type,attr"`
			} `xml:"note"`
			Abstract []string `xml:"abstract"`
			Subject  []struct {
				Text       string `xml:",chardata"`
				Topic      string `xml:"topic"`
				Genre      string `xml:"genre"`
				Geographic string `xml:"geographic"`
				Name       struct {
					Text string `xml:",chardata"`
					Type string `xml:"type,attr"`
				} `xml:"name"`
			} `xml:"subject"`
		} `xml:"mods"`
	} `xml:"metadata"`
}
