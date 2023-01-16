package description

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseDescription(t *testing.T) {
	tests := []struct {
		name string
		args string
		want DescriptionDerivatives
	}{
		{
			name: "It removes parentheses",
			args: "Plutoniumweg (dBD36)",
			want: DescriptionDerivatives{
				Organization: "",
				RoadId:       "",
				RoadOffset:   -1,
				RoadSide:     "",
				Name:         "Plutoniumweg",
			},
		},
		{
			name: "Leaves plain names alone",
			args: "Waterlandlaan",
			want: DescriptionDerivatives{
				Organization: "",
				RoadId:       "",
				RoadOffset:   -1,
				RoadSide:     "",
				Name:         "Waterlandlaan",
			},
		},
		{
			name: "Parses organization names",
			args: "PZH_DRIP65 - Hoefweg Veiling Bleiswijk",
			want: DescriptionDerivatives{
				Organization: "Provincie Zuid-Holland",
				RoadId:       "",
				RoadOffset:   -1,
				RoadSide:     "",
				Name:         "Hoefweg Veiling Bleiswijk",
			},
		},
		{
			name: "Removes blacklisted words",
			args: "201309 - Rechts",
			want: DescriptionDerivatives{
				Organization: "",
				RoadId:       "",
				RoadOffset:   -1,
				RoadSide:     "R",
				Name:         "",
			},
		},
		{
			name: "Parses road id, side and offset",
			args: "PZH_DRIP14 - N211 R 12.7 Poeldijk (9eff9e60-3ece-4abd-84cb-be319504e1)",
			want: DescriptionDerivatives{
				Organization: "Provincie Zuid-Holland",
				RoadId:       "N211",
				RoadOffset:   12700,
				RoadSide:     "R",
				Name:         "Poeldijk",
			},
		},
		{
			name: "Parses road data from dashes",
			args: "A2-Li-62,2",
			want: DescriptionDerivatives{
				Organization: "",
				RoadId:       "A2",
				RoadOffset:   62200,
				RoadSide:     "L",
				Name:         "",
			},
		},
		{
			name: "Handles reverse form offsets",
			args: "A020-21_650-Re-1-4 - A20 Re km 21,650",
			want: DescriptionDerivatives{
				Organization: "",
				RoadId:       "A20",
				RoadOffset:   21650,
				RoadSide:     "R",
				Name:         "",
			},
		},
		{
			name: "Handles partial road data",
			args: "45801 - N471 (c9414daf-9e90-4ba6-b475-89dabde2e8fa)",
			want: DescriptionDerivatives{
				Organization: "",
				RoadId:       "N471",
				RoadOffset:   -1,
				RoadSide:     "",
				Name:         "",
			},
		},
		{
			name: "Handles m suffixes in road offset",
			args: "GDH_QW-18-06 - A4 Re 44,570m parallelbaan voor A4/A12 knp Prins Clausplein (74c46760-51c7-4187-827d-d020dc112133)",
			want: DescriptionDerivatives{
				Organization: "Gemeente Den Haag",
				RoadId:       "A4",
				RoadOffset:   44570,
				RoadSide:     "R",
				Name:         "parallelbaan voor A4/A12 knp Prins Clausplein",
			},
		},
		{
			name: "Handles side suffixes",
			args: "26011 - N223L km 7.7 Twee Pleinenweg (00cf2dd7-a451-4165-8ded-10ab070e7371)",
			want: DescriptionDerivatives{
				Organization: "",
				RoadId:       "N223",
				RoadOffset:   7700,
				RoadSide:     "L",
				Name:         "Twee Pleinenweg",
			},
		},
		{
			name: "Handles offset suffixes",
			args: "A44R_16,700 (dBD179)",
			want: DescriptionDerivatives{
				Organization: "",
				RoadId:       "A44",
				RoadOffset:   16700,
				RoadSide:     "R",
				Name:         "",
			},
		},
		{
			name: "Handles no space format",
			args: "A16L69,900",
			want: DescriptionDerivatives{
				Organization: "",
				RoadId:       "A16",
				RoadOffset:   69900,
				RoadSide:     "L",
				Name:         "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.args); !reflect.DeepEqual(got, tt.want) {
				gotstr, _ := json.MarshalIndent(got, "", "\t")
				wantstr, _ := json.MarshalIndent(tt.want, "", "\t")
				t.Errorf("\nGot:\n%v\nExpected:\n%v", string(gotstr), string(wantstr))
			}
		})
	}
}
