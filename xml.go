package main

import (
	"encoding/base64"
	"encoding/xml"
)

type LocationRecordMap map[string]Location

type Location struct {
	Description, Latitude, Longitude string
}

type XMLPayloadPublication struct {
	Drips []XMLvms `xml:"Body>d2LogicalModel>payloadPublication>vmsUnit"`
}

type XMLvmsUnitReference struct {
	Id string `xml:"id,attr"`
}

type XMLvms struct {
	RefId XMLvmsUnitReference `xml:"vmsUnitReference"`
	Image string              `xml:"vms>vms>vmsMessage>vmsMessage>vmsMessageExtension>vmsMessageExtension>vmsImage>imageData>binary"`
}

type XMLVMSRecord struct {
	Id          string `xml:"id,attr"`
	Description string `xml:"vmsRecord>vmsRecord>vmsDescription>values>value"`
	Latitude    string `xml:"vmsRecord>vmsRecord>vmsLocation>locationForDisplay>latitude"`
	Longitude   string `xml:"vmsRecord>vmsRecord>vmsLocation>locationForDisplay>longitude"`
}

type XMLUnitTable struct {
	Locations LocationRecordMap `xml:"Body>d2LogicalModel>payloadPublication>vmsUnitTable>vmsUnitRecord"`
}

func (l *LocationRecordMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	r := XMLVMSRecord{}
	err := d.DecodeElement(&r, &start)
	if err != nil {
		return err
	}

	(*l)[r.Id] = Location{
		Description: r.Description,
		Latitude:    r.Latitude,
		Longitude:   r.Longitude}

	return nil
}

func parseDripsXML(contentFile, locationFile []byte) ([]Drip, error) {
	payload := XMLPayloadPublication{}
	err := xml.Unmarshal(contentFile, &payload)
	if err != nil {
		return []Drip{}, err
	}

	locations := XMLUnitTable{
		Locations: make(LocationRecordMap, len(payload.Drips)),
	}
	err2 := xml.Unmarshal(locationFile, &locations)
	if err2 != nil {
		return []Drip{}, err2
	}

	drips := make([]Drip, len(payload.Drips))

	for i, d := range payload.Drips {
		loc := locations.Locations[d.RefId.Id]

		drips[i] = Drip{
			Id:          d.RefId.Id,
			Lat:         loc.Latitude,
			Lon:         loc.Longitude,
			Description: loc.Description,
		}

		img, err := base64.StdEncoding.DecodeString(d.Image)
		if err != nil {
			continue // Ignore faulty images
		}

		if len(img) == 0 {
			continue // Ignore empty images
		}

		drips[i].image = img
	}

	return drips, nil
}
