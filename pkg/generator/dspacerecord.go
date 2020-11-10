package generator

import "encoding/xml"

// DspaceRecord MODS to struct mappings.
type DspaceRecord struct {
	XMLName xml.Name `xml:"record"`
	Xmlns   string   `xml:"xmlns,attr"`
	Xsi     string   `xml:"xsi,attr"`
	Header  struct {
		Identifier string    `xml:"identifier"` // oai:dspace.mit.edu:1721...
		Datestamp  string    `xml:"datestamp"`  // 2019-07-03T19:09:32Z, 201...
		SetSpecs   []SetSpec `xml:"setSpec"`    // com_1721.1_...
	} `xml:"header"`
	Metadata struct {
		Mets struct {
			Xmlns              string `xml:"xmlns,attr"`
			Doc                string `xml:"doc,attr"`
			Xlink              string `xml:"xlink,attr"`
			MetsSchemaLocation string `xml:"schemaLocation,attr"`
			DmdSec             struct {
				MdWrap struct {
					XMLData struct {
						ModsNS string `xml:"mods,attr"`
						Mods   struct {
							Contributors []Contributor `xml:"name"`
							OriginInfo   struct {
								DateIssued string `xml:"dateIssued"`
							} `xml:"originInfo"`
							Identifiers []Identifier `xml:"identifier"`
							Abstract    string       `xml:"abstract"`
							Languages   []Language   `xml:"language"`
							Subjects    []Subject    `xml:"subject"`
							TitleInfo   struct {
								Title string `xml:"title"`
							} `xml:"titleInfo"`
							Genre string `xml:"genre"`
						} `xml:"mods"`
					} `xml:"xmlData"`
				} `xml:"mdWrap"`
			} `xml:"dmdSec"`
		} `xml:"mets"`
	} `xml:"metadata"`
}

// Contributor field mapping
type Contributor struct {
	NamePart struct {
		Name string `xml:",chardata"`
	} `xml:"namePart"`
	Role struct {
		RoleTerm struct {
			RoleName string `xml:",chardata"`
		} `xml:"roleTerm"`
	} `xml:"role"`
}

// Identifier field mapping
type Identifier struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

// Language field mapping
type Language struct {
	LanguageTerm string `xml:"languageTerm"`
}

// SetSpec field mapping
type SetSpec struct {
	SetSpec string `xml:",chardata"`
}

// Subject field mapping
type Subject struct {
	Topic string `xml:"topic"`
}
