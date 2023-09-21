package main

import (
	"os"
	"testing"
)

func assert[T comparable](t *testing.T, real, expected T) {
	t.Helper()

	if real != expected {
		t.Errorf("Expected %v to equal %v\n", real, expected)
	}
}

func TestXMLParsing(t *testing.T) {
	vmsUnits, err := os.ReadFile("./testdata/vmsUnit.xml")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	vmsRecords, err := os.ReadFile("./testdata/vmsRecord.xml")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	drips, err := ParseDripsXML(vmsUnits, vmsRecords)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(drips) != 3 {
		t.Errorf("Expected a lenght of 3, not %v\n", len(drips))
	}

	drip1, drip2, drip3 := drips[0], drips[1], drips[2]

	assert(t, drip1.Id, "ID_1")
	assert(t, drip2.Id, "ID_2")
	assert(t, drip3.Id, "ID_3")

	assert(t, drip2.hasImage(), true)
	assert(t, drip2.ImageWidth, 40)
	assert(t, drip2.ImageHeight, 40)

	assert(t, drip1.Lat, "52.1")
	assert(t, drip1.Lon, "4.2")

	assert(t, drip1.Name, "Description 1")

	assert(t, drip1.Working, true)
	assert(t, drip2.Working, true)
	assert(t, drip3.Working, false)

	assert(t, drip1.TextLines[0], "Textline 1")
	assert(t, drip1.TextLines[1], "Textline 2")
	assert(t, drip1.TextLines[2], "Textline 3")

}
