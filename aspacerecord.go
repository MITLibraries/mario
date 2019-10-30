package main

// AspaceRecord ead to struct mappings. This was partically generated by zek
// but then modified to remove bits we don't care about and in some places
// change how zek intepretted the data structures due to some interesting
// decisions in how EAD works.
type AspaceRecord struct {
	Text   string `xml:",chardata"`
	Xmlns  string `xml:"xmlns,attr"`
	Xsi    string `xml:"xsi,attr"`
	Header struct {
		Text       string `xml:",chardata"`
		Identifier string `xml:"identifier"` // oai:mit//repositories/2/r...
		Datestamp  string `xml:"datestamp"`  // 2019-07-03T19:09:32Z, 201...
	} `xml:"header"`
	Metadata struct {
		Text string `xml:",chardata"`
		Ead  struct {
			Text           string `xml:",chardata"`
			Xmlns          string `xml:"xmlns,attr"`
			Xlink          string `xml:"xlink,attr"`
			Xsi            string `xml:"xsi,attr"`
			SchemaLocation string `xml:"schemaLocation,attr"`
			Eadheader      struct {
				Text               string `xml:",chardata"`
				Countryencoding    string `xml:"countryencoding,attr"`
				Dateencoding       string `xml:"dateencoding,attr"`
				Findaidstatus      string `xml:"findaidstatus,attr"`
				Langencoding       string `xml:"langencoding,attr"`
				Repositoryencoding string `xml:"repositoryencoding,attr"`
				Eadid              struct {
					Text           string `xml:",chardata"` // AC 1, AC 2, AC 3, AC 4, A...
					Countrycode    string `xml:"countrycode,attr"`
					Mainagencycode string `xml:"mainagencycode,attr"`
					URL            string `xml:"url,attr"`
				} `xml:"eadid"`
			} `xml:"eadheader"`
			Archdesc struct {
				Text  string `xml:",chardata"`
				Level string `xml:"level,attr"`
				Did   struct {
					Text         string `xml:",chardata"`
					Langmaterial []struct {
						Text     string `xml:",chardata"` // The material is primarily...
						ID       string `xml:"id,attr"`
						Language struct {
							Text     string `xml:",chardata"` // English, English, English...
							Langcode string `xml:"langcode,attr"`
						} `xml:"language"`
					} `xml:"langmaterial"`
					Repository struct {
						Text     string `xml:",chardata"`
						Corpname string `xml:"corpname"` // Massachusetts Institute o...
					} `xml:"repository"`
					Unittitle struct {
						Text  string `xml:",chardata"` // Charles J. Connick Staine...
						Title struct {
							Text   string `xml:",chardata"` // Science, Technology and H...
							Render string `xml:"render,attr"`
						} `xml:"title"`
					} `xml:"unittitle"`
					Unitid   string `xml:"unitid"` // (ROTCH LIBRARY).Connick, ...
					Physdesc []struct {
						Text      string `xml:",chardata"` // The records of the MIT Fl...
						Altrender string `xml:"altrender,attr"`
						ID        string `xml:"id,attr"`
						Label     string `xml:"label,attr"`
						Extent    []struct {
							Text      string `xml:",chardata"` // 31 box(es), 40 Megabytes,...
							Altrender string `xml:"altrender,attr"`
						} `xml:"extent"`
						Physfacet struct {
							Text  string `xml:",chardata"` // Minutes in volumes 1 to 1...
							ID    string `xml:"id,attr"`
							Label string `xml:"label,attr"`
						} `xml:"physfacet"`
					} `xml:"physdesc"`
					Unitdate []struct {
						Text string `xml:",innerxml"` // 1905-2012, 1865-2013, 187...
					} `xml:"unitdate"`
					Origination []struct {
						Text     string `xml:",chardata"`
						Label    string `xml:"label,attr"`
						Corpname struct {
							Text           string `xml:",chardata"` // Massachusetts Institute o...
							Rules          string `xml:"rules,attr"`
							Source         string `xml:"source,attr"`
							Role           string `xml:"role,attr"`
							Authfilenumber string `xml:"authfilenumber,attr"`
						} `xml:"corpname"`
						Persname struct {
							Text           string `xml:",chardata"` // Pounds, W. F. (William F....
							Rules          string `xml:"rules,attr"`
							Source         string `xml:"source,attr"`
							Role           string `xml:"role,attr"`
							Authfilenumber string `xml:"authfilenumber,attr"`
						} `xml:"persname"`
						Famname struct {
							Text   string `xml:",chardata"` // Rogers, Wigglesworth
							Rules  string `xml:"rules,attr"`
							Source string `xml:"source,attr"`
							Role   string `xml:"role,attr"`
						} `xml:"famname"`
					} `xml:"origination"`
					Abstract []struct {
						Text string `xml:",innerxml"`
					} `xml:"abstract"`
					Physloc struct {
						Text string `xml:",chardata"` // Materials are stored off-...
						ID   string `xml:"id,attr"`
					} `xml:"physloc"`
					Materialspec struct {
						Text string `xml:",chardata"` // 10 1/2 in reels, 7 1/2 ip...
						ID   string `xml:"id,attr"`
					} `xml:"materialspec"`
				} `xml:"did"`
				Bioghist []struct {
					Text string   `xml:",chardata"` // The Biotechnology Process...
					ID   string   `xml:"id,attr"`
					Head string   `xml:"head"` // Biographical / Historical...
					P    []string `xml:"p"`    // Access to collections in ...
					List []struct {
						Text    string `xml:",chardata"`
						Type    string `xml:"type,attr"`
						Head    string `xml:"head"` // Personnel in the Presiden...
						Defitem []struct {
							Text  string `xml:",chardata"`
							Label string `xml:"label"` // Bowditch, Ebenezer Franci...
							Item  string `xml:"item"`  // Special Advisor to the Pr...
						} `xml:"defitem"`
						Item []string `xml:"item"` // William L. Campbell, 1945...
					} `xml:"list"`
					Chronlist []struct {
						Text      string `xml:",chardata"`
						Head      string `xml:"head"` // Presidents of the Institu...
						Chronitem []struct {
							Text     string `xml:",chardata"`
							Date     string `xml:"date"` // 1862-1970, 1870-1878, 187...
							Eventgrp struct {
								Text  string `xml:",chardata"`
								Event []struct {
									Text  string `xml:",chardata"` // William Barton Rogers, Jo...
									Title []struct {
										Text   string `xml:",chardata"` // Encyclopedia Americana, B...
										Render string `xml:"render,attr"`
									} `xml:"title"`
								} `xml:"event"`
							} `xml:"eventgrp"`
						} `xml:"chronitem"`
					} `xml:"chronlist"`
					Extref struct {
						Text string `xml:",chardata"` // homepage/
						Href string `xml:"href,attr"`
					} `xml:"extref"`
				} `xml:"bioghist"`

				Dsc struct {
					Text string `xml:",innerxml"`
				} `xml:"dsc"`

				Accessrestrict []struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
					Head string `xml:"head"` // Access note, Access note,...
					P    string `xml:"p"`    // Access to collections in ...
				} `xml:"accessrestrict"`
				Userestrict []struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
					Head string `xml:"head"` // Intellectual Property Rig...
					P    string `xml:"p"`    // Access to collections in ...
				} `xml:"userestrict"`
				Prefercite struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
					Head string `xml:"head"` // Citation, Citation, Citat...
					P    struct {
						Text string `xml:",innerxml"` // Massachusetts Institute o...
					} `xml:"p"`
				} `xml:"prefercite"`
				Controlaccess struct {
					Text    string `xml:",chardata"`
					Subject []struct {
						Text   string `xml:",chardata"` // curricula, Massachusetts ...
						Source string `xml:"source,attr"`
					} `xml:"subject"`
					Function []struct {
						Text   string `xml:",chardata"` // faculty governance, admin...
						Source string `xml:"source,attr"`
					} `xml:"function"`
					Corpname []struct {
						Text           string `xml:",chardata"` // Massachusetts Institute o...
						Rules          string `xml:"rules,attr"`
						Source         string `xml:"source,attr"`
						Authfilenumber string `xml:"authfilenumber,attr"`
						Role           string `xml:"role,attr"`
					} `xml:"corpname"`
					Genreform []struct {
						Text   string `xml:",chardata"` // speeches, speeches, speec...
						Source string `xml:"source,attr"`
					} `xml:"genreform"`
					Persname []struct {
						Text           string `xml:",chardata"` // Voss, Walter C. (Walter C...
						Rules          string `xml:"rules,attr"`
						Source         string `xml:"source,attr"`
						Authfilenumber string `xml:"authfilenumber,attr"`
						Role           string `xml:"role,attr"`
					} `xml:"persname"`
					Famname []struct {
						Text   string `xml:",chardata"` // Compton, Rogers, Wigglesw...
						Rules  string `xml:"rules,attr"`
						Source string `xml:"source,attr"`
					} `xml:"famname"`
					Geogname []struct {
						Text   string `xml:",chardata"` // Somerville (Mass.) -- Ind...
						Source string `xml:"source,attr"`
					} `xml:"geogname"`
					Title struct {
						Text   string `xml:",chardata"` // Louisiana Purchase Exposi...
						Source string `xml:"source,attr"`
					} `xml:"title"`
				} `xml:"controlaccess"`
			} `xml:"archdesc"`
		} `xml:"ead"`
	} `xml:"metadata"`
}
