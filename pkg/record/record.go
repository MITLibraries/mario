package record

// Record struct stores our internal mappings of data and is used when
// mapping various external data sources before sending to OpenSearch
type Record struct {
	AlternateTitles        []*AlternateTitle `json:"alternate_titles,omitempty"`
	CallNumbers            []string          `json:"call_numbers,omitempty"`
	Citation               string            `json:"citation,omitempty"`
	ContentType            []string          `json:"content_type,omitempty"`
	Contents               []string          `json:"contents,omitempty"`
	Contributors           []*Contributor    `json:"contributors,omitempty"`
	Dates                  []*Date           `json:"dates,omitempty"`
	Edition                string            `json:"edition,omitempty"`
	FileFormats            []string          `json:"file_formats,omitempty"`
	Format                 string            `json:"format,omitempty"`
	FundingInformation     []*Funding        `json:"funding_information,omitempty"`
	Holdings               []*Holding        `json:"holdings,omitempty"`
	Identifiers            []*Identifier     `json:"identifiers,omitempty"`
	Languages              []string          `json:"languages,omitempty"`
	Links                  []Link            `json:"links,omitempty"`
	LiteraryForm           string            `json:"literary_form,omitempty"`
	Locations              []*Location       `json:"locations,omitempty"`
	Notes                  []*Note           `json:"notes,omitempty"`
	Numbering              string            `json:"numbering,omitempty"`
	PhysicalDescription    string            `json:"physical_description,omitempty"`
	PublicationFrequency   []string          `json:"publication_frequency,omitempty"`
	PublicationInformation []string          `json:"publication_information,omitempty"`
	RelatedItems           []*RelatedItem    `json:"related_items,omitempty"`
	Rights                 []*Right          `json:"rights,omitempty"`
	Source                 string            `json:"source"`
	SourceLink             string            `json:"source_link"`
	Subjects               []*Subject        `json:"subjects,omitempty"`
	Summary                []string          `json:"summary,omitempty"`
	TimdexRecordId         string            `json:"timdex_record_id"`
	Title                  string            `json:"title"`
}

// AlternateTitle object
type AlternateTitle struct {
	Kind  string `json:"kind,omitempty"`
	Value string `json:"value"`
}

// Contributor object
type Contributor struct {
	Affiliation   string `json:"affiliation,omitempty"`
	Kind          string `json:"kind,omitempty"`
	Identifier    string `json:"identifier,omitempty"`
	MitAffiliated bool   `json:"mit_affiliated,omitempty"`
	Value         string `json:"value"`
}

// Date object
type Date struct {
	Kind  string `json:"kind,omitempty"`
	Note  string `json:"note,omitempty"`
	Range *Range `json:"range,omitempty"`
	Value string `json:"value,omitempty"`
}

// Funding object
type Funding struct {
	AwardNumber          string `json:"award_number,omitempty"`
	AwardUri             string `json:"award_uri,omitempty"`
	FunderIdentifier     string `json:"funder_identifier,omitempty"`
	FunderIdentifierType string `json:"funder_identifier_type,omitempty"`
	FunderName           string `json:"funder_name,omitempty"`
}

// Holding object
type Holding struct {
	CallNumber string `json:"call_number,omitempty"`
	Collection string `json:"collection,omitempty"`
	Format     string `json:"format,omitempty"`
	Location   string `json:"location,omitempty"`
	Note       string `json:"notes,omitempty"`
	Summary    string `json:"summary,omitempty"`
}

// Identifier object
type Identifier struct {
	Kind  string `json:"kind,omitmempty"`
	Value string `json:"value"`
}

// Link object
type Link struct {
	Kind         string `json:"kind,omitempty"`
	Restrictions string `json:"restrictions,omitempty"`
	Text         string `json:"text,omitempty"`
	Url          string `json:"url"`
}

// Location object
type Location struct {
	Geopoint []float32 `json:"geopoint,omitempty"`
	Kind     string    `json:"kind,omitempty"`
	Value    string    `json:"value,omitempty"`
}

// Note object
type Note struct {
	Kind  string   `json:"kind,omitempty"`
	Value []string `json:"value"`
}

// Range object
type Range struct {
	Gt  string `json:"gt,omitempty"`
	Gte string `json:"gte,omitempty"`
	Lt  string `json:"lt,omitempty"`
	Lte string `json:"lte,omitempty"`
}

// RelatedItem object
type RelatedItem struct {
	Description  string `json:"description,omitempty"`
	ItemType     string `json:"item_type,omitempty"`
	Relationship string `json:"relationship,omitempty"`
	Uri          string `json:"uri,omitempty"`
}

// Right object
type Right struct {
	Description string `json:"desription,omitempty"`
	Kind        string `json:"kind,omitempty"`
	Uri         string `json:"uri,omitempty"`
}

// Subject object
type Subject struct {
	Kind  string   `json:"kind,omitempty"`
	Value []string `json:"value"`
}

// Rule defines where the rules are in JSON
type Rule struct {
	Label  string   `json:"label"`
	Array  bool     `json:"array"`
	Fields []*Field `json:"fields"`
}

// Field defines where the Fields within a Rule are in JSON
type Field struct {
	Tag       string `json:"tag"`
	Subfields string `json:"subfields"`
	Bytes     string `json:"bytes"`
	Kind      string `json:"kind"`
}

// Parser defines an interface common to parsers
type Parser interface {
	Parse(chan Record)
}

// Processor is an interface that allows converting from custom data into
// our Record structure
type Processor interface {
	Process()
}

var ingested int
