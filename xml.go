package main

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"image/png"
)

type vmsUnitReference struct {
	Id string `xml:"id,attr"`
}

type vms struct {
	RefId vmsUnitReference `xml:"vmsUnitReference"`
	Image string           `xml:"vms>vms>vmsMessage>vmsMessage>vmsMessageExtension>vmsMessageExtension>vmsImage>imageData>binary"`
}

type location struct {
	Id          string `xml:"id,attr"`
	Description string `xml:"vmsRecord>vmsRecord>vmsDescription>values>value"`
	Latitude    string `xml:"vmsRecord>vmsRecord>vmsLocation>locationForDisplay>latitude"`
	Longitude   string `xml:"vmsRecord>vmsRecord>vmsLocation>locationForDisplay>longitude"`
}

// Used to parse XML directly into a map instead of a slice
type locationRecordMap map[string]location

func (l *locationRecordMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	r := location{}
	err := d.DecodeElement(&r, &start)
	if err != nil {
		return err
	}

	(*l)[r.Id] = r

	return nil
}

// Renders images without the location file
func imagesFromFile(file []byte) (map[string][]byte, error) {
	vmsUnits, err := parseVMsUnits(file)
	if err != nil {
		return nil, err
	}

	out := make(map[string][]byte, len(vmsUnits))
	for _, vmsUnit := range vmsUnits {
		if len(vmsUnit.Image) == 0 {
			continue
		}

		img, err := base64.StdEncoding.DecodeString(vmsUnit.Image)
		if err != nil {
			return nil, err
		}

		if len(img) == 0 {
			continue
		}
		out[vmsUnit.RefId.Id] = img
	}

	return out, nil

}

func parseLocations(locationFile []byte, expectedSize int) (locationRecordMap, error) {
	locations := struct {
		Locations locationRecordMap `xml:"Body>d2LogicalModel>payloadPublication>vmsUnitTable>vmsUnitRecord"`
	}{
		Locations: make(locationRecordMap, expectedSize),
	}

	err := xml.Unmarshal(locationFile, &locations)
	if err != nil {
		return nil, err
	}

	return locations.Locations, nil
}

func parseVMsUnits(contentFile []byte) ([]vms, error) {
	payload := struct {
		Drips []vms `xml:"Body>d2LogicalModel>payloadPublication>vmsUnit"`
	}{}
	err := xml.Unmarshal(contentFile, &payload)
	if err != nil {
		return nil, err
	}

	return payload.Drips, nil
}

func ParseDripsXML(contentFile, locationFile []byte) ([]Drip, error) {
	vmsUnits, err := parseVMsUnits(contentFile)
	if err != nil {
		return nil, err
	}
	locations, err := parseLocations(locationFile, len(vmsUnits))
	if err != nil {
		return nil, err
	}

	drips := make([]Drip, len(vmsUnits))

	for i, d := range vmsUnits {
		loc := locations[d.RefId.Id]

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

		image, err := png.Decode(bytes.NewReader(img))

		if err != nil {
			continue // Ignore faulty images
		}

		drips[i].image = img
		drips[i].ImageWidth = image.Bounds().Dx()
		drips[i].ImageHeight = image.Bounds().Dy()
	}

	return drips, nil
}
