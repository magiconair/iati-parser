package main

import (
	"encoding/csv"
	"encoding/xml"
	"log"
	"os"
)

type Value struct {
	Value string `xml:",innerxml"`
	Date  string `xml:"value-date,attr"`
}

type TransactionType struct {
	Code string `xml:"code,attr"`
}

type Transaction struct {
	Type     TransactionType `xml:"transaction-type"`
	Date     string          `xml:"iso-date,attr"`
	Value    Value           `xml:"value"`
	Receiver string          `xml:"receiver-org>narrative"`
}

type CountryOrRegion struct {
	Code string `xml:"code,attr"`
	Name string `xml:"narrative"`
}

type Sector struct {
	Code string `xml:"code,attr"`
}

type PolicyMarker struct {
	Code         string `xml:"code,attr"`
	Significance string `xml:"significance,attr"`
}

type Activity struct {
	Currency      string          `xml:"default-currency,attr"`
	ID            string          `xml:"iati-identifier"`
	Title         string          `xml:"title>narrative"`
	Description   string          `xml:"description>narrative"`
	Country       CountryOrRegion `xml:"recipient-country"`
	Region        CountryOrRegion `xml:"recipient-region"`
	Transactions  []Transaction   `xml:"transaction"`
	Sector        Sector          `xml:"sector"`
	PolicyMarkers []PolicyMarker  `xml:"policy-marker"`
}

type Doc struct {
	XMLName    xml.Name   `xml:"iati-activities"`
	Activities []Activity `xml:"iati-activity"`
}

func main() {
	var d Doc
	if err := xml.NewDecoder(os.Stdin).Decode(&d); err != nil {
		log.Fatal("Decode: ", err)
	}

	w := csv.NewWriter(os.Stdout)
	w.Comma = ';'
	hdr := []string{
		"Activity Identification String",
		"CRS Sector Code",
		"Marker Gender Equal Code",
		"Marker Climate Mitigation Score",
		"Marker Climate Adaptation Score",
		"Activity Title Text",
		"Activity Description Text",
		"Country Code",
		"Country Name",
		"Region Code",
		"Region Name",
		"TX Type Code",
		"TX Date",
		"TX Value",
	}
	if err := w.Write(hdr); err != nil {
		log.Fatal("Write header: ", err)
	}

	for _, act := range d.Activities {
		// skip other countries
		cc := act.Country.Code
		if !(cc == "GH" || cc == "ET" || cc == "") {
			continue
		}

		gender, mitigation, adaptation := "0", "0", "0"
		for _, p := range act.PolicyMarkers {
			switch p.Code {
			case "1":
				gender = p.Significance
			case "6":
				mitigation = p.Significance
			case "7":
				adaptation = p.Significance
			}
		}

		for _, tx := range act.Transactions {
			rec := []string{
				act.ID,
				act.Sector.Code,
				gender,
				mitigation,
				adaptation,
				act.Title,
				act.Description,
				act.Country.Code,
				act.Country.Name,
				act.Region.Code,
				act.Region.Name,
				tx.Receiver,
				tx.Type.Code,
				tx.Value.Date,
				tx.Value.Value,
			}

			if err := w.Write(rec); err != nil {
				log.Fatal("Write: ", err)
			}
		}
	}
	w.Flush()
}
